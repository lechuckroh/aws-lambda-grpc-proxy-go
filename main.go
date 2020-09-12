package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

// LambdaRequest is incoming event type
type LambdaRequest struct {
	Service  string `json:"service"`
	Method   string `json:"method"`
	Message  string `json:"message"`
	Metadata string `json:"metadata"`
}

// LambdaResponse is response type
type LambdaResponse struct {
	Message string `json:"message"`
}

// init function is executed when handler is loaded.
func init() {
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

	return &LambdaResponse{Message: "not implemented"}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
