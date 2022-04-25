package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	BotServer struct {
		Schema string `yaml:"schema"`
		Host   string `yaml:"host"`
		Listen string `yaml:"listen"`
		ApiKey string `yaml:"apiKey"`
	} `yaml:"bot_server"`
	CheveretoServer struct {
		Schema string `yaml:"schema"`
		Host   string `yaml:"host"`
		ApiKey string `yaml:"apiKey"`
	} `yaml:"chevereto_server"`
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

func ParseFlags() (string, error) {
	var configPath string

	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")
	flag.Parse()

	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	return configPath, nil
}
