package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"google.golang.org/grpc"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/golang/protobuf/proto"
	"github.com/lechuckroh/aws-lambda-grpc-proxy-go/examples/grpc-client/go/hello"
)

const (
	defaultName        = "world"
	grpcService        = "lechuckroh.service.hello.Hello"
	grpcMethod         = "Call"
	lambdaFunctionName = "grpc-proxy-go"
	lambdaRegion       = "ap-northeast-2"
)

type LambdaRequest struct {
	Service string `json:"service"`
	Method  string `json:"method"`
	Data    string `json:"data"`
}

type LambdaResponse struct {
	StatusCode int16  `json:"statusCode"`
	Data       string `json:"data"`
}

func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func decodeBase64(b64Str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(b64Str)
}

// serializeMessage serializes protobuf message
func serializeMessage(m proto.Message) ([]byte, error) {
	return proto.Marshal(m)
}

// deserializeMessage deserializes protobuf message from byte array
func deserializeMessage(data []byte, m proto.Message) error {
	return proto.Unmarshal(data, m)
}

func callLambda(reqBytes []byte) []byte {
	b64Req := encodeBase64(reqBytes)

	// create lambda session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	client := lambda.New(sess, &aws.Config{Region: aws.String(lambdaRegion)})

	// create request payload
	request := LambdaRequest{grpcService, grpcMethod, b64Req}
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

	respBytes, err := decodeBase64(resp.Data)
	if err != nil {
		log.Fatal("failed to decode lambda response: ", err)
	}
	return respBytes
}

func callGrpcServer(serverAddr string, reqBytes []byte) []byte {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(BytesCodec{})),
	}
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatal("failed to connect to gRPC server: ", err)
	}
	defer func(){
		_ = conn.Close()
	}()

	var resp BytesCodecResponse
	method := fmt.Sprintf("/%s/%s", grpcService, grpcMethod)
	if err := conn.Invoke(context.Background(), method, reqBytes, &resp); err != nil {
		log.Fatal("failed to invoke gRPC method: ", err)
	}

	return resp.Data
}

func main() {
	argsCount := len(os.Args)

	name := defaultName
	if argsCount > 1 {
		name = os.Args[1]
	}
	grpcServer := ""
	if argsCount > 2 {
		grpcServer = os.Args[2]
	}

	// serialize message
	reqBytes, err := serializeMessage(&hello.CallRequest{Name: name})
	if err != nil {
		log.Fatal("fail to serialize message: ", err)
	}

	var respBytes []byte
	if grpcServer == "" {
		// call lambda
		respBytes = callLambda(reqBytes)
	} else {
		// call gRPC server
		respBytes = callGrpcServer(grpcServer, reqBytes)
	}

	// decode protobuf message
	var callResp hello.CallResponse
	if err := deserializeMessage(respBytes, &callResp); err != nil {
		log.Fatal("failed to deserialize protobuf response: ", err)
	} else {
		log.Printf("protobuf response: %+v", callResp)
	}
}
