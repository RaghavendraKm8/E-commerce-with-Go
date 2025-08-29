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

type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

var (
	db *pgxpool.Pool

	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "productsvc_http_requests_total",
			Help: "Number of HTTP requests received",
		},
		[]string{"path", "method"},
	)
	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "productsvc_http_request_duration_seconds",
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
		url = "postgres://postgres:postgres@localhost:5432/productsvc?sslmode=disable"
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
	r.POST("/products", createProduct)
	r.GET("/products", listProducts)
	r.GET("/products/:id", getProduct)

	// prometheus metrics
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	srv := &http.Server{Addr: ":8082", Handler: r}
	go srv.ListenAndServe()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctxShut, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctxShut)
}

// Create product
func createProduct(c *gin.Context) {
	var p Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.QueryRow(context.Background(),
		"INSERT INTO products(name, price) VALUES($1,$2) RETURNING id",
		p.Name, p.Price).Scan(&p.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

// List all products
func listProducts(c *gin.Context) {
	rows, err := db.Query(context.Background(),
		"SELECT id, name, price FROM products ORDER BY id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err == nil {
			products = append(products, p)
		}
	}
	c.JSON(http.StatusOK, products)
}

// Get product by ID
func getProduct(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var p Product
	err := db.QueryRow(context.Background(),
		"SELECT id, name, price FROM products WHERE id=$1", id).
		Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}
