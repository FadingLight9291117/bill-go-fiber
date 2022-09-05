package lib

import (
	"bill/model"
	"gopkg.in/yaml.v2"
	"os"
)

func GetClsAndLabel(path string) (*model.Class, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	class := new(model.Class)
	err = yaml.Unmarshal(yamlFile, class)
	if err != nil {
		return nil, err
	}
	return class, nil
}
