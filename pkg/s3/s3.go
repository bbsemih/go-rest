package s3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/bbsemih/gobank/pkg/util"
	"github.com/rs/zerolog/log"
	"os"
)

var (
	region string
	sess   *session.Session
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func init() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load config")
	}
	region = config.AWSRegion
	bucketName := config.AWSBucketName

	sess, err = session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		exitErrorf("Unable to create session, %v", err)
	}

	svc := s3.New(sess)

	ListBuckets(svc)
	CreateBucket(svc, bucketName)
	UploadFile(svc, bucketName, "test.txt")
}

func ListBuckets(svc *s3.S3) {
	result, err := svc.ListBuckets(nil)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	log.Info().Msg("Buckets:")

	for _, b := range result.Buckets {
		log.Info().Msg(*b.Name + "\t" + b.CreationDate.String())
	}

	fmt.Printf("\n")
}

func CreateBucket(svc *s3.S3, bucketName string) {
	if bucketName == "" {
		exitErrorf("Bucket name cannot be empty")
	}
	log.Info().Msg("Creating bucket " + bucketName + "...")

	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(region),
		},
	})
	if err != nil {
		exitErrorf("Unable to create bucket, %v", err)
	}

	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		exitErrorf("Error occurred while waiting for bucket to be created, %v", bucketName)
	}
	log.Info().Msg("Bucket " + bucketName + " created")
}

func UploadFile(svc *s3.S3, bucketName string, fileName string) {
	if bucketName == "" || fileName == "" {
		exitErrorf("Bucket name and file name cannot be empty")
	}

	file, err := os.Open(fileName)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)
	}

	defer file.Close()

	log.Info().Msg("Uploading file " + fileName + " to bucket " + bucketName + "...")

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		exitErrorf("Unable to upload file, %v", err)
	}

	log.Info().Msg("File " + fileName + " uploaded to bucket " + bucketName)
}
