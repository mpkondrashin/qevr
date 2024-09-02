package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

/*
YAML configuration sample

version: 1
source: "c:\tmp\a.csv"
filter:
  type: nofilter
  networks:
    - 10.10.40.0/24
    - 172.15.0.0/16
output:
  target: SMS
  file:
    prefix: "dmz"
    folder: "C:\tmp"
  sms:
    address: "1.2.3.4"
    api_key: "abcd"
	ignore_tls_errors: true
*/

const configFileName = "config.yaml"

// Include only entities that mack the list below
// Include all entities that do no match the list below
//
//go:generate enum -type=FilterType -names=NoFilter,Include,Exclude

var FilterTypesLabels = []string{
	"No filtering",
	"Include only entities that mack the list below",       //"Leave entities with these IPs only"
	"Include all entities that do no match the list below", //"Remove entities for these IPs"
}

var MapFilterTypeLabelFromString = map[string]FilterType{
	FilterTypesLabels[FilterTypeNoFilter]: FilterTypeNoFilter,
	FilterTypesLabels[FilterTypeInclude]:  FilterTypeInclude,
	FilterTypesLabels[FilterTypeExclude]:  FilterTypeExclude,
}

//go:generate enum -type=Target -names=File,SMS

var TargetLabels = []string{
	"Save to file",
	"Upload to SMS server",
}

var MapTargetLabelFromString = map[string]Target{
	TargetLabels[TargetFile]: TargetFile,
	TargetLabels[TargetSMS]:  TargetSMS,
}

type (
	Filter struct {
		Type     FilterType `yaml:"type"`
		Networks []string   `yaml:"networks,omitempty"`
	}

	File struct {
		Prefix string `yaml:"prefix,omitempty"`
		Folder string `yaml:"folder,omitempty"`
	}

	SMS struct {
		Address         string        `yaml:"address,omitempty"`
		IgnoreTLSErrors bool          `yaml:"ignore_tls_errors"`
		APIKey          string        `yaml:"api_key,omitempty"`
		Timeout         time.Duration `yaml:"timeout"`
	}

	Output struct {
		Target Target `yaml:"target"`
		File   File   `yaml:"file,omitempty"`
		SMS    SMS    `yaml:"sms"`
	}

	Config struct {
		Version int    `yaml:"version"`
		Source  string `yaml:"source,omitempty"`
		Filter  Filter `yaml:"filter"`
		Output  Output `yaml:"output"`
	}
)

func ConfigFilePath() (string, error) {
	folder, err := ExecutableFolder()
	if err != nil {
		return "", err
	}
	return filepath.Join(folder, configFileName), nil
}

func DefaultConfig() Config {
	return Config{
		Version: 1,
		Filter:  Filter{Type: FilterTypeNoFilter},
		Output: Output{
			Target: TargetFile,
			SMS: SMS{
				IgnoreTLSErrors: false,
				Timeout:         5 * time.Second,
			},
		},
	}
}

func LoadConfig(configPath string) (*Config, error) {
	config := DefaultConfig()
	data, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &config, nil
		}
		return nil, fmt.Errorf("%s: failed to read config file: %w", configPath, err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to read config file: %w", configPath, err)
	}
	return &config, nil
}

func (c *Config) SaveConfig(folder string) error {
	configPath := filepath.Join(folder, configFileName)
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	err = os.WriteFile(configPath, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}
