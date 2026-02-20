package utils

import (
	"testing"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
)

func TestResponseContainsStatusMessage(t *testing.T) {
	cases := []struct {
		Name               string
		Error              error
		StatusMessageRegex string
		Result             bool
	}{
		{
			Name:               "ProjectNotFound",
			Error:              GetError(400, "VS800075: The project with id"),
			StatusMessageRegex: "VS800075",
			Result:             true,
		},
		{
			Name:               "RandomError",
			Error:              GetError(400, "VS800075: The project with id"),
			StatusMessageRegex: "VS800076",
			Result:             false,
		},
		{
			Name:               "MissingMessage",
			Error:              azuredevops.WrappedError{},
			StatusMessageRegex: "VS800076",
			Result:             false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			result := ResponseContainsStatusMessage(tc.Error, tc.StatusMessageRegex)

			if result != tc.Result {
				t.Errorf("ResponseContainsStatusMessage returned unexpected result.")
			}
		})
	}
}

func TestResponseWasNotFound(t *testing.T) {
	cases := []struct {
		Name   string
		Error  error
		Result bool
	}{
		{
			Name:   "NoStatus",
			Error:  azuredevops.WrappedError{},
			Result: false,
		},
		{
			Name:   "404NotFound",
			Error:  GetError(404, ""),
			Result: true,
		},
		{
			Name:   "400NotFound",
			Error:  GetError(400, "VS800075: The project with id"),
			Result: true,
		},
		{
			Name:   "400FieldNotFound",
			Error:  GetError(400, "VS402806: Work item type does not contain field"),
			Result: true,
		},
		{
			Name:   "400Different",
			Error:  GetError(400, "Some different issue"),
			Result: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			result := ResponseWasNotFound(tc.Error)

			if result != tc.Result {
				t.Errorf("ResponseWasNotFound returned unexpected result.")
			}
		})
	}
}

func GetError(statusCode int, message string) azuredevops.WrappedError {
	return azuredevops.WrappedError{
		StatusCode: &statusCode,
		Message:    &message,
	}
}
