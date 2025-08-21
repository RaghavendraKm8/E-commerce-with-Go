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

type User struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

var (
	db  *pgxpool.Pool
	req = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "usersvc_http_requests_total", Help: "Total HTTP requests"},
		[]string{"method", "path", "status"},
	)
	dur = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Name: "usersvc_http_duration_seconds", Help: "Request duration", Buckets: prometheus.DefBuckets},
		[]string{"method", "path"},
	)
)

func init() { prometheus.MustRegister(req, dur) }

func main() {
	lg, _ := zap.NewProduction()
	defer lg.Sync()

	pg := os.Getenv("PG_URL")
	if pg == "" {
		pg = "postgres://postgres:postgres@localhost:5432/ecom_users?sslmode=disable"
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
	r.POST("/users", create)
	r.GET("/users/:id", get)

	srv := &http.Server{Addr: ":8081", Handler: r}
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
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := db.QueryRow(c,
		"INSERT INTO users(email, age) VALUES($1,$2) RETURNING id",
		u.Email, u.Age).Scan(&u.ID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, u)
}

func get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var u User
	if err := db.QueryRow(c,
		"SELECT id, email, age FROM users WHERE id=$1", id).
		Scan(&u.ID, &u.Email, &u.Age); err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	c.JSON(200, u)
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
