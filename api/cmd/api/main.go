package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Reazy-ai/incident-tracker/internal/database"
	"github.com/Reazy-ai/incident-tracker/internal/handlers"
	"github.com/Reazy-ai/incident-tracker/internal/metrics"
	"github.com/Reazy-ai/incident-tracker/internal/middleware"
	"github.com/Reazy-ai/incident-tracker/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func main() {
	metrics.Register()

	_ = godotenv.Load()

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DB_URL")

	database.RunMigrations(dbURL)

	db, err := database.NewPostgresConnection()
	if err != nil {
		panic(err)
	}

	defer db.Close()
	incidentRepo := repositories.NewIncidentRepository(db)

	incidentHandler := handlers.NewIncidentHandler(incidentRepo)

	router := gin.Default()

	router.Use(middleware.MetricsMiddleware())

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.POST("/incidents", incidentHandler.CreateIncident)
	router.GET("/incidents", incidentHandler.GetIncidents)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	router.GET("/ready", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := db.Ping(ctx)

		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "database unavailable",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {

		log.Info().
			Str("port", port).
			Msg("starting API server")

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			log.Fatal().
				Err(err).
				Msg("failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit

	log.Info().Msg("shutdown signal received")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {

		log.Error().
			Err(err).
			Msg("graceful shutdown failed")
	}

	db.Close()

	log.Info().Msg("database connection closed")
	log.Info().Msg("server exited cleanly")
}
