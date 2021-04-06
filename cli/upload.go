package cli

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/urfave/cli/v2"

	grunts3 "awsgrunt/s3"
)

func UploadTemplatesToS3(c *cli.Context) error {

	cfg, err := ParseAWSGruntOptions()

	awscfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	bucketName := grunts3.FindOrCreateBucket(awscfg, cfg.BucketName)

	client := s3.NewFromConfig(awscfg)

	fmt.Println("Uploading files to s3...")
	grunts3.AddFilesToS3(client, bucketName, cfg.TemplateFiles)
	fmt.Println("Completed.")

	return nil
}
