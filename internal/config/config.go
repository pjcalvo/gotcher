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

// Config is the top level configuration
type Config struct {
	TargetURL      string         `yaml:"target_url"`
	Authentication Authentication `yaml:"authentication"`
	Intercept      InterceptGroup `yaml:"intercept"`
	Token          string         `yaml:"token"`
}

// Authentication is provided in cases where the authentication
// has to be fetched by the proxy or hardcoded. The main use case
// for this is when testing Oauth and the authentication is sent by
// cookies so the proxy won't receive them
type Authentication struct {
	Basic struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"basic"`
	Bearer struct {
		Type  string `yaml:"type"`
		Token string `yaml:"token"`
	} `yaml:"bearer"`
}

type InterceptGroup struct {
	Responses []Intercept `yaml:"responses"`
	Requests  []Intercept `yaml:"requests"`
}

type Intercept struct {
	Match Match `yaml:"match"`
	Patch Patch `yaml:"patch"`
}

type Patch struct {
	Status int      `yaml:"status"`
	Body   string   `yaml:"body"`
	Type   BodyType `yaml:"type"`
}

type Match struct {
	Uri    string `yaml:"uri"`
	Params []struct {
		Name  string `yaml:"name"`
		Value string `yaml:"value"`
	} `yaml:"params"`
	Methods []string `yaml:"methods"`
}

// LoadConfig reads the given file and returns a clean Config
func LoadConfig(filePath string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML data: %w", err)
	}
	config.Intercept.Requests = cleanMatches(config.Intercept.Requests)
	config.Intercept.Responses = cleanMatches(config.Intercept.Responses)

	return &config, nil
}

// cleanMatches is a function that converts * to
// acceptable .* patterns
func cleanMatches(patches []Intercept) []Intercept {
	mPatches := make([]Intercept, len(patches))
	for index, v := range patches {
		v.Match.Uri = strings.ReplaceAll(v.Match.Uri, "*", ".*")
		mPatches[index] = v
	}
	return mPatches
}
