.PHONY=compile-protos

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
