package errorutil

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func ImportAsExistsError(id string) diag.ErrorDiagnostic {
	return diag.NewErrorDiagnostic("The resource already exists", fmt.Sprintf("The resource with the Id %q already exists - to be managed via Terraform this resource needs to be imported into the State. Please see the resource documentation for for more information", id))
}
