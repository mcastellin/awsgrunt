package cli

import (
	"awsgrunt/utils"
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/urfave/cli/v2"
)

func GetUpdatableStackStatus() []types.StackStatus {
	return []types.StackStatus{
		types.StackStatusCreateFailed,
		types.StackStatusCreateComplete,
		types.StackStatusRollbackFailed,
		types.StackStatusRollbackComplete,
		types.StackStatusDeleteFailed,
		types.StackStatusUpdateComplete,
		types.StackStatusUpdateRollbackFailed,
		types.StackStatusUpdateRollbackComplete,
		types.StackStatusImportComplete,
		types.StackStatusImportRollbackFailed,
		types.StackStatusImportRollbackComplete,
	}
}

func GetCreatedStackStatuses() []types.StackStatus {
	return []types.StackStatus{
		types.StackStatusCreateInProgress,
		types.StackStatusCreateFailed,
		types.StackStatusCreateComplete,
		types.StackStatusRollbackInProgress,
		types.StackStatusRollbackFailed,
		types.StackStatusRollbackComplete,
		types.StackStatusDeleteInProgress,
		types.StackStatusDeleteFailed,
		types.StackStatusUpdateInProgress,
		types.StackStatusUpdateCompleteCleanupInProgress,
		types.StackStatusUpdateComplete,
		types.StackStatusUpdateRollbackInProgress,
		types.StackStatusUpdateRollbackFailed,
		types.StackStatusUpdateRollbackCompleteCleanupInProgress,
		types.StackStatusUpdateRollbackComplete,
		types.StackStatusReviewInProgress,
		types.StackStatusImportInProgress,
		types.StackStatusImportComplete,
		types.StackStatusImportRollbackInProgress,
		types.StackStatusImportRollbackFailed,
		types.StackStatusImportRollbackComplete,
	}
}

func ApplyStack(c *cli.Context) error {

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

	isUpdate := false
	if err != nil {
		return err
	} else if stackInfo != nil {
		// see if stack is in an updatable state
		updatableCheck, err := findStack(cf, cfg.StackName, GetUpdatableStackStatus())
		if err != nil {
			return err
		} else if updatableCheck != nil {
			isUpdate = true
		} else {
			return errors.New(fmt.Sprintf("The stack %s exists with status [%s] and cannot be updated. Aborting", cfg.StackName, stackInfo.StackStatus))
		}
	}

	templateBody, err := utils.ReadTemplateBodyFromFile(cfg.StackTemplateFile)
	if err != nil {
		return err
	}

	capabilities, err := resolveCapabilities(cfg.Capabilities)
	if err != nil {
		return err
	}

	var stackId string
	if isUpdate {
		fmt.Println("Updating existing stack with name", cfg.StackName)
		output, err := cf.UpdateStack(context.TODO(), &cloudformation.UpdateStackInput{
			StackName:    &cfg.StackName,
			Capabilities: *capabilities,
			TemplateBody: templateBody,
			Parameters:   prepareCfParameters(cfg.Parameters),
		})
		if err != nil {
			return err
		} else {
			stackId = *output.StackId
		}
	} else {
		fmt.Println("Creating the new stack with name", cfg.StackName)
		output, err := cf.CreateStack(context.TODO(), &cloudformation.CreateStackInput{
			StackName:    &cfg.StackName,
			Capabilities: *capabilities,
			TemplateBody: templateBody,
			Parameters:   prepareCfParameters(cfg.Parameters),
		})
		if err != nil {
			return err
		} else {
			stackId = *output.StackId
		}
	}

	fmt.Println(fmt.Sprintf("StackId: %s", stackId))

	return nil
}

func findStack(client *cloudformation.Client, stackName string, statusFilter []types.StackStatus) (*types.StackSummary, error) {
	listStacksOut, err := client.ListStacks(context.TODO(), &cloudformation.ListStacksInput{
		StackStatusFilter: statusFilter,
	})
	if err != nil {
		return nil, err
	}
	for _, stack := range listStacksOut.StackSummaries {
		if *stack.StackName == stackName {
			return &stack, nil
		}
	}
	return nil, nil
}

func prepareCfParameters(params map[string]string) []types.Parameter {
	cfParams := []types.Parameter{}

	for k, v := range params {
		key, value := k, v // force creation of a new address space
		cfParams = append(cfParams, types.Parameter{
			ParameterKey:   &key,
			ParameterValue: &value,
		})
	}

	return cfParams
}

func resolveCapabilities(in []string) (*[]types.Capability, error) {
	resolvedValues := []types.Capability{}
	for _, capability := range in {
		var resolved types.Capability
		switch capability {
		case "CAPABILITY_IAM":
			resolved = types.CapabilityCapabilityIam
		case "CAPABILITY_NAMED_IAM":
			resolved = types.CapabilityCapabilityNamedIam
		case "CAPABILITY_AUTO_EXPAND":
			resolved = types.CapabilityCapabilityAutoExpand
		default:
			return nil, fmt.Errorf("Invalid configuration attribute. Capability with name [%s] cannot be resolved.", capability)
		}

		resolvedValues = append(resolvedValues, resolved)
	}

	return &resolvedValues, nil
}
