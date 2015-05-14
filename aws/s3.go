package aws

import (
	"github.com/AdRoll/goamz/aws"
	"github.com/AdRoll/goamz/s3"
	"log"
	"os"
	"time"
)

var bucket *s3.Bucket

func Init() error {
	log.Println("Getting auth")
	auth, err := aws.GetAuth(os.Getenv("AWS_ACCESS"), os.Getenv("AWS_SECRET"), "", time.Now().Add(time.Hour))
	if err != nil {
		return err
	}
	client := s3.New(auth, aws.APSoutheast)
	switch os.Getenv("ENV") {
	case "dev":
		bucket = client.Bucket("lemonades-staging")
	case "staging":
		bucket = client.Bucket("lemonades-staging")
	case "prod":
		bucket = client.Bucket("lemonades-prod")
	}
	return nil
}

func Bucket() *s3.Bucket {
	return bucket
}
