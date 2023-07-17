package approvalsandchecks

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/pipelineschecksextras"
)

// ResourceCheckExclusiveLock schema and implementation for Exclusive Lock resource
func ResourceCheckExclusiveLock() *schema.Resource {
	r := genBaseCheckResource(flattenExclusiveLock, expandExclusiveLock)

	r.Schema["timeout"] = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Default:  43200,
	}

	return r
}

func flattenExclusiveLock(d *schema.ResourceData, check *pipelineschecksextras.CheckConfiguration, projectID string) error {
	err := doBaseFlattening(d, check, projectID)
	if err != nil {
		return err
	}

	if check.Settings == nil {
		return fmt.Errorf("settings nil")
	}

	check.Type = approvalAndCheckType.ExclusiveLock

	if check.Timeout != nil {
		d.Set("timeout", *check.Timeout)
	}

	return nil
}

func expandExclusiveLock(d *schema.ResourceData) (*pipelineschecksextras.CheckConfiguration, string, error) {
	settings := make(map[string]interface{})

	return doBaseExpansion(d, approvalAndCheckType.ExclusiveLock, settings, converter.ToPtr(d.Get("timeout").(int)))
}
