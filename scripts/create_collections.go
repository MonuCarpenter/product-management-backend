package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
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

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	dbName := "pms" // Set your DB name here
	database := client.Database(dbName)
	collections := []string{"users", "products", "user_product_changes"}
	ctx := context.Background()
	for _, col := range collections {
		err := database.CreateCollection(ctx, col)
		if err != nil && !isCollectionExistsError(err) {
			log.Fatalf("failed to create collection %s: %v", col, err)
		}
	}
	fmt.Println("Collections created successfully.")
}

func isCollectionExistsError(err error) bool {
	return err != nil && (err.Error() == "(NamespaceExists) a collection 'users' already exists" || err.Error() == "(NamespaceExists) a collection 'products' already exists" || err.Error() == "(NamespaceExists) a collection 'user_product_changes' already exists")
}
