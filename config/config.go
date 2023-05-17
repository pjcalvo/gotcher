package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type BodyType string

const (
	BodyTypeFile   BodyType = "file"
	BodyTypeString BodyType = "string"
	BodyTypeJson   BodyType = "json"
)

type Config struct {
	TargetURL string      `yaml:"target_url"`
	Patches   PatchGroups `yaml:"patches"`
	Token     string      `yaml:"token"`
}

type PatchGroups struct {
	Responses []Patch `yaml:"responses"`
	Requests  []Patch `yaml:"requests"`
}

type Patch struct {
	Pattern string   `yaml:"pattern"`
	Status  int      `yaml:"status"`
	Body    string   `yaml:"body"`
	Type    BodyType `yaml:"type"`
}

func LoadConfig() (*Config, error) {
	// read patches
	yamlFile, err := ioutil.ReadFile("pigo.yaml")
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML data: %w", err)
	}
	config.Patches.Requests = cleanMatches(config.Patches.Requests)
	config.Patches.Responses = cleanMatches(config.Patches.Responses)

	return &config, nil
}

func cleanMatches(patches []Patch) []Patch {
	mPatches := make([]Patch, len(patches))
	for index, v := range patches {
		v.Pattern = strings.ReplaceAll(v.Pattern, "*", ".*")
		mPatches[index] = v
	}
	return mPatches
}
