package azuredevops

import (
	"bytes"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/licensing"
	"github.com/microsoft/azure-devops-go-api/azuredevops/memberentitlementmanagement"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
)

func resourceUserEntitlement() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserEntitlementCreate,
		Read:   resourceUserEntitlementRead,
		Delete: resourceUserEntitlementDelete,

		Schema: map[string]*schema.Schema{
			"principal_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"origin_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"origin": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "aad",
				ValidateFunc: validation.StringInSlice([]string{"aad", "ghb"}, false),
			},
			"account_license_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "express",
				ValidateFunc: validation.StringInSlice([]string{"advanced", "earlyAdopter", "express", "none", "professional", "stakeholder"}, false),
			},
			"descriptor": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceUserEntitlementCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)
	userEntitlement, err := expandUserEntitlement(d)
	if err != nil {
		return err
	}

	addedUserEntitlement, err := addUserEntitlement(clients, userEntitlement)
	if err != nil {
		return err
	}

	flattenUserEntitlement(d, addedUserEntitlement)
	return nil
}

func expandUserEntitlement(d *schema.ResourceData) (*memberentitlementmanagement.UserEntitlement, error) {
	origin := d.Get("origin").(string)
	originID := d.Get("origin_id").(string)
	principalName := d.Get("principal_name").(string)

	if len(originID) > 0 && len(principalName) > 0 {
		return nil, fmt.Errorf("Use origin_id or principal_name. You can not use both: origin_id: %s principal_name %s", originID, principalName)
	}

	if len(originID) == 0 && len(principalName) == 0 {
		return nil, fmt.Errorf("Use origin_id or principal_name")
	}

	subjectKind := "user"

	accountLicenseType, err := converter.AccountLicenseType(d.Get("account_license_type").(string))
	if err != nil {
		return nil, err
	}

	return &memberentitlementmanagement.UserEntitlement{

		AccessLevel: &licensing.AccessLevel{
			AccountLicenseType: accountLicenseType,
		},

		// TODO check if it works in both case for GitHub and AzureDevOps
		User: &graph.GraphUser{
			Origin:        &origin,
			OriginId:      &originID,
			PrincipalName: &principalName,
			SubjectKind:   &subjectKind,
		},
	}, nil
}

func flattenUserEntitlement(d *schema.ResourceData, userEntitlement *memberentitlementmanagement.UserEntitlement) {
	d.SetId(userEntitlement.Id.String())
	d.Set("descriptor", *userEntitlement.User.Descriptor)
}

func addUserEntitlement(clients *aggregatedClient, userEntitlement *memberentitlementmanagement.UserEntitlement) (*memberentitlementmanagement.UserEntitlement, error) {
	userEntitlementsPostResponse, err := clients.MemberEntitleManagementClient.AddUserEntitlement(clients.ctx, memberentitlementmanagement.AddUserEntitlementArgs{
		UserEntitlement: userEntitlement,
	})

	if err != nil {
		return nil, fmt.Errorf("Error adding user entitlement in Azure DevOps: %+v", err)
	}

	if *userEntitlementsPostResponse.IsSuccess == false {
		var buffer bytes.Buffer
		for _, e := range *userEntitlementsPostResponse.OperationResult.Errors {
			buffer.WriteString((*e.Value).(string))
		}
		return nil, fmt.Errorf("%s", buffer.String())
	}

	return userEntitlementsPostResponse.UserEntitlement, nil
}

func resourceUserEntitlementRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)
	userEntitlementID := d.Id()
	id, err := uuid.Parse(userEntitlementID)
	if err != nil {
		return fmt.Errorf("Error parsing UserEntitlement ID, got %s: %v", userEntitlementID, err)
	}

	userEntitlement, err := readUserEntitlement(clients, &id)

	if err != nil {
		return err
	}

	flattenUserEntitlement(d, userEntitlement)
	return nil
}

func readUserEntitlement(clients *aggregatedClient, id *uuid.UUID) (*memberentitlementmanagement.UserEntitlement, error) {
	return clients.MemberEntitleManagementClient.GetUserEntitlement(clients.ctx, memberentitlementmanagement.GetUserEntitlementArgs{
		UserId: id,
	})
}

func resourceUserEntitlementDelete(d *schema.ResourceData, m interface{}) error {
	if d.Id() == "" {
		return nil
	}
	userEntitlementID := d.Id()
	id, err := uuid.Parse(userEntitlementID)
	if err != nil {
		return fmt.Errorf("Error parsing UserEntitlement ID, got %s: %v", userEntitlementID, err)
	}

	clients := m.(*aggregatedClient)

	err = clients.MemberEntitleManagementClient.DeleteUserEntitlement(m.(*aggregatedClient).ctx, memberentitlementmanagement.DeleteUserEntitlementArgs{
		UserId: &id,
	})

	return err
}
