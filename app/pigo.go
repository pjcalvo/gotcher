package app

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"

	"github.com/fatih/color"
	"github.com/pjcalvo/pigo/config"
	"github.com/rs/cors"
)

// single config
var patchConfig *config.Config
var port int

func Run() error {
	// flags
	flag.IntVar(&port, "p", 8443, "port number to run the proxy")
	flag.Parse()

	// config
	cons, err := config.LoadConfig()
	if err != nil {
		return err
	}
	patchConfig = cons

	// Parse the target URL that we want to proxy to.
	targetURL, err := url.Parse(patchConfig.TargetURL)
	if err != nil {
		return err
	}

	// Create a new reverse proxy that will forward requests to the target URL.
	// Create a new reverse proxy with a custom Director function.
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			// Modify the request as needed before forwarding it to the target.
			req.Host = targetURL.Host
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host

			// Todo: parameterize
			if _, ok := req.Header["Authorization"]; !ok {
				req.Header["Authorization"] = []string{fmt.Sprintf("Bearer %s", patchConfig.Token)}
			}
		},
		ModifyResponse: func(resp *http.Response) error {
			// modify the request otherwise return it as it is
			if ok, status, body := shouldPatchResponse(resp); ok {
				// Handle the intercepted request and return a custom response.
				fmt.Printf("Patching RESPONSE for: %s\n	status: %v\n", resp.Request.URL.String(), status)
				resp.Body = ioutil.NopCloser(bytes.NewReader(body))
				resp.ContentLength = int64(len(body))
				resp.StatusCode = status
				return nil
			}

			return nil
		},
	}

	// Create an HTTP handler function that will serve as our proxy.
	proxyHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Handle preflight requests
		if req.Method == http.MethodOptions {
			// Set the necessary CORS headers
			w.WriteHeader(http.StatusOK)
			return
		}

		if ok, status, body := shouldPatchRequest(req.RequestURI); ok {
			// Handle the intercepted request and return a custom response.
			fmt.Printf("Patching REQUEST for: %s\n	status: %v\n", req.RequestURI, status)
			handleInterceptedRequest(w, status, body)
			return
		}
		// Forward the request to the target.
		proxy.ServeHTTP(w, req)
	})

	c := cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowCredentials: true,
		AllowedHeaders:   []string{"x-internal-session-id", "content-type"},
		AllowedMethods:   []string{"GET", "PUT", "OPTIONS", "POST"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})
	// Create a new HTTPS server with the TLS configuration and proxy handler.
	server := &http.Server{
		Addr: fmt.Sprintf(":%v", port),
		// TLSConfig: tlsConfig,
		Handler: c.Handler(proxyHandler),
	}

	// Start the HTTPS server and listen on port 8443.
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("[ %s ] on http://localhost:%v\n", green("ready"), port)
	log.Fatal(server.ListenAndServe())
	return nil
}

// Handle an intercepted request and return a custom response.
func handleInterceptedRequest(w http.ResponseWriter, status int, body []byte) {
	w.WriteHeader(status)
	w.Write(body)
}

func shouldPatchRequest(uri string) (ok bool, status int, body []byte) {
	for _, patch := range patchConfig.Patches.Requests {
		if patch.Pattern == "" {
			return
		}

		matched, err := regexp.MatchString(patch.Pattern, uri)
		if err != nil {
			return
		}
		if matched {
			switch patch.Type {
			case config.BodyTypeFile:
				body, err = ioutil.ReadFile(patch.Body)
				if err != nil {
					return
				}
			case config.BodyTypeString, config.BodyTypeJson:
				body = []byte(patch.Body)
				// override the body with the content file
			}

			if patch.Status != 0 {
				status = patch.Status
			}

			return true, status, body
		}
	}
	return
}

func shouldPatchResponse(resp *http.Response) (ok bool, status int, body []byte) {
	// naked returns for the win
	for _, patch := range patchConfig.Patches.Responses {
		if patch.Pattern == "" {
			return
		}
		matched, err := regexp.MatchString(patch.Pattern, resp.Request.URL.String())
		if err != nil {
			return
		}
		if matched {
			switch patch.Type {
			case config.BodyTypeFile:
				body, err = ioutil.ReadFile(patch.Body)
				if err != nil {
					return
				}
			case config.BodyTypeString, config.BodyTypeJson:
				body = []byte(patch.Body)
				// override the body with the content file
			}

			if patch.Status != 0 {
				status = patch.Status
			}

			return true, status, body
		}
	}
	return
}
