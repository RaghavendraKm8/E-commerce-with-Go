package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Order struct {
	ID        int64 `json:"id"`
	UserID    int64 `json:"user_id"`
	ProductID int64 `json:"product_id"`
	Qty       int   `json:"qty"`
	Total     int64 `json:"total"`
}

var (
	db  *pgxpool.Pool
	req = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "ordersvc_http_requests_total", Help: "Total HTTP requests"},
		[]string{"method", "path", "status"},
	)
	dur = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Name: "ordersvc_http_duration_seconds", Help: "Request duration", Buckets: prometheus.DefBuckets},
		[]string{"method", "path"},
	)
)

func init() { prometheus.MustRegister(req, dur) }

func main() {
	lg, _ := zap.NewProduction()
	defer lg.Sync()

	pg := os.Getenv("PG_URL")
	if pg == "" {
		pg = "postgres://postgres:postgres@localhost:5432/ecom_orders?sslmode=disable"
	}
	ctx := context.Background()
	var err error
	db, err = pgxpool.New(ctx, pg)
	if err != nil {
		lg.Fatal("db", zap.Error(err))
	}
	defer db.Close()

	r := gin.New()
	r.Use(gin.Recovery(), logmw(lg), metmw())
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.POST("/orders", create)
	r.GET("/orders/:id", get)

	srv := &http.Server{Addr: ":8083", Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			lg.Fatal("srv", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctxShut, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctxShut)
}

func create(c *gin.Context) {
	var o Order
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := db.QueryRow(c,
		"INSERT INTO orders(user_id, product_id, qty, total) VALUES($1,$2,$3,$4) RETURNING id",
		o.UserID, o.ProductID, o.Qty, o.Total).Scan(&o.ID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, o)
}

func get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var o Order
	if err := db.QueryRow(c,
		"SELECT id, user_id, product_id, qty, total FROM orders WHERE id=$1", id).
		Scan(&o.ID, &o.UserID, &o.ProductID, &o.Qty, &o.Total); err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	c.JSON(200, o)
}

func logmw(l *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		st := time.Now()
		c.Next()
		l.Info("req",
			zap.String("m", c.Request.Method),
			zap.String("p", c.FullPath()),
			zap.Int("s", c.Writer.Status()),
			zap.Duration("lat", time.Since(st)),
		)
	}
}
func metmw() gin.HandlerFunc {
	return func(c *gin.Context) {
		st := time.Now()
		c.Next()
		req.WithLabelValues(c.Request.Method, c.FullPath(), strconv.Itoa(c.Writer.Status())).Inc()
		dur.WithLabelValues(c.Request.Method, c.FullPath()).Observe(time.Since(st).Seconds())
	}
}
