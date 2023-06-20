package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
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
	FilePath	   string  
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

func NewConfig(filePath string) *Config {
	return &Config{FilePath: filePath}
  }

// LoadConfig reads the given file and returns a clean Config
func (config *Config) LoadConfig() error {
	log.Println("Loading testing configuration")
	yamlFile, err := ioutil.ReadFile(config.FilePath)
	if err != nil {
		return fmt.Errorf("error reading YAML file: %w", err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return fmt.Errorf("error unmarshaling YAML data: %w", err)
	}
	config.Intercept.Requests = cleanMatches(config.Intercept.Requests)
	config.Intercept.Responses = cleanMatches(config.Intercept.Responses)

	return nil
}

func (config *Config) Watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("NewWatcher failed: ", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		defer close(done)
		log.Println("Watching for changes in", config.FilePath)
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op.String() == "WRITE" {
					log.Println("Config file changed:", config.FilePath)
					config.LoadConfig()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(config.FilePath)
	if err != nil {
		log.Fatal("Add failed:", err)
	}
	<-done
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