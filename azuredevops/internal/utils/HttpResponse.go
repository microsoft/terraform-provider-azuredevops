package utils

import (
	"net/http"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
)

// ResponseWasNotFound was used for check if error status was 404
func ResponseWasNotFound(err error) bool {
	return ResponseWasStatusCode(err, http.StatusNotFound)
}

// ResponseWasStatusCode was used for check if error status code was specific http status code
func ResponseWasStatusCode(err error, statusCode int) bool {
	if err == nil {
		return false
	}
	if wrapperErr, ok := err.(azuredevops.WrappedError); ok {
		if wrapperErr.StatusCode != nil && *wrapperErr.StatusCode == statusCode {
			return true
		}
	}
	return false
}
