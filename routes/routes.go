package routes

import (
	"product-management-backend/controllers"
	"product-management-backend/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	// Get user by token
	e.GET("/api/user/me", controllers.GetUserByToken, middleware.JWTMiddleware)

	// Get product by id (already exists as /api/products/:id)
	
	// Auth
	e.POST("/api/auth/login", controllers.Login)
	e.POST("/api/auth/register", controllers.RegisterSalesman, middleware.JWTMiddleware, middleware.RoleMiddleware("admin"))

	// Users (admin only)
	e.GET("/api/users", controllers.GetUsers, middleware.JWTMiddleware, middleware.RoleMiddleware("admin"))
	e.GET("/api/users/:id", controllers.GetUserByID, middleware.JWTMiddleware, middleware.RoleMiddleware("admin"))
	e.DELETE("/api/users/:id", controllers.DeleteUser, middleware.JWTMiddleware, middleware.RoleMiddleware("admin"))

	// Products
	// Create product (admin only)
	e.POST("/api/products/create", controllers.CreateProduct, middleware.JWTMiddleware, middleware.RoleMiddleware("admin"))
	e.GET("/api/products", controllers.GetProducts, middleware.JWTMiddleware)
	e.GET("/api/products/:id", controllers.GetProductByID, middleware.JWTMiddleware)
	e.POST("/api/products", controllers.AddProduct, middleware.JWTMiddleware, middleware.RoleMiddleware("admin"))
	e.POST("/api/products/bulk", controllers.AddBulkProducts, middleware.JWTMiddleware, middleware.RoleMiddleware("admin"))
	e.PUT("/api/products/:id", controllers.UpdateProduct, middleware.JWTMiddleware, middleware.RoleMiddleware("admin", "salesman"))
	e.DELETE("/api/products/:id", controllers.DeleteProduct, middleware.JWTMiddleware, middleware.RoleMiddleware("admin"))
}
