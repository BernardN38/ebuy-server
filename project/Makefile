PRODUCT_BINARY=productApp
AUTH_BINARY=authApp
## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_product build_auth
	@echo "Stopping docker images (if running...)"
	docker-compose down 
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_product: builds the product service binary as a linux executable
build_auth:
	@echo "Building auth service binary..."
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"

## build_product: builds the product service binary as a linux executable
build_product:
	@echo "Building product service binary..."
	cd ../product-service && env GOOS=linux CGO_ENABLED=0 go build -o ${PRODUCT_BINARY} ./cmd/api
	@echo "Done!"
