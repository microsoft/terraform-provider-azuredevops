package workitemtrackingprocess

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func convertColorToApi(d *schema.ResourceData) *string {
	return converter.String(
		strings.ReplaceAll(d.Get("color").(string), "#", ""))
}

func convertColorToResource(apiFormattedColor string) string {
	return fmt.Sprintf("#%s", apiFormattedColor)
}
