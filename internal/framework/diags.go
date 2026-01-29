package framework

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/errorutil"
)

var _ diag.Diagnostic = DiagSdkError{}

// DiagSdkError represents an ADO API Error.
type DiagSdkError struct {
	summary string
	err     error
}

func (e DiagSdkError) Inner() error {
	return e.err
}

func NewDiagSdkError(summary string, err error) DiagSdkError {
	return DiagSdkError{summary, err}
}

func NewDiagSdkErrorWithCode(summary string, code int) DiagSdkError {
	return DiagSdkError{summary, azuredevops.WrappedError{StatusCode: &code}}
}

func IsDiagResourceNotFound(o diag.Diagnostic) bool {
	err, ok := o.(DiagSdkError)
	if !ok {
		return false
	}
	return errorutil.WasNotFound(err.err)
}

// Summary implements diag.Diagnostic.
func (d DiagSdkError) Summary() string {
	return "AzureDevops SDK call: " + d.summary
}

// Detail implements diag.Diagnostic.
func (d DiagSdkError) Detail() string {
	if e, ok := d.err.(azuredevops.WrappedError); ok {
		if b, err := json.Marshal(e); err == nil {
			return string(b)
		}
	}
	return d.err.Error()
}

// Equal implements diag.Diagnostic.
func (d DiagSdkError) Equal(o diag.Diagnostic) bool {
	if do, ok := o.(DiagSdkError); ok {
		return d == do
	}
	return false
}

// Severity implements diag.Diagnostic.
func (d DiagSdkError) Severity() diag.Severity {
	return diag.SeverityError
}
