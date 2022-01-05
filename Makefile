PRINT= printf "%b"
BLUE=\033[0;94m
NC=\033[0m

build:
	go build -o ./target/product-api ./cmd/api

.PHONY: build