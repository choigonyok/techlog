all: run unit integration build
run:
	go run main.go
format:
	gofmt -l .
unit:
	go test ./...
integration:
	go test ./integration/transformer_test.go
build:
	docker build . -t shadowshotx/product-go-micro