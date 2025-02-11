import * as crypto from 'crypto';
import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as apigateway from 'aws-cdk-lib/aws-apigateway';

export class RekognitionCdkStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const randomSuffix = crypto.randomBytes(3).toString('hex');

    const bucket = new s3.Bucket(this, 'RekognitionImageBucket', {
      bucketName: 'rekognition-image-bucket-' + randomSuffix,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      autoDeleteObjects: true,
    });

    const lambdaRole = new iam.Role(this, 'LambdaExecutionRole', {
      assumedBy: new iam.ServicePrincipal('lambda.amazonaws.com'),
      managedPolicies: [
        iam.ManagedPolicy.fromAwsManagedPolicyName('AmazonRekognitionFullAccess'),
        iam.ManagedPolicy.fromAwsManagedPolicyName('AmazonS3ReadOnlyAccess'),
        iam.ManagedPolicy.fromAwsManagedPolicyName('service-role/AWSLambdaBasicExecutionRole'),
      ],
    });

    const rekognitionLambda = new lambda.Function(this, 'RekognitionLambda', {
      runtime: lambda.Runtime.PROVIDED_AL2,
      code: lambda.Code.fromAsset('lambda'),
      handler: 'bootstrap',
      role: lambdaRole,
      timeout: cdk.Duration.seconds(30),
    });

    const api = new apigateway.RestApi(this, 'RekognitionAPI', {});
    const resource = api.root.addResource('analyze');
    resource.addMethod('POST', new apigateway.LambdaIntegration(rekognitionLambda));

    rekognitionLambda.addEnvironment('S3_BUCKET', bucket.bucketName);

    new cdk.CfnOutput(this, 'S3BucketName', {
      value: bucket.bucketName,
      description: 'This is the S3 bucket where images should be uploaded.',
    });
  }
}
