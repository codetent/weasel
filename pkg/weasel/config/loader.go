package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/codetent/weasel/pkg/weasel/config/versions"
	"github.com/codetent/weasel/pkg/weasel/config/versions/v1alpha1"
	"gopkg.in/yaml.v3"
)

type ConfigFile struct {
	Path string
}

func LocateConfigFile() (*ConfigFile, error) {
	next, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	current := ""
	for next != current {
		current = next
		configPath := filepath.Join(current, "weasel.yml")

		if _, err := os.Stat(configPath); err == nil {
			config := &ConfigFile{
				Path: configPath,
			}
			return config, nil
		}

		next = filepath.Dir(current)
	}

	return nil, fmt.Errorf("configuration not found")
}

func (file *ConfigFile) Content() (*v1alpha1.Config, error) {
	data, err := os.ReadFile(file.Path)
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
