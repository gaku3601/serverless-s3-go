package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Data struct {
	ID       int
	FileName string
	URL      string
}

func (d *Data) createUrl() {
	host := "https://s3-ap-northeast-1.amazonaws.com/"
	bucket := os.Getenv("IMAGE_BUCKET")
	d.URL = fmt.Sprintf("%s%s/%s", host, bucket, d.FileName)
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	startID := request.QueryStringParameters["start"]
	endID := request.QueryStringParameters["end"]
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	svc := dynamodb.New(sess)
	params := &dynamodb.ScanInput{
		TableName:            aws.String(os.Getenv("IMAGE_META_DATA_TABLE")),
		ProjectionExpression: aws.String("ID, FileName"),
		FilterExpression:     aws.String("#key BETWEEN :startKey AND :endKey"),
		ExpressionAttributeNames: map[string]*string{
			"#key": aws.String("ID"), // 項目名をプレースホルダに入れる
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":startKey": {
				N: &startID,
			},
			":endKey": {
				N: &endID,
			},
		},
	}
	resp, err := svc.Scan(params)
	if err != nil {
		fmt.Errorf("failed to make Query API call, %v", err)
	}
	obj := []Data{}
	dynamodbattribute.UnmarshalListOfMaps(resp.Items, &obj)
	for i, _ := range obj {
		obj[i].createUrl()
	}
	sort.Slice(obj, func(i, j int) bool {
		return obj[i].ID > obj[j].ID
	})
	j, _ := json.Marshal(obj)

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("%v", string(j)),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
