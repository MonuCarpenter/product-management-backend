package tasks

import (
	"context"
	"fmt"
	"os"
	"product-management-backend/db"
	"product-management-backend/models"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func getDBNameFromURL() (string, error) {
	url := os.Getenv("MONGODB_URI")
	if url == "" {
		return "", fmt.Errorf("MONGODB_URI must be set in env")
	}
	parts := strings.Split(url, "/")
	if len(parts) < 4 {
		return "", fmt.Errorf("invalid MONGODB_URI format")
	}
	dbAndParams := parts[len(parts)-1]
	dbName := dbAndParams
	if idx := strings.Index(dbAndParams, "?"); idx != -1 {
		dbName = dbAndParams[:idx]
	}
	if dbName == "" {
		return "", fmt.Errorf("database name not found in MONGODB_URI")
	}
	return dbName, nil
}

// CreateAdminUser creates an admin user in the database if not exists
func CreateAdminUser() error {
	dbName, err := getDBNameFromURL()
	if err != nil {
		return err
	}
	ctx := context.Background()
	usersCol := db.Client.Database(dbName).Collection("users")
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminEmail == "" || adminPassword == "" {
		return fmt.Errorf("ADMIN_EMAIL and ADMIN_PASSWORD must be set in env")
	}
	var existing models.User

	err = usersCol.FindOne(ctx, bson.M{"email": adminEmail}).Decode(&existing)
	if err == nil {
		return nil // already exists
	}
	if err != mongo.ErrNoDocuments {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	admin := models.User{
		Name:     "Admin",
		Email:    adminEmail,
		Phone:    "",
		Password: string(hash),
		Role:     models.RoleAdmin,
	}
	_, err = usersCol.InsertOne(ctx, admin)
	return err
}
