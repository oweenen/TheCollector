package datastore

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var client *s3.S3
var bucket = aws.String("tft-stats-match-data")

func SetupConnection() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION"))},
	)
	if err != nil {
		log.Fatalf("Failed to connect to AWS: %v\n", err)
	}

	client = s3.New(sess)

	fmt.Println("Successfully connected to Amazon S3!")
}
