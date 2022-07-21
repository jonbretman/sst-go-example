package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	rds "github.com/aws/aws-sdk-go/service/rdsdataservice"
)

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	s, err := session.NewSession()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}, nil
	}

	client := rds.New(s)

	result, err := client.ExecuteStatement(&rds.ExecuteStatementInput{
		SecretArn:       aws.String(os.Getenv("DATABASE_SECRET_ARN")),
		Database:        aws.String(os.Getenv("DATABASE_NAME")),
		ResourceArn:     aws.String(os.Getenv("DATABASE_RESOURCE_ARN")),
		Sql:             aws.String("select * from people"),
		FormatRecordsAs: aws.String("JSON"),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}, nil
	}

	if result.FormattedRecords == nil {
		return events.APIGatewayProxyResponse{
			Body:       "No results from database",
			StatusCode: 400,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       *result.FormattedRecords,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
