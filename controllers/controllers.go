package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"product-management-backend/auth"
	"product-management-backend/db"
	"product-management-backend/models"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	v2options "go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Helper to get DB name
func getDBName() string {
	return "pms"
}

// Auth
func Login(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	col := db.Client.Database(getDBName()).Collection("users")
	var user models.User
	err := col.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}
	token, err := auth.GenerateJWT(user.ID.Hex(), string(user.Role))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "token error"})
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func RegisterSalesman(c echo.Context) error {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	col := db.Client.Database(getDBName()).Collection("users")
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "password error"})
	}
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: string(hash),
		Role:     models.RoleSalesman,
	}
	_, err = col.InsertOne(context.Background(), user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "salesman registered"})
}

// Users
// GetUsers godoc
// @Summary List users
// @Description Get a paginated list of users (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {array} models.User
// @Failure 500 {object} map[string]string
// @Router /api/users [get]
func GetUsers(c echo.Context) error {
	col := db.Client.Database(getDBName()).Collection("users")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}
	skip := int64((page - 1) * limit)
	filter := bson.M{"deletedAt": bson.M{"$exists": false}}
	opts := v2options.Find().SetSkip(skip).SetLimit(int64(limit))
	cur, err := col.Find(context.Background(), filter, opts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
	}
	var users []models.User
	if err := cur.All(context.Background(), &users); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
	}
	return c.JSON(http.StatusOK, users)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get a user by ID (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Router /api/users/{id} [get]
func GetUserByID(c echo.Context) error {
	col := db.Client.Database(getDBName()).Collection("users")
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	var user models.User
	err = col.FindOne(context.Background(), bson.M{"_id": id, "deletedAt": bson.M{"$exists": false}}).Decode(&user)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	return c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Soft delete a user (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/users/{id} [delete]
func DeleteUser(c echo.Context) error {
	col := db.Client.Database(getDBName()).Collection("users")
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	now := time.Now()
	_, err = col.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": bson.M{"deletedAt": now}})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "user deleted"})
}

// Products
// GetProducts godoc
// @Summary List products
// @Description Get a paginated list of products
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {array} models.Product
// @Failure 500 {object} map[string]string
// @Router /api/products [get]
func GetProducts(c echo.Context) error {
	col := db.Client.Database(getDBName()).Collection("products")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}
	skip := int64((page - 1) * limit)
	filter := bson.M{"deletedAt": bson.M{"$exists": false}}
	opts := v2options.Find().SetSkip(skip).SetLimit(int64(limit))
	cur, err := col.Find(context.Background(), filter, opts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
	}
	var products []models.Product
	if err := cur.All(context.Background(), &products); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
	}
	return c.JSON(http.StatusOK, products)
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Failure 404 {object} map[string]string
// @Router /api/products/{id} [get]
func GetProductByID(c echo.Context) error {
	col := db.Client.Database(getDBName()).Collection("products")
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	var product models.Product
	err = col.FindOne(context.Background(), bson.M{"_id": id, "deletedAt": bson.M{"$exists": false}}).Decode(&product)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	return c.JSON(http.StatusOK, product)
}

func AddProduct(c echo.Context) error {
	col := db.Client.Database(getDBName()).Collection("products")
	var req models.Product
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	req.DeletedAt = nil
	_, err := col.InsertOne(context.Background(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "product added"})
}

func AddBulkProducts(c echo.Context) error {
	col := db.Client.Database(getDBName()).Collection("products")
	var req []models.Product
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	var docs []interface{}
	for i := range req {
		req[i].DeletedAt = nil
		docs = append(docs, req[i])
	}
	_, err := col.InsertMany(context.Background(), docs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "bulk products added"})
}

func UpdateProduct(c echo.Context) error {
	col := db.Client.Database(getDBName()).Collection("products")
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	var req models.Product
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	_, err = col.UpdateOne(context.Background(), bson.M{"_id": id, "deletedAt": bson.M{"$exists": false}}, bson.M{"$set": req})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "product updated"})
}

func DeleteProduct(c echo.Context) error {
	col := db.Client.Database(getDBName()).Collection("products")
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	now := time.Now()
	_, err = col.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": bson.M{"deletedAt": now}})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "product deleted"})
}

// User-Product Changes
func GetChanges(c echo.Context) error {
	col := db.Client.Database(getDBName()).Collection("user_product_changes")
	cur, err := col.Find(context.Background(), bson.M{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
	}
	var changes []models.UserProductChange
	if err := cur.All(context.Background(), &changes); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
	}
	return c.JSON(http.StatusOK, changes)
}

func GetChangeByID(c echo.Context) error {
	col := db.Client.Database(getDBName()).Collection("user_product_changes")
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	var change models.UserProductChange
	err = col.FindOne(context.Background(), bson.M{"_id": id}).Decode(&change)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	return c.JSON(http.StatusOK, change)
}
