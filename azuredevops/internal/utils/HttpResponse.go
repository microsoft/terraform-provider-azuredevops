package utils

import (
	"net/http"
	"strings"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
)

// ResponseWasNotFound was used for check if error is due to resource not found
func ResponseWasNotFound(err error) bool {
	// If API returns 404, resource was not found
	statusNotFound := ResponseWasStatusCode(err, http.StatusNotFound)
	if statusNotFound {
		return statusNotFound
	}

	// Some APIs return 400 BadRequest with the VS800075 error message if
	// DevOps Project doesn't exist. If parent project doesn't exist, all
	// child resources are considered "doesn't exist".
	statusBadRequest := ResponseWasStatusCode(err, http.StatusBadRequest)
	if statusBadRequest {
		return ResponseContainsStatusMessage(err, "VS800075")
	}
	return false
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

// ResponseContainsStatusMessage is used for check if error message contains specific message
func ResponseContainsStatusMessage(err error, statusMessage string) bool {
	if err == nil {
		return false
	}
	if wrapperErr, ok := err.(azuredevops.WrappedError); ok {
		if wrapperErr.Message == nil {
			return false
		}
		return strings.Contains(*wrapperErr.Message, statusMessage)
	}
	return false
}

// ResponseWasTypeKey is used to check if error has a specific TypeKey
func ResponseWasTypeKey(err error, typeKey string) bool {
	if err == nil {
		return false
	}
	if wrapperErr, ok := err.(azuredevops.WrappedError); ok {
		if wrapperErr.TypeKey != nil && *wrapperErr.TypeKey == typeKey {
			return true
		}
	}
	return false
}
