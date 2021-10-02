############
# Role
############

resource "aws_iam_role" "lambda_exec" {
  name               = "lambda_execution"
  description        = "Allows Lambda functions to call AWS services on your behalf"
  assume_role_policy = jsonencode({
    Version   = "2012-10-17",
    Statement = [{
      Effect    = "Allow",
      Principal = {
        Service = "lambda.amazonaws.com"
      },
      Action    = "sts:AssumeRole"
    }]
  })

  tags = {
    Terraform = true
  }
}

##################
# Lambda Function
##################

resource "aws_lambda_function" "proxy_lambda" {
  role          = aws_iam_role.lambda_exec.arn
  function_name = var.lambda_function_name
  handler       = var.handler
  runtime       = "go1.x"
  memory_size   = var.memory
  filename      = var.filename
  architectures = [var.architecture]

  environment {
    variables = {
      "SERVER_ADDR" = var.grpc_server_addr
    }
  }

  # without this, you'll get the following error:
  # - Lambda was unable to decrypt the environment variables because KMS access was denied.
  depends_on = [
    aws_iam_role_policy_attachment.lambda_logs,
    aws_cloudwatch_log_group.proxy_lambda,
  ]

  tags = {
    Terraform = true
  }
}

#######################
# Cloudwatch Log Group
#######################

resource "aws_cloudwatch_log_group" "proxy_lambda" {
  name              = "/aws/lambda/${var.lambda_function_name}"
  retention_in_days = 14

  tags = {
    Terraform = true
  }
}

resource "aws_iam_policy" "lambda_logging" {
  name        = "lambda_logging"
  path        = "/"
  description = "IAM policy for logging from a lambda"

  policy = jsonencode({
    Version   = "2012-10-17",
    Statement = [{
      Action   = [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      Resource = "arn:aws:logs:*:*:*",
      Effect   = "Allow"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_logs" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = aws_iam_policy.lambda_logging.arn
}
