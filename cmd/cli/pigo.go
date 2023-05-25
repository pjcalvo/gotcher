package cli

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/fatih/color"
	"github.com/pjcalvo/rigo/internal/config"
	"github.com/pjcalvo/rigo/internal/service"
	"github.com/rs/cors"
)

// single config with default values
var (
	port       int    = 8443
	configPath string = "rigo.yaml"
	verbose    bool   = false
	record     bool   = false
)

func Run() error {
	// flags
	flag.IntVar(&port, "p", port, "port number to run the proxy")
	flag.StringVar(&configPath, "f", configPath, "file path to be used as the config")
	flag.BoolVar(&verbose, "v", verbose, "file path to be used as the config")
	flag.BoolVar(&record, "record", record, "wheter or not to record instead of intercept")
	flag.Parse()

	interceptConfig, err := config.LoadConfig(configPath)
	if err != nil {
		return err
	}

	interceptService := service.NewInterceptService(*interceptConfig, record)

	// Parse the target URL that we want to proxy to.
	targetURL, err := url.Parse(interceptConfig.TargetURL)
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
				req.Header["Authorization"] = []string{fmt.Sprintf("%s %s", interceptConfig.Authentication.Bearer.Type, interceptConfig.Authentication.Bearer.Token)}
			}
		},
		ModifyResponse: func(resp *http.Response) error {
			// modify the request otherwise return it as it is
			interceptService.HandleResponse(resp)
			return nil
		},
	}

	// Create an HTTP handler function that will serve as our proxy.
	proxyHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Handle preflight requests always OK
		if req.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if !interceptService.HandleRequest(w, req) {
			// Forward the request to the target.
			proxy.ServeHTTP(w, req)
		}
	})

	// handle cors issues
	c := cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowCredentials: true,
		AllowedHeaders:   []string{"x-internal-session-id", "content-type"},
		AllowedMethods:   []string{"GET", "PUT", "OPTIONS", "POST"},
		// Enable Debugging for testing, consider disabling in production
		Debug: verbose,
	})

	// create a new HTTPS server with the TLS configuration and proxy handler.
	server := &http.Server{
		Addr: fmt.Sprintf(":%v", port),
		// TLSConfig: tlsConfig,
		Handler: c.Handler(proxyHandler),
	}

	// start the HTTPS server and listen on port default to 8443
	// but overriden by the port flag.
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("[ %s ] on http://localhost:%v\n", green("ready"), port)
	log.Fatal(server.ListenAndServe())
	return nil
}
