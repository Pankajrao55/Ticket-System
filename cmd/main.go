package main

import (
	"log"
	"net/http"
	"os"
	"ticket-system/internal/database"
	"ticket-system/internal/handlers"
	"ticket-system/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database.Init()

	// Setup Gin router
	router := gin.Default()

	// Health check endpoint - GET and HEAD both supported
	healthHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
	router.GET("/health", healthHandler)
	router.HEAD("/health", healthHandler)

	// Auth routes - public (no JWT needed)
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", handlers.Register)
		authRoutes.POST("/login", handlers.Login)
	}

	// Ticket routes - protected (JWT required)
	ticketRoutes := router.Group("/tickets")
	ticketRoutes.Use(middleware.AuthRequired())
	{
		ticketRoutes.POST("", handlers.CreateTicket)
		ticketRoutes.GET("", handlers.ListTickets)
		ticketRoutes.GET("/:id", handlers.GetTicket)
		ticketRoutes.PATCH("/:id/status", handlers.UpdateTicketStatus)
	}

	// Port from env or default 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
