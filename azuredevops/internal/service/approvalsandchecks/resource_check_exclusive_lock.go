package approvalsandchecks

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/pipelineschecksextras"
)

var exclusiveLockId = "2EF31AD6-BAA0-403A-8B45-2CBC9B4E5563"

// ResourceCheckExclusiveLock schema and implementation for Exclusive Lock resource
func ResourceCheckExclusiveLock() *schema.Resource {
	r := genBaseCheckResource(flattenExclusiveLock, expandExclusiveLock)

	r.Schema["display_name"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "Managed by Terraform",
	}

	return r
}

func flattenExclusiveLock(d *schema.ResourceData, exclusiveLockCheck *pipelineschecksextras.CheckConfiguration, projectID string) error {
	err := doBaseFlattening(d, exclusiveLockCheck, projectID)
	if err != nil {
		return err
	}

	exclusiveLockCheck.Type.Id = converter.UUID(exclusiveLockId)
	return nil
}

func expandExclusiveLock(d *schema.ResourceData) (*pipelineschecksextras.CheckConfiguration, string, error) {
	settings := make(map[string]interface{})
	settings["displayName"] = d.Get("display_name").(string)

	return doBaseExpansion(d, approvalAndCheckType.ExclusiveLock, settings, nil)
}
