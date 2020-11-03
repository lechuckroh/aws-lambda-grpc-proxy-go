#!/usr/bin/env node
import * as cdk from '@aws-cdk/core';
import { GrpcProxyLambdaStack } from '../lib/grpc-proxy-lambda-stack';

const app = new cdk.App();
new GrpcProxyLambdaStack(app, 'GrpcProxyLambdaStack');
