package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"product-management-backend/models"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	_ = godotenv.Load()
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI not set in .env")
	}
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	dbName := "pms"
	usersCol := client.Database(dbName).Collection("users")
	ctx := context.Background()

	adminEmail := "admin@pms.com"
	adminPassword := "pms@admin"

	var existing models.User
	err = usersCol.FindOne(ctx, map[string]interface{}{"email": adminEmail}).Decode(&existing)
	if err == nil {
		fmt.Println("Admin user already exists.")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	admin := models.User{
		ID:        primitive.NewObjectID(),
		Name:      "Admin",
		Email:     adminEmail,
		Phone:     "",
		Password:  string(hash),
		Role:      models.RoleAdmin,
		DeletedAt: nil,
	}
	_, err = usersCol.InsertOne(ctx, admin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Admin user created successfully.")
}
