package errorutil

import (
	"net/http"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/pointer"
)

// WasBadRequest returns true if the error is an WrappedError and has a status code of BadRequest
func WasBadRequest(err error) bool {
	return WasStatusCode(err, http.StatusBadRequest)
}

// WasConflict returns true if the error is an WrappedError and has a status code of Conflict
func WasConflict(err error) bool {
	return WasStatusCode(err, http.StatusConflict)
}

// WasForbidden returns true if the error is an WrappedError and has a status code of Forbidden
func WasForbidden(err error) bool {
	return WasStatusCode(err, http.StatusForbidden)
}

// WasNotFound returns true if the error is an WrappedError and has a status code of NotFound
func WasNotFound(err error) bool {
	return WasStatusCode(err, http.StatusNotFound)
}

// WasStatusCode returns true if the error is an WrappedError and matches the Status Code
// It's recommended to use WasBadRequest/WasConflict/WasNotFound where possible instead
func WasStatusCode(err error, statusCode int) bool {
	if err == nil {
		return false
	}
	werr, ok := err.(azuredevops.WrappedError)
	if !ok {
		return false
	}
	return pointer.To(werr.StatusCode) == statusCode
}
