package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	str := "sawa"

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello, %s", str),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
