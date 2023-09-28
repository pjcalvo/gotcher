package config

import (
	"strings"
)

type BodyType string

const (
	BodyTypeFile   BodyType = "file"
	BodyTypeString BodyType = "string"
	BodyTypeJson   BodyType = "json"
)

// ConfigValues is the top level configuration
type ConfigValues struct {
	TargetURL      string         `yaml:"target_url"`
	Authentication Authentication `yaml:"authentication"`
	Intercept      InterceptGroup `yaml:"intercept"`
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
