package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func updateSequence(svc *dynamodb.DynamoDB, tableName string) *string {
	putParams := &dynamodb.UpdateItemInput{
		TableName: aws.String(os.Getenv("SEQUENCE_TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"TableName": {
				S: aws.String(tableName),
			},
		},
		AttributeUpdates: map[string]*dynamodb.AttributeValueUpdate{
			"CurrentNumber": {
				Value: &dynamodb.AttributeValue{
					N: aws.String("1"),
				},
				Action: aws.String("ADD"),
			},
		},
		// 返却内容を記載するのを忘れない！！！！
		ReturnValues: aws.String("UPDATED_NEW"),
	}
	putItem, putErr := svc.UpdateItem(putParams)
	if putErr != nil {
		panic(fmt.Sprintf("error:%#v", putErr))
	}
	return putItem.Attributes["CurrentNumber"].N
}

func storeMetaData(i imageData) {
	// session
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	svc := dynamodb.New(sess)
	id := updateSequence(svc, os.Getenv("SEQUENCE_TABLE"))

	putParams := &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("IMAGE_META_DATA_TABLE")),
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				N: id,
			},
			"FileName": {
				S: aws.String(i.objectKey),
			},
		},
	}

	_, putErr := svc.PutItem(putParams)
	if putErr != nil {
		panic(putErr)
	}
}
