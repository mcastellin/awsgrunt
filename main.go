package main

import (
	"fmt"
	"os"

	cli2 "awsgrunt/cli"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "conf",
				Aliases: []string{"c"},
				Usage:   "Tests the configuration",
				Action:  cli2.TestConfigurationAction,
			},
			{
				Name:    "validate",
				Aliases: []string{"v"},
				Usage:   "Validates all configured cloudformation templates",
				Action:  cli2.ValidateTemplates,
			},
			{
				Name:    "upload",
				Aliases: []string{"u"},
				Usage:   "Uploads the configured template files to the S3 bucket",
				Action:  cli2.UploadTemplatesToS3,
			},
			{
				Name:    "apply",
				Aliases: []string{"a"},
				Usage:   "Creates or updates the cloudformation stack",
				Action:  cli2.ApplyStack,
			},
			{
				Name:   "destroy",
				Usage:  "Deletes the configured stack",
				Action: cli2.DestroyStack,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
