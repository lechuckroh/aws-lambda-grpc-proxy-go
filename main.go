package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

// LambdaRequest is incoming event type
type LambdaRequest struct {
	Service string `json:"service"`
	Method  string `json:"method"`
	Data    string `json:"data"`
}

// LambdaResponse is response type
type LambdaResponse struct {
	StatusCode int16  `json:"statusCode"`
	Data       string `json:"data"`
}

// init function is executed when handler is loaded.
func init() {
}

// encodeBase64 returns the base64 encoding of data.
func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// decodeBase64 returns the bytes represented by the base64 string.
func decodeBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// HandleRequest is a handler function
func HandleRequest(ctx context.Context, req LambdaRequest) (*LambdaResponse, error) {
	// context
	lc, _ := lambdacontext.FromContext(ctx)
	log.Printf("AwsRequestID: %s", lc.AwsRequestID)

	// environment variables
	serverAddr := os.Getenv("SERVER_ADDR")

	log.Printf("serverAddr: %s", serverAddr)
	log.Printf("request: %+v", req)

	// connect to gRPC server
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(BytesCodec{})),
	}
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		return &LambdaResponse{500, err.Error()}, nil
	}
	defer func() {
		_ = conn.Close()
	}()

	// decode base64 encoded request data
	if reqDataBytes, err := decodeBase64(req.Data); err != nil {
		return &LambdaResponse{500, err.Error()}, nil
	} else {
		var resp BytesCodecResponse
		method := fmt.Sprintf("/%s/%s", req.Service, req.Method)

		// invoke gRPC method
		if err := conn.Invoke(ctx, method, reqDataBytes, &resp); err != nil {
			return &LambdaResponse{500, err.Error()}, nil
		}

		// base64 encode gRPC response
		b64RespData := encodeBase64(resp.Data)

		// return response
		return &LambdaResponse{200, b64RespData}, nil
	}
}

func main() {
	lambda.Start(HandleRequest)
}
