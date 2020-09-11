package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

// MyEvent is incoming event type
type MyEvent struct {
	Key1 string `json:"key1"`
	Key2 string `json:"key2"`
	Key3 string `json:"key3"`
}

// MyResponse is response type
type MyResponse struct {
	Message string `json:"message"`
}

// init function is executed when handler is loaded.
func init() {
}

// HandleRequest is a handler function
func HandleRequest(ctx context.Context, evt MyEvent) (*MyResponse, error) {
	// context
	lc, _ := lambdacontext.FromContext(ctx)
	log.Printf("AwsRequestID: %s", lc.AwsRequestID)

	// environment variables
	for _, e := range os.Environ() {
		log.Println(e)
	}

	log.Printf("Key1: %s", evt.Key1)
	log.Printf("Key2: %s", evt.Key2)
	log.Printf("Key3: %s", evt.Key3)

	if evt.Key3 == "" {
		return nil, errors.New("key3 is empty")
	}
	return &MyResponse{Message: evt.Key1}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
