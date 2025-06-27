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
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/accounts"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/licensing"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/memberentitlementmanagement"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
)

// ResourceServicePrincipalEntitlement schema and implementation for service principal entitlement resource
func ResourceServicePrincipalEntitlement() *schema.Resource {
	return &schema.Resource{
		Create: resourceServicePrincipalEntitlementCreate,
		Read:   resourceServicePrincipalEntitlementRead,
		Delete: resourceServicePrincipalEntitlementDelete,
		Update: resourceServicePrincipalEntitlementUpdate,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			State: importServicePrincipalEntitlement,
		},
		Schema: map[string]*schema.Schema{
			"origin_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"origin": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validation.StringIsNotWhiteSpace,
			},
			"account_license_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  licensing.AccountLicenseTypeValues.Express,
				ValidateFunc: validation.StringInSlice([]string{
					string(licensing.AccountLicenseTypeValues.Advanced),
					string(licensing.AccountLicenseTypeValues.EarlyAdopter),
					string(licensing.AccountLicenseTypeValues.Express),
					"basic",
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
			"display_name": {
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

func resourceServicePrincipalEntitlementCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	servicePrincipalEntitlement, err := expandServicePrincipalEntitlement(d)
	if err != nil {
		return fmt.Errorf("Creating service principal entitlement: %v", err)
	}

	addedServicePrincipalEntitlement, err := addServicePrincipalEntitlement(clients, servicePrincipalEntitlement)
	if err != nil {
		return fmt.Errorf("Creating service principal entitlement: %v", err)
	}

	d.SetId(addedServicePrincipalEntitlement.Id.String())
	return resourceServicePrincipalEntitlementRead(d, m)
}

func resourceServicePrincipalEntitlementRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	servicePrincipalEntitlementID := d.Id()
	id, err := uuid.Parse(servicePrincipalEntitlementID)
	if err != nil {
		return fmt.Errorf("Parsing ServicePrincipalEntitlementID: %s. %v", servicePrincipalEntitlementID, err)
	}

	servicePrincipalEntitlement, err := clients.MemberEntitleManagementClient.GetServicePrincipalEntitlement(clients.Ctx, memberentitlementmanagement.GetServicePrincipalEntitlementArgs{
		ServicePrincipalId: &id,
	})

	if err != nil {
		if utils.ResponseWasNotFound(err) || isServicePrincipalDeleted(servicePrincipalEntitlement) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Reading service principal entitlement: %v", err)
	}

	if servicePrincipalEntitlement == nil || servicePrincipalEntitlement.Id == nil ||
		(servicePrincipalEntitlement.ServicePrincipal != nil && servicePrincipalEntitlement.ServicePrincipal.IsDeletedInOrigin != nil && *servicePrincipalEntitlement.ServicePrincipal.IsDeletedInOrigin) {
		log.Println(" Service Principal has been deleted")
		d.SetId("")
		return nil
	}

	flattenServicePrincipalEntitlement(d, servicePrincipalEntitlement)
	return nil
}

func resourceServicePrincipalEntitlementUpdate(d *schema.ResourceData, m interface{}) error {
	servicePrincipalEntitlementID := d.Id()
	id, err := uuid.Parse(servicePrincipalEntitlementID)
	if err != nil {
		return fmt.Errorf("Parsing ServicePrincipalEntitlement ID. ServicePrincipalEntitlementID: %s. %v", servicePrincipalEntitlementID, err)
	}

	accountLicenseType, err := converter.AccountLicenseType(d.Get("account_license_type").(string))
	if err != nil {
		return fmt.Errorf("Convert AccountLicenseType: %v", err)
	}
	licensingSource, _ := d.GetOk("licensing_source")

	clients := m.(*client.AggregatedClient)

	patchResponse, err := clients.MemberEntitleManagementClient.UpdateServicePrincipalEntitlement(clients.Ctx,
		memberentitlementmanagement.UpdateServicePrincipalEntitlementArgs{
			ServicePrincipalId: &id,
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
		return fmt.Errorf("Updating service principal entitlement: %v", err)
	}

	if patchResponse != nil && patchResponse.IsSuccess != nil && !*patchResponse.IsSuccess {
		return fmt.Errorf("Updating service principal entitlement: %s", getServicePrincipalEntitlementAPIErrorMessage(patchResponse.OperationResults))
	}
	return resourceServicePrincipalEntitlementRead(d, m)
}

func resourceServicePrincipalEntitlementDelete(d *schema.ResourceData, m interface{}) error {
	if d.Id() == "" {
		return nil
	}

	servicePrincipalEntitlementID := d.Id()
	id, err := uuid.Parse(servicePrincipalEntitlementID)
	if err != nil {
		return fmt.Errorf("Parsing ServicePrincipalEntitlement ID. ServicePrincipalEntitlementID: %s. %v", servicePrincipalEntitlementID, err)
	}

	clients := m.(*client.AggregatedClient)

	err = clients.MemberEntitleManagementClient.DeleteServicePrincipalEntitlement(m.(*client.AggregatedClient).Ctx, memberentitlementmanagement.DeleteServicePrincipalEntitlementArgs{
		ServicePrincipalId: &id,
	})

	if err != nil {
		return fmt.Errorf("Deleting service principal entitlement: %v", err)
	}

	return nil
}

func expandServicePrincipalEntitlement(d *schema.ResourceData) (*memberentitlementmanagement.ServicePrincipalEntitlement, error) {
	return &memberentitlementmanagement.ServicePrincipalEntitlement{
		AccessLevel: &licensing.AccessLevel{
			AccountLicenseType: converter.ToPtr(licensing.AccountLicenseType(d.Get("account_license_type").(string))),
			LicensingSource:    converter.ToPtr(licensing.LicensingSource(d.Get("licensing_source").(string))),
		},
		ServicePrincipal: &graph.GraphServicePrincipal{
			Origin:      converter.ToPtr(d.Get("origin").(string)),
			OriginId:    converter.ToPtr(d.Get("origin_id").(string)),
			DisplayName: converter.ToPtr(d.Get("display_name").(string)),
			Descriptor:  converter.ToPtr(d.Get("descriptor").(string)),
			SubjectKind: converter.String("servicePrincipal"),
		},
	}, nil
}

func flattenServicePrincipalEntitlement(d *schema.ResourceData, servicePrincipalEntitlement *memberentitlementmanagement.ServicePrincipalEntitlement) {
	if servicePrincipalEntitlement != nil {
		if servicePrincipalEntitlement.ServicePrincipal != nil {
			d.Set("origin", *servicePrincipalEntitlement.ServicePrincipal.Origin)
			if servicePrincipalEntitlement.ServicePrincipal.OriginId != nil {
				d.Set("origin_id", *servicePrincipalEntitlement.ServicePrincipal.OriginId)
			}
			d.Set("display_name", *servicePrincipalEntitlement.ServicePrincipal.DisplayName)
			d.Set("descriptor", *servicePrincipalEntitlement.ServicePrincipal.Descriptor)
		}
		if servicePrincipalEntitlement.AccessLevel != nil {
			d.Set("account_license_type", string(*servicePrincipalEntitlement.AccessLevel.AccountLicenseType))
			d.Set("licensing_source", *servicePrincipalEntitlement.AccessLevel.LicensingSource)
		}
	}
}

func addServicePrincipalEntitlement(clients *client.AggregatedClient, servicePrincipalEntitlement *memberentitlementmanagement.ServicePrincipalEntitlement) (*memberentitlementmanagement.ServicePrincipalEntitlement, error) {
	servicePrincipalEntitlementsPostResponse, err := clients.MemberEntitleManagementClient.AddServicePrincipalEntitlement(clients.Ctx, memberentitlementmanagement.AddServicePrincipalEntitlementArgs{
		ServicePrincipalEntitlement: servicePrincipalEntitlement,
	})

	if err != nil {
		return nil, err
	}

	if !*servicePrincipalEntitlementsPostResponse.IsSuccess {
		opResults := []memberentitlementmanagement.ServicePrincipalEntitlementOperationResult{}
		if servicePrincipalEntitlementsPostResponse.OperationResult != nil {
			opResults = append(opResults, *servicePrincipalEntitlementsPostResponse.OperationResult)
		}
		return nil, fmt.Errorf("Adding service principal entitlement: %s", getServicePrincipalEntitlementAPIErrorMessage(&opResults))
	}

	return servicePrincipalEntitlementsPostResponse.ServicePrincipalEntitlement, nil
}

func importServicePrincipalEntitlement(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	servicePrincipalEntitlementId := d.Id()
	id, err := uuid.Parse(servicePrincipalEntitlementId)

	if err != nil {
		return nil, fmt.Errorf("Only UUID values can used for import [%s]", servicePrincipalEntitlementId)
	}

	clients := m.(*client.AggregatedClient)
	resp, err := clients.MemberEntitleManagementClient.GetServicePrincipalEntitlement(clients.Ctx, memberentitlementmanagement.GetServicePrincipalEntitlementArgs{
		ServicePrincipalId: &id,
	})
	if err != nil {
		return nil, fmt.Errorf("Getting the service principal entitlement with supplied id %s: %s", servicePrincipalEntitlementId, err)
	}

	if resp == nil || resp.Id == nil {
		return nil, fmt.Errorf("Service Principal entitlement with ID: %s not found", servicePrincipalEntitlementId)
	}

	d.SetId((*resp).Id.String())

	return []*schema.ResourceData{d}, nil
}

func getServicePrincipalEntitlementAPIErrorMessage(operationResults *[]memberentitlementmanagement.ServicePrincipalEntitlementOperationResult) string {
	errMsg := "Unknown API error"
	if operationResults != nil && len(*operationResults) > 0 {
		errMsg = linq.From(*operationResults).
			Where(func(elem interface{}) bool {
				ueo := elem.(memberentitlementmanagement.ServicePrincipalEntitlementOperationResult)
				return !*ueo.IsSuccess
			}).
			SelectMany(func(elem interface{}) linq.Query {
				ueo := elem.(memberentitlementmanagement.ServicePrincipalEntitlementOperationResult)
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

func isServicePrincipalDeleted(servicePrincipalEntitlement *memberentitlementmanagement.ServicePrincipalEntitlement) bool {
	if servicePrincipalEntitlement == nil {
		return true
	}

	return *servicePrincipalEntitlement.AccessLevel.Status == accounts.AccountUserStatusValues.Deleted ||
		*servicePrincipalEntitlement.AccessLevel.Status == accounts.AccountUserStatusValues.None
}
