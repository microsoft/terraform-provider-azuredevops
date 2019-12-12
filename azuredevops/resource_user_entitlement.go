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
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
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
	clients := m.(*config.AggregatedClient)
	userEntitlement, err := expandUserEntitlement(d)
	if err != nil {
		return fmt.Errorf("Error creating user entitlement: %v", err)
	}

	addedUserEntitlement, err := addUserEntitlement(clients, userEntitlement)
	if err != nil {
		return fmt.Errorf("Error creating user entitlement: %v", err)
	}

	flattenUserEntitlement(d, addedUserEntitlement)
	return resourceUserEntitlementRead(d, m)
}

func expandUserEntitlement(d *schema.ResourceData) (*memberentitlementmanagement.UserEntitlement, error) {
	origin := d.Get("origin").(string)
	originID := d.Get("origin_id").(string)
	principalName := d.Get("principal_name").(string)

	if len(originID) > 0 && len(principalName) > 0 {
		return nil, fmt.Errorf("Error both origin_id and principal_name set. You can not use both: origin_id: %s principal_name %s", originID, principalName)
	}

	if len(originID) == 0 && len(principalName) == 0 {
		return nil, fmt.Errorf("Error neither origin_id and principal_name set. Use origin_id or principal_name")
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

func addUserEntitlement(clients *config.AggregatedClient, userEntitlement *memberentitlementmanagement.UserEntitlement) (*memberentitlementmanagement.UserEntitlement, error) {
	userEntitlementsPostResponse, err := clients.MemberEntitleManagementClient.AddUserEntitlement(clients.Ctx, memberentitlementmanagement.AddUserEntitlementArgs{
		UserEntitlement: userEntitlement,
	})

	if err != nil {
		return nil, err
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
	clients := m.(*config.AggregatedClient)
	userEntitlementID := d.Id()
	id, err := uuid.Parse(userEntitlementID)
	if err != nil {
		return fmt.Errorf("Error parsing UserEntitlementID: %s. %v", userEntitlementID, err)
	}

	userEntitlement, err := readUserEntitlement(clients, &id)

	if err != nil {
		return fmt.Errorf("Error reading user entitlement: %v", err)
	}

	flattenUserEntitlement(d, userEntitlement)
	return nil
}

func readUserEntitlement(clients *config.AggregatedClient, id *uuid.UUID) (*memberentitlementmanagement.UserEntitlement, error) {
	return clients.MemberEntitleManagementClient.GetUserEntitlement(clients.Ctx, memberentitlementmanagement.GetUserEntitlementArgs{
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
		return fmt.Errorf("Error parsing UserEntitlement ID. UserEntitlementID: %s. %v", userEntitlementID, err)
	}

	clients := m.(*config.AggregatedClient)

	err = clients.MemberEntitleManagementClient.DeleteUserEntitlement(m.(*config.AggregatedClient).Ctx, memberentitlementmanagement.DeleteUserEntitlementArgs{
		UserId: &id,
	})

	if err != nil {
		return fmt.Errorf("Error deleting user entitlement: %v", err)
	}

	return nil
}
