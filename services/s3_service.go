package services

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

var AccessKeyID string
var SecretAccessKey string
var MyRegion string

func ConnectAWS() *session.Session {
	AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	MyRegion = os.Getenv("AWS_REGION")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(MyRegion),
		Credentials: credentials.NewStaticCredentials(
			AccessKeyID,
			SecretAccessKey,
			"",
		),
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return sess
}
