# Welcome to your CDK TypeScript project

This is a blank project for CDK development with TypeScript.

The `cdk.json` file tells the CDK Toolkit how to execute your app.

## Useful commands

* `npm run build`   compile typescript to js
* `npm run watch`   watch for changes and compile
* `npm run test`    perform the jest unit tests
* `npx cdk deploy`  deploy this stack to your default AWS account/region
* `npx cdk diff`    compare deployed stack with current state
* `npx cdk synth`   emits the synthesized CloudFormation template

## Deploy

```bash
make all
```

After the deployment, you will see the API endpoint and the S3 bucket in the output:

```
Outputs:
AwsRekognitionGolangExampleStack.RekognitionAPIEndpointEDF9EFE7 = https://xxxxxxxxxx.execute-api.ap-northeast-1.amazonaws.com/prod/
AwsRekognitionGolangExampleStack.S3BucketName = rekognition-image-bucket-awsrekognitiongolangexamplestack
```

## How to use the API

Upload an image to the S3 bucket and call the API to analyze the image.

```bash
export S3_BUCKET=rekognition-image-bucket-awsrekognitiongolangexamplestack
aws s3 cp sample.jpg s3://${S3_BUCKET}/
```

Call the API to analyze the image.

```bash 
export API_ENDPOINT=https://w7pmhcu98e.execute-api.ap-northeast-1.amazonaws.com/prod/
curl -X POST ${API_ENDPOINT}analyze \
     -H "Content-Type: application/json" \
     -d '{"key": "sample.jpg"}'
```
