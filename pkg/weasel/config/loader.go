package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/codetent/weasel/pkg/weasel/config/versions"
	"github.com/codetent/weasel/pkg/weasel/config/versions/v1alpha1"
	"gopkg.in/yaml.v3"
)

func Locate() (string, error) {
	next, err := os.Getwd()
	if err != nil {
		return "", err
	}

	current := ""
	for next != current {
		current = next
		configPath := filepath.Join(current, "weasel.yml")

		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}

		next = filepath.Dir(current)
	}

	return "", fmt.Errorf("configuration not found")
}

func Load(path string) (*v1alpha1.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	raw := &versions.Config{}
	yaml.Unmarshal(data, raw)

	if raw.Version != "v1alpha1" {
		return nil, fmt.Errorf("invalid version %s", raw.Version)
	}

	config := &v1alpha1.Config{}
	yaml.Unmarshal(data, config)

	return config, nil
}
