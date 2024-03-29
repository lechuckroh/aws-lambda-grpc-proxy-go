variable "architecture" {
  default = "x86_64"
}

variable "lambda_function_name" {
  default = "grpc-proxy-go"
}

variable "grpc_server_addr" {
  default = "grpc.lechuckcgx.com:59090"
}

variable "handler" {
  description = "the function entrypoint"
  default = "app"
}

variable "memory" {
  description = "amount of memory in MB your Lambda Function can use at runtime"
  default = 128
}

variable "executable" {
  description = "the path to the function's binary"
  default = "../../app"
}

variable "zip_filename" {
  description = "the path to the function's deployment package"
  default = "../../app.zip"
}

