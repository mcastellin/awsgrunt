package cli

import (
	"awsgrunt/utils"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/urfave/cli/v2"
)

func ValidateTemplates(c *cli.Context) error {

	cfg, err := ParseAWSGruntOptions()
	if err != nil {
		return err
	}

	awscfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	cf := cloudformation.NewFromConfig(awscfg)

	// validating main and supporting templates
	allTemplates := append(cfg.TemplateFiles, cfg.StackTemplateFile)
	for _, path := range allTemplates {
		err = validateTemplateFile(cf, path)

		if err != nil {
			return err
		}
	}

	fmt.Println("All templates are valid!")

	return nil
}

func validateTemplateFile(client *cloudformation.Client, path string) error {
	templateBody, err := utils.ReadTemplateBodyFromFile(path)
	if err != nil {
		return err
	}

	output, err := client.ValidateTemplate(context.TODO(), &cloudformation.ValidateTemplateInput{
		TemplateBody: templateBody,
	})
	if err != nil {
		return err
	}

	if output == nil {
		return fmt.Errorf("Validation returned empty results. Template is invalid.:")
	}

	return nil
}
