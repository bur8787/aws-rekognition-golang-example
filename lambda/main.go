package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	rekognitiontypes "github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"log"
	"os"
)

type RekognitionRequest struct {
	ImageKey string `json:"key"`
}

type RekognitionResponse struct {
	Labels []string `json:"labels"`
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var rekognitionRequest RekognitionRequest
	err := json.Unmarshal([]byte(request.Body), &rekognitionRequest)
	if err != nil {
		log.Printf("Error parsing request: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid request"}, nil
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to load AWS SDK config: %v", err)
	}

	bucketName, ok := os.LookupEnv("S3_BUCKET")
	fmt.Printf("bucketName: %v\n", bucketName)
	if !ok {
		log.Fatalf("S3_BUCKET environment variable not set")
	}

	client := rekognition.NewFromConfig(cfg)
	input := &rekognition.DetectLabelsInput{
		Image: &rekognitiontypes.Image{
			S3Object: &rekognitiontypes.S3Object{
				Bucket: aws.String(bucketName),
				Name:   &rekognitionRequest.ImageKey,
			},
		},
		MaxLabels:     aws.Int32(10),
		MinConfidence: aws.Float32(75),
	}

	result, err := client.DetectLabels(ctx, input)
	if err != nil {
		log.Printf("Rekognition API error: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error processing image"}, nil
	}

	labels := []string{}
	for _, label := range result.Labels {
		labels = append(labels, *label.Name)
	}

	response := RekognitionResponse{Labels: labels}
	responseBody, _ := json.Marshal(response)
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(responseBody)}, nil
}

func main() {
	lambda.Start(handleRequest)
}
