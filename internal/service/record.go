package service

import (
	"net/http"

	"github.com/pjcalvo/rigo/internal/config"
)

type RecordService struct {
	interceptConfig config.Config
}

func newRecordService(c config.Config) RecordService {
	return RecordService{
		interceptConfig: c,
	}
}

// intercept handles the logic to match and return the proper response
// should split more
func record(intercepts []config.Intercept, uri string) (ok bool, status int, body []byte) {
	panic("not implemented")
}

func (i RecordService) HandleRequest(w http.ResponseWriter, r *http.Request) bool {
	// return intercept(i.interceptConfig.Intercept.Requests, uri)
	return false
}

func (i RecordService) HandleResponse(r *http.Response) {
	// return intercept(i.interceptConfig.Intercept.Responses, uri)
}
