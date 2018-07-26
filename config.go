package main

import (
	"errors"
	"io/ioutil"
	"regexp"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SlackToken string    `yaml:"slack_token"`
	Handlers   []Handler `yml:"handlers"`
}

func readConfig(path string) (*Config, error) {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := yaml.Unmarshal(body, config); err != nil {
		return nil, err
	}
	if config.SlackToken == "" {
		return nil, errors.New("Slack token is not provided")
	}

	for idx, h := range config.Handlers {
		re, err := regexp.Compile(h.Pattern)
		if err != nil {
			return nil, err
		}
		config.Handlers[idx].re = re
	}

	return config, nil
}

func (c *Config) FindHandler(input string) (*Handler, []string) {
	for _, h := range c.Handlers {
		match, args := h.Match(input)
		if match {
			return &h, args
		}
	}
	return nil, nil
}
