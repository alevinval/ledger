.PHONY=compile-protos

compile-protos:
	protoc -I=proto/ --go_out=plugins=grpc:. proto/*.proto

