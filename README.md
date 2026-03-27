# Install dependencies

GO111MODULE=on go mod tidy

# Generate Swagger docs

swag init --parseDependency --parseInternal

# Run the server

PORT=8080 go run main.go
