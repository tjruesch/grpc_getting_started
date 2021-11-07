.PHONY: protos run client

protos:
	 protoc protos/translations.proto --go-grpc_out=protos --go_out=protos

run:
	go run main.go

client:
	 protoc protos/translations.proto --python_out=clients