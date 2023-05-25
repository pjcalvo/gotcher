package service

import (
	"regexp"

	"github.com/pjcalvo/rigo/internal/config"
	"github.com/pjcalvo/rigo/internal/stuff"
)

type verifyInterceptParams struct {
	method string
	uri    string
	// Todo: complete this verification as well
	params []struct {
		name  string
		value string
	}
}

func isPatchable(v verifyInterceptParams, intercept config.Intercept) (bool, error) {
	// conditions to break the matching process
	if len(intercept.Match.Methods) > 0 && !stuff.InArray(v.method, intercept.Match.Methods) {
		return false, nil
	}
	if intercept.Match.Uri == "" {
		return false, nil
	}
	matched, err := regexp.MatchString(intercept.Match.Uri, v.uri)
	if err != nil {
		return false, err
	}
	// matched
	return matched, nil
}
