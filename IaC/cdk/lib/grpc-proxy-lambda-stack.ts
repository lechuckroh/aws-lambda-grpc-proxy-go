import * as iam from "@aws-cdk/aws-iam";
import { Effect } from "@aws-cdk/aws-iam";
import * as cdk from "@aws-cdk/core";
import { RemovalPolicy } from "@aws-cdk/core";
import * as lambda from "@aws-cdk/aws-lambda";
import * as logs from "@aws-cdk/aws-logs";

export class GrpcProxyLambdaStack extends cdk.Stack {
  constructor(scope: cdk.App, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const lambdaFunctionName = "grpc-proxy-go";
    const handler = "app";
    const grpcServerAddr = "grpc.lechuckcgx.com:59090";
    const memory = 128;
    const filename = `${process.cwd()}/../../app.zip`;

    // #####################
    // Role
    // #####################

    const lambdaExecRole = new iam.Role(this, "lambdaExecRole", {
      roleName: "lambda_execution",
      description:
        "Allows Lambda functions to call AWS services on your behalf",
      assumedBy: new iam.ServicePrincipal("lambda.amazonaws.com"),
    });

    // #####################
    // Lambda Function
    // #####################

    new lambda.Function(this, "proxyLambda", {
      role: lambdaExecRole,
      functionName: lambdaFunctionName,
      handler: handler,
      runtime: lambda.Runtime.GO_1_X,
      memorySize: memory,
      code: new lambda.AssetCode(filename),
      environment: {
        SERVER_ADDR: grpcServerAddr,
      },
      architectures: [lambda.Architecture.X86_64],
    });

    // #####################
    // Cloudwatch Log Group
    // #####################

    new logs.LogGroup(this, "logGroup", {
      logGroupName: `/aws/lambda/${lambdaFunctionName}`,
      retention: 14,
      removalPolicy: RemovalPolicy.DESTROY,
    });

    lambdaExecRole.addToPolicy(
      new iam.PolicyStatement({
        resources: ["arn:aws:logs:*:*:*"],
        actions: [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
        ],
        effect: Effect.ALLOW,
      })
    );
  }
}
