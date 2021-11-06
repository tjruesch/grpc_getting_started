.PHONY: protos

protos:
	 protoc -I protos/ protos/translations.proto --go-grpc_out=protos --go_out=protos