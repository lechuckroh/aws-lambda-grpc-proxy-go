.PHONY: proto
proto:
	protoc \
		--go_out=plugins=grpc:../grpc-client/go/hello \
		hello.proto
	protoc \
		--go_out=plugins=grpc:../grpc-server/go/hello \
		hello.proto
