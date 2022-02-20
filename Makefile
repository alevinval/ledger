cover:
	go test ./... -count=1 -cover -coverprofile coverage
	go tool cover -html coverage

compile-protos:
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		pkg/proto/*.proto

run-server:
	go run ./cmd/remote/server

run-writer:
	go run ./cmd/remote/writer

run-reader:
	go run ./cmd/remote/reader

.PHONY: cover compile-protos run-server run-writer run-reader
