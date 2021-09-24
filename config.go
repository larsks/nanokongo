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
		Controls []ControlSpec
	}

	ControlSpec struct {
		Control    uint8
		Type       string
		ScaleRange []int    `yaml:"scaleRange,omitempty"`
		OnRelease  []Action `yaml:"onRelease,omitempty"`
		OnPress    []Action `yaml:"onPress,omitempty"`
		OnChange   []Action `yaml:"onChange,omitempty"`
	}

	Action struct {
		Command   []string    `yaml:"command,omitempty"`
		SendKeys  []string    `yaml:"sendKeys,omitempty"`
		SendMouse []MouseSpec `yaml:"sendMouse,omitempty"`
	}

	MouseSpec struct {
		Press   string `yaml:"press,omitempty"`
		Release string `yaml:"release,omitempty"`
		Click   string `yaml:"click,omitempty"`
		X       int    `yaml:"x,omitempty"`
		Y       int    `yaml:"y,omitempty"`
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
