package main

import (
	"context"
	"log"
	"net/http"
	"postgres-listen-notify-producer/db"
	"postgres-listen-notify-producer/handlers"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConn, err := db.ConnectDB(ctx)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer dbConn.Close()

	router := gin.Default()
	router.POST("/bet", handlers.PostBetHandler(dbConn))

	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
