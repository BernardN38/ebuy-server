PRODUCT_BINARY=productApp
AUTH_BINARY=authApp
ORDER_BINARY=orderApp
## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_product build_auth build_image_proccessing_service build_media build_order_service build_image_proccessing_service
	@echo "Stopping docker images (if running...) test change"
	docker-compose down product-service authentication-service media-service image-proccessing-service nginx 
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build product-service authentication-service media-service nginx image-proccessing-service order-service
	@echo "Docker images built and started!"

up_build_all_d: build_product build_auth build_image_proccessing_service build_media build_order_service build_image_proccessing_service
	@echo "Stopping docker images (if running...)"
	docker-compose down 
	@echo "Building (when required) and starting docker images..."
	docker-compose up -d  --build product-service authentication-service media-service nginx order-service image-proccessing-service rabbitmq postgres minio
	@echo "Docker images built and started!"

up_build_d: build_product build_auth build_image_proccessing_service build_media
	@echo "Stopping docker images (if running...)"
	docker-compose down product-service authentication-service media-service image-proccessing-service nginx 
	@echo "Building (when required) and starting docker images..."
	docker-compose up -d --build product-service authentication-service media-service image-proccessing-service nginx
	@echo "Docker images built and started!"

up_infra:
	@echo "Stopping docker images (if running...)"
	docker-compose down rabbitmq postgres minio
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d rabbitmq postgres minio
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

## build_media: builds the media service binary as a linux executable
build_media:
	@echo "Building media service binary..."
	cd ../media-service && ./mvnw install -DskipTests -Dmaven.compiler.useIncrementalCompilation=true
	@echo "Done!"

build_order_service:
	@echo "Building order service binary..."
	cd ../order-service && env GOOS=linux CGO_ENABLED=0 go build -o ${ORDER_BINARY} ./cmd/api
	@echo "Done!"

build_image_proccessing_service:
	@echo "Building image_proccessing binary..."
	cd ../image-proccessing-service && docker build -t image-proccessing-service -f image-proccessing-service.dockerfile .
	@echo "Done!"