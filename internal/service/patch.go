package service

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pjcalvo/rigo/internal/config"
)

type PatchService struct {
	interceptConfig config.Config
}

func newPatchService(c config.Config) PatchService {
	return PatchService{
		interceptConfig: c,
	}
}

// getPatchDetails checks if a request should be patched and return the proper details
func getPatchDetails(request *http.Request, intercepts []config.Intercept) (ok bool, status int, body []byte) {
	for _, intercept := range intercepts {
		matched, err := isPatchable(verifyInterceptParams{
			uri:    request.URL.String(),
			method: request.Method,
		}, intercept)
		if err != nil {
			return
		}
		if matched {
			// define body to return
			switch intercept.Patch.Type {
			case config.BodyTypeFile:
				body, err = ioutil.ReadFile(intercept.Patch.Body)
				if err != nil {
					return
				}
			case config.BodyTypeString, config.BodyTypeJson:
				body = []byte(intercept.Patch.Body)
			}
			// define status to return (default to 200)
			status = 200
			if intercept.Patch.Status != 0 {
				status = intercept.Patch.Status
			}

			return true, status, body
		}
	}
	return
}

func (i PatchService) HandleRequest(w http.ResponseWriter, r *http.Request) bool {
	if ok, status, body := getPatchDetails(r, i.interceptConfig.Intercept.Requests); ok {
		// Handle the intercepted request and return a custom response.
		fmt.Printf("Patching REQUEST for: %s\n	status: %v\n", r.RequestURI, status)

		w.WriteHeader(status)
		w.Write(body)
		return true
	}
	return false
}

func (i PatchService) HandleResponse(r *http.Response) {
	if ok, status, body := getPatchDetails(r.Request, i.interceptConfig.Intercept.Responses); ok {
		// Handle the intercepted request and return a custom response.
		// Tood: implement logging
		fmt.Printf("Patching RESPONSE for: %s\n	status: %v\n", r.Request.URL.String(), status)
		r.Body = io.NopCloser(bytes.NewReader(body))
		r.ContentLength = int64(len(body))
		r.StatusCode = status
	}
}
