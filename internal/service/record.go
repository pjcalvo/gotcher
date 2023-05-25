package service

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pjcalvo/rigo/internal/config"
)

const (
	maxCharsFileName     = 200
	defaultFileExtension = "json"
)

type RecordService struct {
	interceptConfig config.Config
}

func newRecordService(c config.Config) RecordService {
	// record invert the requests
	c.Intercept.Requests, c.Intercept.Responses = c.Intercept.Responses, c.Intercept.Requests
	return RecordService{
		interceptConfig: c,
	}
}

// shouldRecord records the response to a file
func shouldRecord(response *http.Response, intercept config.Intercept) (ok bool) {
	ok, err := isPatchable(verifyInterceptParams{
		uri:    response.Request.URL.String(),
		method: response.Request.Method,
	}, intercept)
	if err != nil {
		return
	}

	if ok {
		// write response to file
		// if the filename is specified then use it as the name
		// otherwise build one based on
		var filename string
		if intercept.Patch.Type == config.BodyTypeFile {
			filename = intercept.Patch.Body
		} else {
			filename = fmt.Sprintf("%s_%s.%s", response.Request.Method, response.Request.URL, defaultFileExtension)
			filename = strings.ReplaceAll(filename, "/", "_")
			if len(filename) > maxCharsFileName {
				filename = filename[:maxCharsFileName]
			}
		}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Println("Error reading the body:", err)
			return
		}

		err = os.WriteFile(filename, body, 0644)
		if err != nil {
			log.Println("Error writing to the file:", err)
			return
		}

		// re-append the body to the response
		bodyReader := io.NopCloser(bytes.NewBuffer(body))
		response.Body = bodyReader
		return true
	}
	return
}

func (i RecordService) HandleRequest(w http.ResponseWriter, r *http.Request) bool {
	// not implemented we might want to limit this functionality
	for range i.interceptConfig.Intercept.Requests {
		return false
	}
	return false
}

func (i RecordService) HandleResponse(r *http.Response) {
	for _, intercept := range i.interceptConfig.Intercept.Responses {
		if ok := shouldRecord(r, intercept); ok {
			fmt.Printf("Recording RESPONSE for: %s\n	status: %v\n", r.Request.URL.String(), r.Status)
		}
	}
}
