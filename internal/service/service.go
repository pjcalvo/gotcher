package service

import (
	"net/http"

	"github.com/pjcalvo/rigo/internal/config"
)

type ServiceType int

type InterceptService interface {
	HandleRequest(http.ResponseWriter, *http.Request) bool
	HandleResponse(*http.Response)
}

// Factory pattern used to build the right service based on the record flag
func NewInterceptService(c config.Config, record bool) InterceptService {
	if record {
		return newRecordService(c)
	}
	return newPatchService(c)
}
