package approvalsandchecks

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/pipelineschecksextras"
)

// ResourceCheckExclusiveLock schema and implementation for exclusive lock resource
func ResourceCheckExclusiveLock() *schema.Resource {
	return genBaseCheckResource(flattenExclusiveLock, expandExclusiveLock)
}

func flattenExclusiveLock(d *schema.ResourceData, exclusiveLockCheck *pipelineschecksextras.CheckConfiguration, projectID string) error {
	// TODO: implement this function
	return nil
}

func expandExclusiveLock(d *schema.ResourceData) (*pipelineschecksextras.CheckConfiguration, string, error) {
	// TODO: implement this function
	return nil, "", nil
}
