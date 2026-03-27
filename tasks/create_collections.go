package tasks

import (
	"context"
	"fmt"
	"os"
	"product-management-backend/db"
	"strings"
)

func getDBNameFromURI() (string, error) {
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

// CreateCollections creates required collections in the MongoDB database
func CreateCollections() error {
	dbName, err := getDBNameFromURI()
	if err != nil {
		return err
	}
	ctx := context.Background()
	database := db.Client.Database(dbName)
	collections := []string{"users", "products", "user_product_changes"}
	for _, col := range collections {
		if err := database.CreateCollection(ctx, col); err != nil && !isCollectionExistsError(err) {
			return fmt.Errorf("failed to create collection %s: %w", col, err)
		}
	}
	return nil
}

// isCollectionExistsError checks if the error is due to collection already existing
func isCollectionExistsError(err error) bool {
	return err != nil && (err.Error() == "(NamespaceExists) a collection 'users' already exists" || err.Error() == "(NamespaceExists) a collection 'products' already exists" || err.Error() == "(NamespaceExists) a collection 'user_product_changes' already exists")
}
