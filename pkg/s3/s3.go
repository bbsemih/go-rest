package s3

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rs/zerolog/log"
)

var (
	region string
	sess   *session.Session
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func InitS3(bucketName string, region string) {
	sess, err := session.NewSession(&aws.Config{
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

func DeleteItem(svc *s3.S3, bucketName string, fileName string) {
	if bucketName == "" || fileName == "" {
		exitErrorf("Bucket name and file name cannot be empty")
	}

	log.Info().Msg("Deleting file " + fileName + " from bucket " + bucketName + "...")

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucketName), Key: aws.String(fileName)})
	if err != nil {
		exitErrorf("Unable to delete object %q from bucket %q, %v", fileName, bucketName, err)
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		exitErrorf("Error occurred while waiting for object %q to be deleted, %v", fileName, err)
	}

	log.Info().Msg("File " + fileName + " deleted from bucket " + bucketName)
}

func DeleteAllItems(svc *s3.S3, bucketName string) {
	if bucketName == "" {
		exitErrorf("Bucket name cannot be empty")
	}

	log.Info().Msg("Deleting all files from bucket " + bucketName + "...")

	iter := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	})

	if err := s3manager.NewBatchDeleteWithClient(svc).Delete(aws.BackgroundContext(), iter); err != nil {
		exitErrorf("Unable to delete objects from bucket %q, %v", bucketName, err)
	}

	log.Info().Msg("Deleted object(s) from bucket: %s" + bucketName)
}

func DeleteBucket(svc *s3.S3, bucketName string) {
	if bucketName == "" {
		exitErrorf("Bucket name cannot be empty")
	}

	_, err := svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		exitErrorf("Unable to delete bucket %q, %v", bucketName, err)
	}

	fmt.Printf("Waiting for bucket %q to be deleted...\n", bucketName)

	err = svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		exitErrorf("Error occurred while waiting for bucket to be deleted, %v", bucketName)
	}

	log.Info().Msg("Bucket " + bucketName + " deleted successfully")
}
