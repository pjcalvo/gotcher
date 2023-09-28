package config

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
)

type Config struct {
	filePath string
	Values   *ConfigValues
}

func NewConfig(filePath string) (*Config, error) {
	config := &Config{filePath: filePath}

	// load config internally
	err := config.loadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	return config, nil
}

// LoadConfig reads the given file and returns a clean Config
func (c *Config) loadConfig() error {
	yamlFile, err := ioutil.ReadFile(c.filePath)
	if err != nil {
		return fmt.Errorf("error reading YAML file: %w", err)
	}

	//
	var config ConfigValues
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return fmt.Errorf("error unmarshaling YAML data: %w", err)
	}
	config.Intercept.Requests = cleanMatches(config.Intercept.Requests)
	config.Intercept.Responses = cleanMatches(config.Intercept.Responses)
	c.Values = &config

	return nil
}

func (c *Config) Watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("NewWatcher failed: ", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		defer close(done)
		log.Println("Watching for changes in", c.filePath)
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op.String() == "WRITE" {
					log.Println("Config file changed:", c.filePath)
					c.loadConfig()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(c.filePath)
	if err != nil {
		log.Fatal("Add failed:", err)
	}
	<-done
}
