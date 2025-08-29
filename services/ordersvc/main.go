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
)

type Order struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"userId"`
	ProductID int64     `json:"productId"`
	Quantity  int       `json:"quantity"`
	Total     int64     `json:"total"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

var (
	db *pgxpool.Pool

	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ordersvc_http_requests_total",
			Help: "Number of HTTP requests received",
		},
		[]string{"path", "method"},
	)
	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ordersvc_http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)
)

func main() {
	// Register metrics
	prometheus.MustRegister(httpRequests, httpDuration)

	// connect DB
	url := os.Getenv("PG_URL")
	if url == "" {
		url = "postgres://postgres:postgres@localhost:5432/ordersvc?sslmode=disable"
	}
	ctx := context.Background()
	var err error
	db, err = pgxpool.New(ctx, url)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := gin.Default()

	// middleware for metrics
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()
		httpRequests.WithLabelValues(c.FullPath(), c.Request.Method).Inc()
		httpDuration.WithLabelValues(c.FullPath(), c.Request.Method).Observe(duration)
	})

	// health check
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	// routes
	r.POST("/orders", createOrder)
	r.GET("/orders/:id", getOrder)

	// prometheus metrics
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	srv := &http.Server{Addr: ":8083", Handler: r}
	go srv.ListenAndServe()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctxShut, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctxShut)
}

// Create order handler
func createOrder(c *gin.Context) {
	var o Order
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hardcoded price simulation
	o.Total = int64(o.Quantity) * 75000

	err := db.QueryRow(context.Background(),
		"INSERT INTO orders(user_id, product_id, quantity, total) VALUES($1,$2,$3,$4) RETURNING id, created_at",
		o.UserID, o.ProductID, o.Quantity, o.Total).Scan(&o.ID, &o.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, o)
}

// Get order by ID
func getOrder(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var o Order
	err := db.QueryRow(context.Background(),
		"SELECT id, user_id, product_id, quantity, total, created_at FROM orders WHERE id=$1", id).
		Scan(&o.ID, &o.UserID, &o.ProductID, &o.Quantity, &o.Total, &o.CreatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, o)
}
