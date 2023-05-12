package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	TargetURL string  `yaml:"target_url"`
	Patches   []Patch `yaml:"patches"`
}

type Patch struct {
	Pattern  string `yaml:"pattern"`
	Status   int    `yaml:"status"`
	BodyFile string `yaml:"body_file"`
	Body     string `yaml:"body"`
}

// single config
var config Config

func main() {

	// read patches
	yamlFile, err := ioutil.ReadFile("gotcher.yaml")
	if err != nil {
		log.Fatalf("error reading YAML file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("error unmarshaling YAML data: %v", err)
	}
	config.Patches = cleanMatches(config.Patches)

	// Parse the target URL that we want to proxy to.
	targetURL, err := url.Parse(config.TargetURL)
	if err != nil {
		panic(err)
	}

	// Create a new reverse proxy that will forward requests to the target URL.
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Create an HTTP handler function that will serve as our proxy.
	proxyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request should be intercepted and handled separately.
		if ok, status, body := shouldPatch(r.RequestURI); ok {
			// Handle the intercepted request and return a custom response.
			handleInterceptedRequest(w, status, body)
			return
		}

		// Modify the request as needed before forwarding it to the target.
		r.Host = targetURL.Host
		r.URL.Scheme = targetURL.Scheme
		r.URL.Host = targetURL.Host

		// Forward the request to the target.
		proxy.ServeHTTP(w, r)
	})

	// Create a new HTTPS server with the TLS configuration and proxy handler.
	server := &http.Server{
		Addr: ":8443",
		// TLSConfig: tlsConfig,
		Handler: proxyHandler,
	}

	// Start the HTTPS server and listen on port 8443.
	log.Fatal(server.ListenAndServe())
}

// Handle an intercepted request and return a custom response.
func handleInterceptedRequest(w http.ResponseWriter, status int, body string) {
	// Add your custom response logic here.
	w.WriteHeader(status)
	w.Write([]byte(body))
}

func shouldPatch(uri string) (ok bool, status int, body string) {
	for _, patch := range config.Patches {
		matched, err := regexp.MatchString(patch.Pattern, uri)
		if err != nil {
			return
		}
		if matched {
			// patch if found
			if patch.Body != "" {
				body = patch.Body
			}
			// override the body with the content file
			if patch.BodyFile != "" {
				bytes, err := ioutil.ReadFile(patch.BodyFile)
				// when naked returns become nasty
				if err != nil {
					return
				}
				body = string(bytes)
			}
			if patch.Status != 0 {
				status = patch.Status
			}

			return true, status, body
		}
	}
	return
}

func cleanMatches(patches []Patch) []Patch {

	mPatches := make([]Patch, len(patches))
	for index, v := range patches {
		v.Pattern = strings.ReplaceAll(v.Pattern, "*", ".*")
		mPatches[index] = v
	}
	return mPatches
}
