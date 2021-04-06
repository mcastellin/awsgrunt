package cli

import "io/ioutil"

func readTemplateBodyFromFile(path string) (*string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	template := string(data)
	return &template, nil
}