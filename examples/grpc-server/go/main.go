package main

import (
	"context"
	"fmt"
	"github.com/lechuckroh/aws-lambda-grpc-proxy-go/examples/grpc-server/go/hello"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":9090"
)

type HelloServer struct {
}

func (s *HelloServer) Call(ctx context.Context, req *hello.CallRequest) (*hello.CallResponse, error) {
	resp := hello.CallResponse{
		Msg: fmt.Sprintf("hello %s", req.Name),
	}
	log.Printf("request: %+v", req)
	log.Printf("response: %+v", resp)
	return &resp, nil
}

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", port)
	}
	log.Printf("listening on %v", port)

	server := grpc.NewServer()
	hello.RegisterHelloServer(server, &HelloServer{})
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %+v", err)
	}
}
