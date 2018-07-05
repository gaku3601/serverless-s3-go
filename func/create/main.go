package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/tidwall/gjson"
)

type imageData struct {
	original    string
	name        string
	objectKey   string
	file        io.Reader
	contentType string
}

func putData(file io.Reader, objectKey string, contentType string) {
	svc := s3.New(session.New())
	input := &s3.PutObjectInput{
		ACL:         aws.String("public-read"),
		Body:        aws.ReadSeekCloser(file),
		Bucket:      aws.String("gakustestbuckets2"),
		Key:         aws.String(objectKey),
		ContentType: aws.String(contentType),
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

func decodeData(inlineImageData string) (io.Reader, string) {
	ary := strings.Split(inlineImageData, ",")
	contentType := getContentType(ary[0])
	unbased, _ := base64.StdEncoding.DecodeString(ary[1])
	res := bytes.NewReader(unbased)
	return res, contentType
}

func getContentType(rawMeta string) string {
	r := strings.Replace(rawMeta, "data:", "", 1)
	r = strings.Replace(r, ";base64", "", 1)
	return r
}

func createObjectKey(name string) string {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	t := time.Now().In(jst)
	ft := t.Format("2006-01-02-15:04:05.000")
	return fmt.Sprintf("[%s]%s", ft, name)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	i := []imageData{}
	dataList := gjson.Get(request.Body, "dataList")
	dataList.ForEach(func(key, value gjson.Result) bool {
		file, contentType := decodeData(value.Get("data").String())
		d := imageData{
			original:    value.Get("data").String(),
			name:        value.Get("name").String(),
			objectKey:   createObjectKey(value.Get("name").String()),
			file:        file,
			contentType: contentType,
		}
		i = append(i, d)
		return true // keep iterating
	})
	for _, v := range i {
		putData(v.file, v.objectKey, v.contentType)
		storeMetaData(v)
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
