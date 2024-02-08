package env

import (
	"os"

	"gopkg.in/yaml.v3"
)

// LoadYaml
func LoadYaml(file_path string, obj interface{}) error {
	data, err := os.ReadFile(file_path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, obj)
	if err != nil {
		return err
	}
	return nil
}
