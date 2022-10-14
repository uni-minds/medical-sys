package tools

import (
	"gopkg.in/yaml.v3"
	"os"
)

func SaveYaml(file string, data interface{}) error {
	if fp, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600); err != nil {
		return err
	} else {
		defer fp.Close()
		return yaml.NewEncoder(fp).Encode(data)
	}
}

func LoadYaml(file string, data interface{}) error {
	if fp, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm); err != nil {
		return err
	} else {
		defer fp.Close()
		return yaml.NewDecoder(fp).Decode(data)
	}
}
