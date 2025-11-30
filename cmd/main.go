package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pelicanch1k/link-checker/internal/adapter/pdf"
	"github.com/pelicanch1k/link-checker/internal/adapter/repository"
	"github.com/pelicanch1k/link-checker/internal/checker"
	httpController "github.com/pelicanch1k/link-checker/internal/controller/http"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Инициализация слоев
	repo := repository.NewInMemoryTaskRepository()
	pdfGen := pdf.NewPDFGenerator()
	useCase := checker.NewLinkCheckerUseCase(5*time.Second, 10, repo, pdfGen)
	controller := httpController.NewHTTPController(useCase)

	// Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Link Checker Service",
	})

	// Routes
	httpController.SetupRoutes(app, controller)

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down gracefully...")
		app.Shutdown()
	}()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}