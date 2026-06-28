package conf

import (
	"os"

	"gopkg.in/yaml.v3"
)

// LoadYAML loads a YAML configuration file into a struct of type T.
func LoadYAML[T any](path string) (*T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := new(T)
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
