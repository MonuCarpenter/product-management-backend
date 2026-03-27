package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"

	"product-management-backend/db"
	"product-management-backend/routes"
)

func main() {
	_ = godotenv.Load()
	if err := db.ConnectMongo(); err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register API routes
	// Swagger docs endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Register all API routes
	routes.RegisterRoutes(e)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on port", port)
	e.Logger.Fatal(e.Start(":" + port))
}
