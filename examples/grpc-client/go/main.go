package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/golang/protobuf/proto"
	"github.com/lechuckroh/aws-lambda-grpc-proxy-go/examples/grpc-client/go/hello"
)

const (
	defaultName        = "world"
	grpcService        = "lechuckroh.service.hello"
	grpcMethod         = "Call"
	lambdaFunctionName = "grpc-proxy-go"
	lambdaRegion       = "ap-northeast-2"
)

type LambdaRequest struct {
	Service string `json:"service"`
	Method  string `json:"method"`
	Message string `json:"message"`
}

type LambdaResponse struct {
	StatusCode int16  `json:"statusCode"`
	Message    string `json:"message"`
}

// serializeMessage serializes protobuf message and returns base64 encoded string.
func serializeMessage(m proto.Message) (string, error) {
	if data, err := proto.Marshal(m); err != nil {
		return "", err
	} else {
		return base64.StdEncoding.EncodeToString(data), nil
	}
}

func deserializeMessage(b64Str string, m proto.Message) error {
	if data, err := base64.StdEncoding.DecodeString(b64Str); err != nil {
		return err
	} else {
		return proto.Unmarshal(data, m)
	}
}

func main() {
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	// serialize message
	b64Data, err := serializeMessage(&hello.CallRequest{Name: name})
	if err != nil {
		log.Fatal("fail to serialize message: ", err)
	}
	log.Printf("serialized request: %+v", b64Data)

	// create lambda session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	client := lambda.New(sess, &aws.Config{Region: aws.String(lambdaRegion)})

	// create request payload
	request := LambdaRequest{grpcService, grpcMethod, b64Data}
	payload, err := json.Marshal(request)
	if err != nil {
		log.Fatal("failed to marshal lambda request: ", err)
	}

	// invoke lambda
	result, err := client.Invoke(&lambda.InvokeInput{
		FunctionName: aws.String(lambdaFunctionName),
		Payload:      payload,
	})
	if err != nil {
		log.Fatal("failed to invoke lambda function: ", err)
	}

	// parse response
	var resp LambdaResponse
	if err := json.Unmarshal(result.Payload, &resp); err != nil {
		log.Fatal("failed to unmarshal lambda response: ", err)
	}
	log.Printf("lambda response: %+v", resp)

	// decode protobuf message
	var callResp hello.CallResponse
	if  err := deserializeMessage(resp.Message, &callResp); err != nil {
		log.Fatal("failed to deserialize protobuf response: ", err)
	} else {
		log.Printf("protobuf response: %+v", callResp)
	}
}
