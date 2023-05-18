package repository

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/policy"
)

func ResourceRepositoryPolicyCheckCredentials() *schema.Resource {
	resource := genBasePolicyResource(&policyCrudArgs{
		FlattenFunc: checkCredentialsFlattenFunc,
		ExpandFunc:  checkCredentialsExpandFunc,
		PolicyType:  CheckCredentials,
	})
	resource.DeprecationMessage = "This Repository has been deprecated and cannot create this policy anymore." //TODO
	return resource
}

func checkCredentialsFlattenFunc(d *schema.ResourceData, policyConfig *policy.PolicyConfiguration, projectID *string) error {
	err := baseFlattenFunc(d, policyConfig, projectID)
	if err != nil {
		return err
	}
	return nil
}

func checkCredentialsExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	policyConfig, projectID, err := baseExpandFunc(d, typeID)
	if err != nil {
		return nil, nil, err
	}

	return policyConfig, projectID, nil
}
