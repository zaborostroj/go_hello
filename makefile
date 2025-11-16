BIN_DIR=bin

build-orders:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/orders-service ./cmd/orders-service

build-warehouse:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/warehouse-service ./cmd/warehouse-service

build: build-orders build-warehouse

test:
	go test -v ./...

test-orders:
	go test -v ./internal/orders/...

clean:
	rm -rf $(BIN_DIR)
