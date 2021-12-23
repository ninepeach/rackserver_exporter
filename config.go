package main

import (
	"fmt"
	"io/ioutil"
	"sync"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Groups map[string]HostConfig `yaml:"groups"`
}

type SafeConfig struct {
	sync.RWMutex
	C *Config
}

type HostConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	BasicAuth string `yaml:"basicauth"`
}

func (sc *SafeConfig) ReloadConfig(configFile string) error {
	var c = &Config{}

	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(yamlFile, c); err != nil {
		return err
	}

	sc.Lock()
	sc.C = c
	sc.Unlock()

	return nil
}

// HostConfigForGroup checks the configuration for a matching group config and returns the configured HostConfig for
// that matched group.
func (sc *SafeConfig) HostConfigForGroup(group string) (*HostConfig, error) {
	sc.Lock()
	defer sc.Unlock()
	if hostConfig, ok := sc.C.Groups[group]; ok {
		return &hostConfig, nil
	}
	return &HostConfig{}, fmt.Errorf("no credentials found for group %s", group)
}