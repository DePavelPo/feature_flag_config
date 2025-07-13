include .env

.PHONY: run
run:
	go run cmd/main.go

.PHONY: proto-gen
proto-gen:
	protoc -I=./api \
	--go_out=./pkg/pb \
	--go-grpc_out=./pkg/pb \
	--grpc-gateway_out=./pkg/pb \
	featureFlagConfig/feature_flag_config.proto