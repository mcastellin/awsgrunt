package cli

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/urfave/cli/v2"
)

func DestroyStack(c *cli.Context) error {
	cfg, err := ParseAWSGruntOptions()
	if err != nil {
		return err
	}

	awscfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	cf := cloudformation.NewFromConfig(awscfg)

	// see if stack is already created
	stackInfo, err := findStack(cf, cfg.StackName, GetCreatedStackStatuses())

	if err != nil {
		return err
	} else if stackInfo != nil {
		// see if stack is in an updatable state
		updatableCheck, err := findStack(cf, cfg.StackName, GetUpdatableStackStatus())
		if err != nil {
			return err
		} else if updatableCheck == nil {
			return fmt.Errorf("The stack %s exists with status [%s] and cannot be deleted. Aborting", cfg.StackName, stackInfo.StackStatus)
		}
	} else {
		return fmt.Errorf("Stack with name %s does not exists. Aborting.", cfg.StackName)
	}

	fmt.Println("Deleting stack with name", cfg.StackName)
	_, err = cf.DeleteStack(context.TODO(), &cloudformation.DeleteStackInput{
		StackName: &cfg.StackName,
	})
	if err != nil {
		return err
	}
	fmt.Println("Done.")

	return nil
}
