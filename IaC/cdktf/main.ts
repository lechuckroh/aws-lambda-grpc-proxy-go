import {Construct} from 'constructs';
import {App, TerraformStack} from 'cdktf';
import {
    AwsProvider,
    CloudwatchLogGroup,
    IamPolicy,
    IamRole,
    IamRolePolicyAttachment,
    LambdaFunction
} from './.gen/providers/aws';

class LambdaGrpcProxyStack extends TerraformStack {
    constructor(scope: Construct, name: string) {
        super(scope, name);

        // AWS Provider
        new AwsProvider(this, 'aws', {
            region: 'ap-northeast-2',
        });

        // variable
        // TODO: TerraformVariable is not implemented yet.
        const lambdaFunctionName = 'grpc-proxy-go';
        const handler = 'app';
        const grpcServerAddr = 'grpc.lechuckcgx.com:59090';
        const memory = 128;
        const filename = `${process.cwd()}/../../app.zip`;

        // #####################
        // Role
        // #####################

        const lambdaExecRole = new IamRole(this, 'lambdaExecRole', {
            name: "lambda_execution",
            description: "Allows Lambda functions to call AWS services on your behalf",
            assumeRolePolicy: JSON.stringify({
                Version: "2012-10-17",
                Statement: [{
                    Effect: "Allow",
                    Principal: {
                        Service: "lambda.amazonaws.com"
                    },
                    Action: "sts:AssumeRole",
                }]
            }),
            tags: {
                Terraform: 'true',
            }
        });

        // #####################
        // Lambda Function
        // #####################

        new LambdaFunction(this, "proxyLambda", {
            role: lambdaExecRole.arn,
            functionName: lambdaFunctionName,
            handler: handler,
            runtime: "go1.x",
            memorySize: memory,
            filename: filename,
            environment: [{
                variables: {
                    'SERVER_ADDR': grpcServerAddr
                }
            }],
        });

        // #####################
        // Cloudwatch Log Group
        // #####################

        new CloudwatchLogGroup(this, "logGroup", {
            name: `/aws/lambda/${lambdaFunctionName}`,
            retentionInDays: 14,
            tags: {
                Terraform: 'true'
            }
        });

        const loggingPolicy = new IamPolicy(this, 'loggingPolocy', {
            name: 'lambda_logging',
            path: '/',
            description: 'IAM policy for logging from a lambda',
            policy: JSON.stringify({
                Version: "2012-10-17",
                Statement: [{
                    Action: [
                        "logs:CreateLogGroup",
                        "logs:CreateLogStream",
                        "logs:PutLogEvents"
                    ],
                    Resource: "arn:aws:logs:*:*:*",
                    Effect: "Allow"
                }]
            }),
        });

        new IamRolePolicyAttachment(this, 'lambdaLogs', {
            role: lambdaExecRole.name!,
            policyArn: loggingPolicy.arn,
        });
    }
}

const app = new App();
new LambdaGrpcProxyStack(app, 'aws-lambda-grpc-proxy-go');
app.synth();
