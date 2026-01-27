package errorutil

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func DiagToError(d diag.Diagnostic) error {
	if d.Severity() != diag.SeverityError {
		return nil
	}
	return fmt.Errorf("%s: %s", d.Summary(), d.Detail())
}

func DiagsToError(diags diag.Diagnostics) error {
	var errs error

	for _, ediag := range diags.Errors() {
		errs = multierror.Append(errs, DiagToError(ediag))
	}
	return errs
}
