package handler

import (
	"net/http"
	"sync"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"product-management-backend/db"
	"product-management-backend/routes"
)

var (
	echoOnce    sync.Once
	echoHandler http.Handler
)

// Handler is the Vercel serverless entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
	echoOnce.Do(func() {
		_ = godotenv.Load()
		_ = db.ConnectMongo()
		e := echo.New()
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
		e.GET("/swagger/*", echoSwagger.WrapHandler)
		e.GET("/", echoSwagger.WrapHandler)
		routes.RegisterRoutes(e)
		echoHandler = e.Server.Handler
	})
	echoHandler.ServeHTTP(w, r)
}
