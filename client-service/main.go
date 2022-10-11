package main

import (
	"math/rand"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var AccessCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "api_requests_total",
	},
	[]string{"method", "path"},
)

var QueueGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "queue_num_total",
	},
	[]string{"name"},
)

var HttpDurationsHistogram = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_durations_histogram_seconds",
		Buckets: []float64{0.2, 0.5, 1, 2, 5, 10, 30},
	},
	[]string{"path"},
)

var HttpDurations = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name:       "http_durations_seconds",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"path"},
)

func init() {
	prometheus.MustRegister(AccessCounter)
	prometheus.MustRegister(QueueGauge)
	prometheus.MustRegister(HttpDurationsHistogram)
	prometheus.MustRegister(HttpDurations)
}

func main() {

	// init gin
	engine := gin.New()

	// test ep
	engine.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "world")
	})

	/*
		test metrics
	*/

	// counter metric
	engine.GET("/counter", func(c *gin.Context) {
		purl, _ := url.Parse(c.Request.RequestURI)
		AccessCounter.With(prometheus.Labels{
			"method": c.Request.Method,
			"path":   purl.Path,
		}).Add(1)
	})

	// gauge metric
	engine.GET("/gauge", func(c *gin.Context) {
		num := c.Query("num")
		fnum, _ := strconv.ParseFloat(num, 32)
		QueueGauge.With(prometheus.Labels{"name": "queue_eddycjy"}).Set(fnum)
	})

	// histogram metric
	engine.GET("/histogram", func(c *gin.Context) {
		purl, _ := url.Parse(c.Request.RequestURI)
		HttpDurationsHistogram.With(prometheus.Labels{"path": purl.Path}).Observe(float64(rand.Intn(30)))
	})

	// summary metric
	engine.GET("/summary", func(c *gin.Context) {
		purl, _ := url.Parse(c.Request.RequestURI)
		HttpDurations.With(prometheus.Labels{"path": purl.Path}).Observe(float64(rand.Intn(30)))
	})

	// display prometheus metrics
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	engine.Run(":10001")
}
