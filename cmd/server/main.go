// Package main runs the application (starts the server)
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/5fives-to-go/internal/database"
	"github.com/5fives-to-go/internal/server"
)

func main() {
	srv := server.NewHTTPServer()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Error listening and serving: %v\n", err)
		}
	}()

	fmt.Println("Testing db connection")
	db, err := database.GetDatabaseConnection()
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()
	fmt.Println("Database connection successful")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan

	ctx, stop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stop()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("\nServer shutdown with error: %v\n", err)
	} else {
		fmt.Println("\nServer shutdown gracefully")
	}
}
