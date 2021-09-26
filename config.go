package main

import (
	"io"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type (
	Config struct {
		Device   string
		Channel  uint8
		Controls map[uint8]ControlSpec
	}

	ControlSpec struct {
		Type       string
		ScaleRange []int                  `yaml:"scaleRange,omitempty"`
		OnRelease  []map[string]yaml.Node `yaml:"onRelease,omitempty"`
		OnPress    []map[string]yaml.Node `yaml:"onPress,omitempty"`
		OnChange   []map[string]yaml.Node `yaml:"onChange,omitempty"`
	}
)

func ReadConfig(r io.Reader) (*Config, error) {
	var config Config

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func ReadConfigFromFile(path string) (*Config, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return ReadConfig(r)
}
