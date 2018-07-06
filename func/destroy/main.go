package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]

	// session
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	svc := dynamodb.New(sess)
	params := &dynamodb.DeleteItemInput{
		TableName: aws.String(os.Getenv("IMAGE_META_DATA_TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				N: aws.String(id),
			},
		},

		ReturnConsumedCapacity:      aws.String("NONE"),
		ReturnItemCollectionMetrics: aws.String("NONE"),
		ReturnValues:                aws.String("NONE"),
	}

	_, err = svc.DeleteItem(params)
	if err != nil {
		panic(fmt.Sprintf("error:%#v", err))
	}
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Delete Success"),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
