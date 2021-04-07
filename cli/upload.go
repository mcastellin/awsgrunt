package cli

import (
	grunts3 "awsgrunt/s3"
	"awsgrunt/utils"
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/urfave/cli/v2"
)

func UploadTemplatesToS3(c *cli.Context) error {

	cfg, err := ParseAWSGruntOptions()
	if err != nil {
		return err
	}

	awscfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	bucketName := grunts3.FindOrCreateBucket(awscfg, cfg.BucketName)

	client := s3.NewFromConfig(awscfg)

	fmt.Println("Uploading template files to s3...")
	grunts3.AddFilesToS3(client, bucketName, cfg.TemplateFiles)
	fmt.Println("Completed.")

	return uploadLambdaFunctions(cfg, client)
}

func uploadLambdaFunctions(cfg *GruntConf, s3client *s3.Client) error {

	filesToUpload := []string{}
	for _, lambda := range cfg.Lambdas {
		fmt.Println(fmt.Sprintf("generating archive for lambda function [%s]", lambda.Name))
		var source string
		if strings.HasSuffix(lambda.Path, "/") {
			source = lambda.Path
		} else {
			source = lambda.Path + "/"
		}

		outFile, err := utils.CreateZipFile(source, lambda.Name, "_lambda_releases")
		if err != nil {
			return err
		}
		filesToUpload = append(filesToUpload, *outFile)
	}

	fmt.Println("Uploading lambda functions zip files")
	grunts3.AddFilesToS3(s3client, cfg.BucketName, filesToUpload)

	fmt.Println("Done.")
	return nil
}
