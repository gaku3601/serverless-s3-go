package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/tidwall/gjson"
)

type imageData struct {
	original string
	name     string
	binary   string
}

func putData(file io.Reader) {
	svc := s3.New(session.New())
	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(file),
		Bucket: aws.String("gakustestbuckets2"),
		Key:    aws.String("test.jpg"),
	}
	_, err := svc.PutObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println("error:" + aerr.Error() + "||" + err.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
}

func decodeData(inlineImageData string) io.Reader {
	ary := strings.Split(inlineImageData, ",")
	unbased, _ := base64.StdEncoding.DecodeString(ary[1])
	res := bytes.NewReader(unbased)
	return res
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//	data, _ := base64.StdEncoding.DecodeString(request.Body)
	i := []imageData{}
	dataList := gjson.Get(request.Body, "dataList")
	dataList.ForEach(func(key, value gjson.Result) bool {
		d := imageData{
			original: value.Get("data").String(),
			name:     value.Get("name").String(),
		}
		i = append(i, d)
		return true // keep iterating
	})
	for _, v := range i {
		putData(decodeData(v.original))
	}
	return events.APIGatewayProxyResponse{
		Body: fmt.Sprintf("%#v", i),
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
		},
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}