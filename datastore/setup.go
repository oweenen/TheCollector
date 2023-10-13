package datastore

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var client *s3.S3

func SetupConnection() {
	key := os.Getenv("SPACES_KEY")
	secret := os.Getenv("SPACES_SECRET_KEY")

	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:         aws.String("https://sfo2.digitaloceanspaces.com"),
		S3ForcePathStyle: aws.Bool(false),
		Region:           aws.String("sfo2")},
	)
	if err != nil {
		log.Fatalf("Failed to connect to spaces: %v\n", err)
	}

	client = s3.New(sess)

	fmt.Println("Successfully connected to spaces!")
}
