package main

import (
	_ "github.com/Bieroid/dictionary/docs"
	"github.com/Bieroid/dictionary/comand"
	"github.com/Bieroid/dictionary/repository"
	"github.com/Bieroid/dictionary/service"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var (
	// Создаем счетчик для отслеживания количества запросов
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path"},
	)

	// Создаем гистограмму для отслеживания времени обработки запросов
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of latencies for HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {
	// Регистрация метрик в Prometheus
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestDuration)
}

func main() {
	e := echo.New()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Начинаем таймер
			start := time.Now()

			// Обрабатываем запрос
			err := next(c)

			// Записываем метрики
			duration := time.Since(start).Seconds()
			requestCounter.WithLabelValues(c.Request().Method, c.Request().URL.Path).Inc()
			requestDuration.WithLabelValues(c.Request().Method, c.Request().URL.Path).Observe(duration)

			return err
		}
	})

	connStr := "host=localhost port=5432 user=postgres password=123 dbname=postgres sslmode=disable"

	dictRepo := repository.NewRepository("postgres", connStr)
	dictIntegrationService := service.NewOuterService("http", "localhost:8800", "lite", 10)
	dictTranslateService := service.NewTranslateService(dictRepo, dictIntegrationService)
	dictHandle := comand.NewTranslateHandle(dictTranslateService)

	e.GET("/dictionary/api/translate", dictHandle.TranslateWord)
	e.POST("/dictionary/api/words", dictHandle.AddTranslation)
	e.POST("/dictionary/api/languages", dictHandle.AddLanguage)
	e.GET("/dictionary/api/languages", dictHandle.GetLanguages)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.Start(":8080")
}
