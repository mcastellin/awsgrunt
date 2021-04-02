package main

import (
	"fmt"
	"log"
	"os"

	cli2 "awsgrunt/cli"

	"github.com/urfave/cli"
)

func test(c *cli.Context) error {
	fmt.Println("Booooom!!")
	cli2.ParseAWSGruntOptions()
	return nil
}

func main() {

	app := &cli.App{
		Name:   "awsgrunt",
		Usage:  "work with cloudformation nested stacks",
		Action: test,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
