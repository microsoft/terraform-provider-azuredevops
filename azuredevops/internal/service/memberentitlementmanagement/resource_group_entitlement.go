package memberentitlementmanagement

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ahmetb/go-linq"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/licensing"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/memberentitlementmanagement"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
)

func ResourceGroupEntitlement() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupEntitlementCreate,
		Read:   resourceGroupEntitlementRead,
		Update: resourceGroupEntitlementUpdate,
		Delete: resourceGroupEntitlementDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			State: importGroupEntitlement,
		},
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"origin_id", "origin"},
				ExactlyOneOf:  []string{"display_name", "origin_id"},
				ValidateFunc:  validation.StringIsNotWhiteSpace,
			},
			"origin_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"display_name"},
				RequiredWith:  []string{"origin"},
				ExactlyOneOf:  []string{"display_name", "origin_id"},
				ValidateFunc:  validation.StringIsNotWhiteSpace,
			},
			"origin": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"display_name"},
				RequiredWith:  []string{"origin_id"},
				ValidateFunc:  validation.StringIsNotWhiteSpace,
			},
			"account_license_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  licensing.AccountLicenseTypeValues.Express,
				ValidateFunc: validation.StringInSlice([]string{
					string(licensing.AccountLicenseTypeValues.Advanced),
					string(licensing.AccountLicenseTypeValues.EarlyAdopter),
					string(licensing.AccountLicenseTypeValues.Express),
					string(licensing.AccountLicenseTypeValues.None),
					string(licensing.AccountLicenseTypeValues.Professional),
					string(licensing.AccountLicenseTypeValues.Stakeholder),
				}, true),
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					equalEntitlements := []string{
						string(licensing.AccountLicenseTypeValues.EarlyAdopter),
						string(licensing.AccountLicenseTypeValues.Express),
						"basic",
					}
					stringInSlice := func(v string, valid []string) bool {
						for _, str := range valid {
							if strings.EqualFold(v, str) {
								return true
							}
						}
						return false
					}
					return strings.EqualFold(old, new) ||
						(stringInSlice(old, equalEntitlements) && stringInSlice(new, equalEntitlements))
				},
			},
			"licensing_source": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(licensing.LicensingSourceValues.Account),
				ValidateFunc: validation.StringInSlice([]string{
					string(licensing.LicensingSourceValues.None),
					string(licensing.LicensingSourceValues.Account),
					string(licensing.LicensingSourceValues.Msdn),
					string(licensing.LicensingSourceValues.Profile),
					string(licensing.LicensingSourceValues.Auto),
					string(licensing.LicensingSourceValues.Trial),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"principal_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"descriptor": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGroupEntitlementCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	groupEntitlement, err := expandGroupEntitlement(d)
	if err != nil {
		return fmt.Errorf("Creating group entitlement: %v", err)
	}

	addedGroupEntitlement, err := addGroupEntitlement(clients, groupEntitlement)
	if err != nil {
		return fmt.Errorf("Creating group entitlement: %v", err)
	}

	d.SetId(addedGroupEntitlement.Id.String())
	return resourceGroupEntitlementRead(d, m)
}

func resourceGroupEntitlementRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	groupEntitlementID := d.Id()
	id, err := uuid.Parse(groupEntitlementID)
	if err != nil {
		return fmt.Errorf("Parsing GroupEntitlementID: %s. %v", groupEntitlementID, err)
	}
	groupEntitlement, err := clients.MemberEntitleManagementClient.GetGroupEntitlement(clients.Ctx, memberentitlementmanagement.GetGroupEntitlementArgs{
		GroupId: &id,
	})

	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading group entitlement: %v", err)
	}

	if groupEntitlement == nil || groupEntitlement.Id == nil ||
		(groupEntitlement.Group != nil && groupEntitlement.Group.IsDeleted != nil && *groupEntitlement.Group.IsDeleted) {
		log.Println(" Group has been deleted")
		d.SetId("")
		return nil
	}

	flattenGroupEntitlement(d, groupEntitlement)
	return nil
}

func resourceGroupEntitlementUpdate(d *schema.ResourceData, m interface{}) error {
	groupEntitlementID := d.Id()
	id, err := uuid.Parse(groupEntitlementID)
	if err != nil {
		return fmt.Errorf("Parsing GroupEntitlement ID. GroupEntitlementID: %s. %v", groupEntitlementID, err)
	}

	accountLicenseType, err := converter.AccountLicenseType(d.Get("account_license_type").(string))
	if err != nil {
		return err
	}
	licensingSource, ok := d.GetOk("licensing_source")
	if !ok {
		return fmt.Errorf("Reading account licensing source for GroupEntitlementID: %s", groupEntitlementID)
	}

	clients := m.(*client.AggregatedClient)

	patchResponse, err := clients.MemberEntitleManagementClient.UpdateGroupEntitlement(clients.Ctx,
		memberentitlementmanagement.UpdateGroupEntitlementArgs{
			GroupId: &id,
			Document: &[]webapi.JsonPatchOperation{
				{
					Op:   &webapi.OperationValues.Replace,
					From: nil,
					Path: converter.String("/accessLevel"),
					Value: struct {
						AccountLicenseType string `json:"accountLicenseType"`
						LicensingSource    string `json:"licensingSource"`
					}{
						string(*accountLicenseType),
						licensingSource.(string),
					},
				},
			},
		})

	if err != nil {
		return fmt.Errorf("Updating group entitlement: %v", err)
	}

	result := *patchResponse.Results

	if !*result[0].IsSuccess {
		return fmt.Errorf("Updating group entitlement: %s", getGroupEntitlementAPIErrorMessage(&result))
	}
	return resourceGroupEntitlementRead(d, m)
}

func resourceGroupEntitlementDelete(d *schema.ResourceData, m interface{}) error {
	if d.Id() == "" {
		return nil
	}

	groupEntitlementID := d.Id()
	id, err := uuid.Parse(groupEntitlementID)
	if err != nil {
		return fmt.Errorf("Parsing GroupEntitlement ID. GroupEntitlementID: %s. %v", groupEntitlementID, err)
	}

	clients := m.(*client.AggregatedClient)

	_, err = clients.MemberEntitleManagementClient.DeleteGroupEntitlement(m.(*client.AggregatedClient).Ctx, memberentitlementmanagement.DeleteGroupEntitlementArgs{
		GroupId: &id,
	})

	if err != nil {
		return fmt.Errorf("Deleting group entitlement: %v", err)
	}

	// Also delete the org wise group if the group is Azure DevOps local, meaning
	// most likely the local group was created by this resource
	origin := d.Get("origin")
	if origin == "vsts" {
		err = clients.GraphClient.DeleteGroup(clients.Ctx, graph.DeleteGroupArgs{
			GroupDescriptor: converter.String(d.Get("descriptor").(string)),
		})

		if err != nil {
			return fmt.Errorf("Deleting Azure DevOps local group: %v", err)
		}
	}

	return nil
}

func importGroupEntitlement(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	upn := d.Id()
	id, err := uuid.Parse(upn)

	if err != nil {
		return nil, fmt.Errorf("Only UUID values can used for import [%s]", upn)
	}

	clients := m.(*client.AggregatedClient)
	resp, err := clients.MemberEntitleManagementClient.GetGroupEntitlement(clients.Ctx, memberentitlementmanagement.GetGroupEntitlementArgs{
		GroupId: &id,
	})
	if err != nil {
		return nil, fmt.Errorf("Getting the group entitlement with supplied id %s: %s", upn, err)
	}

	if resp == nil || resp.Id == nil {
		return nil, fmt.Errorf("Group entitlement with ID: %s not found", upn)
	}

	d.SetId((*resp).Id.String())

	return []*schema.ResourceData{d}, nil
}

func flattenGroupEntitlement(d *schema.ResourceData, groupEntitlement *memberentitlementmanagement.GroupEntitlement) {
	d.Set("descriptor", *groupEntitlement.Group.Descriptor)
	d.Set("origin", *groupEntitlement.Group.Origin)
	d.Set("principal_name", *groupEntitlement.Group.PrincipalName)
	if groupEntitlement.Group.OriginId != nil {
		d.Set("origin_id", *groupEntitlement.Group.OriginId)
	}
	d.Set("display_name", *groupEntitlement.Group.DisplayName)
	d.Set("account_license_type", string(*groupEntitlement.LicenseRule.AccountLicenseType))
	d.Set("licensing_source", *groupEntitlement.LicenseRule.LicensingSource)
}

func expandGroupEntitlement(d *schema.ResourceData) (*memberentitlementmanagement.GroupEntitlement, error) {
	origin := d.Get("origin").(string)
	originID := d.Get("origin_id").(string)
	displayName := d.Get("display_name").(string)

	accountLicenseType, err := converter.AccountLicenseType(d.Get("account_license_type").(string))
	if err != nil {
		return nil, err
	}
	licensingSource, err := converter.AccountLicensingSource(d.Get("licensing_source").(string))
	if err != nil {
		return nil, err
	}

	return &memberentitlementmanagement.GroupEntitlement{
		LicenseRule: &licensing.AccessLevel{
			AccountLicenseType: accountLicenseType,
			LicensingSource:    licensingSource,
		},

		Group: &graph.GraphGroup{
			Origin:      &origin,
			OriginId:    &originID,
			DisplayName: &displayName,
			SubjectKind: converter.String("group"),
		},
	}, nil
}

func addGroupEntitlement(clients *client.AggregatedClient, groupEntitlement *memberentitlementmanagement.GroupEntitlement) (*memberentitlementmanagement.GroupEntitlement, error) {
	groupEntitlementsPostResponse, err := clients.MemberEntitleManagementClient.AddGroupEntitlement(clients.Ctx, memberentitlementmanagement.AddGroupEntitlementArgs{
		GroupEntitlement: groupEntitlement,
	})

	if err != nil {
		return nil, err
	}

	result := *groupEntitlementsPostResponse.Results

	if !*result[0].IsSuccess {
		opResults := []memberentitlementmanagement.GroupOperationResult{}
		if result[0].Errors != nil {
			opResults = append(opResults, result[0])
		}
		return nil, fmt.Errorf("Adding group entitlement: %s", getGroupEntitlementAPIErrorMessage(&opResults))
	}

	return result[0].Result, nil
}

func getGroupEntitlementAPIErrorMessage(operationResults *[]memberentitlementmanagement.GroupOperationResult) string {
	errMsg := "Unknown API error"
	if operationResults != nil && len(*operationResults) > 0 {
		errMsg = linq.From(*operationResults).
			Where(func(elem interface{}) bool {
				ueo := elem.(memberentitlementmanagement.GroupOperationResult)
				return !*ueo.IsSuccess
			}).
			SelectMany(func(elem interface{}) linq.Query {
				ueo := elem.(memberentitlementmanagement.GroupOperationResult)
				if ueo.Errors == nil {
					key := interface{}("0000")
					value := interface{}("Unknown API error")
					return linq.From([]azuredevops.KeyValuePair{
						{
							Key:   &key,
							Value: &value,
						},
					})
				}
				return linq.From(*ueo.Errors)
			}).
			SelectT(func(err azuredevops.KeyValuePair) string {
				return fmt.Sprintf("(%v) %s", *err.Key, *err.Value)
			}).
			AggregateT(func(agg string, elem string) string {
				return agg + "\n" + elem
			}).(string)
	}
	return errMsg
}
