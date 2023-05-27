package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ReadYML[T any](path string) (T, error) {
	var config T
	file, err := os.Open(path)
	fmt.Println()
	if err != nil {
		return config, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}
