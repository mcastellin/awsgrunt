package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func askConfirmation(msg string) bool {
	fmt.Println(msg)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	return strings.HasPrefix(
		strings.ToLower(response), "y")
}

func askToCreateBucket(client *s3.Client, bucketName string) (string, error) {
	msg := fmt.Sprintf("Bucket with name %s does not exists. Would you like to create (Y/n)?", bucketName)
	if !askConfirmation(msg) {
		return "", errors.New("Cannot proceed without a bucket. Aborting.")
	}

	createBucketInput := &s3.CreateBucketInput{
		Bucket: &bucketName,
	}
	_, err := client.CreateBucket(context.TODO(), createBucketInput)
	if err != nil {
		log.Fatal(err)
	}

	return bucketName, nil
}

func FindOrCreateBucket(cfg aws.Config, bucketName string) string {
	client := s3.NewFromConfig(cfg)

	input := &s3.ListBucketsInput{}
	buckets, err := client.ListBuckets(context.TODO(), input)
	if err != nil {
		log.Fatal(err)
	}

	for _, bucket := range buckets.Buckets {
		if *bucket.Name == bucketName {
			// found valid bucket. returning
			return *bucket.Name
		}
	}

	_, err = askToCreateBucket(client, bucketName)
	if err != nil {
		log.Fatal(err)
	}
	return bucketName
}

func AddFilesToS3(client *s3.Client, bucketName string, filePaths []string) (*s3.PutObjectOutput, error) {
	for _, path := range filePaths {
		fmt.Println(path)
		_, err := AddFileToS3(client, bucketName, path)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func AddFileToS3(client *s3.Client, bucketName string, path string) (*s3.PutObjectOutput, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	return client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:               aws.String(bucketName),
		Key:                  aws.String(path),
		Body:                 bytes.NewReader(buffer),
		ACL:                  types.ObjectCannedACLPrivate,
		ContentLength:        size,
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: types.ServerSideEncryptionAes256,
	})
}
