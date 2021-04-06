package cli

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/urfave/cli/v2"
	yaml "gopkg.in/yaml.v2"
)

type GruntConf struct {
	BucketName    string   `yaml:"BucketName"`
	StackName     string   `yaml:"StackName"`
	TemplateFiles []string `yaml:"TemplateFiles"`
}

func (conf *GruntConf) Parse(data []byte) error {
	return yaml.Unmarshal(data, conf)
}

//Parses the awsgrunt awsgrunt options file
func ParseAWSGruntOptions() (*GruntConf, error) {

	defaultConfigLocation, _ := filepath.Abs("./awsgrunt.yaml")

	data, err := ioutil.ReadFile(defaultConfigLocation)
	if err != nil {
		return nil, err
	}

	var config GruntConf
	err = config.Parse(data)

	if err != nil {
		return nil, err
	}
	return &config, nil
}

func TestConfigurationAction(c *cli.Context) error {
	gruntConf, err := ParseAWSGruntOptions()
	if err != nil {
		return err
	}

	fmt.Println("The following configuration has been loaded.")
	fmt.Printf("%+v\n", gruntConf)

	return nil
}
