package memberentitlementmanagement

import (
	"fmt"
	"log"
	"strings"

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

func ResourceServicePrincipalEntitlement() *schema.Resource {
	return &schema.Resource{
		Create: resourceServicePrincipalEntitlementCreate,
		Read:   resourceServicePrincipalEntitlementRead,
		Update: resourceServicePrincipalEntitlementUpdate,
		Delete: resourceServicePrincipalEntitlementDelete,
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
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      string("aad"),
				ValidateFunc: validation.StringIsNotWhiteSpace,
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
		return fmt.Errorf("Error parsing ServicePrincipalEntitlementID: %s. %v", servicePrincipalEntitlementID, err)
	}
	servicePrincipalEntitlement, err := clients.MemberEntitleManagementClient.GetServicePrincipalEntitlement(clients.Ctx, memberentitlementmanagement.GetServicePrincipalEntitlementArgs{
		ServicePrincipalId: &id,
	})

	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf(" reading service principal entitlement: %v", err)
	}

	if servicePrincipalEntitlement == nil || servicePrincipalEntitlement.Id == nil {
		log.Println(" Service principal entitlement has been deleted")
		d.SetId("")
		return nil
	}

	flattenServicePrincipalEntitlement(d, servicePrincipalEntitlement)
	return nil
}

func resourceServicePrincipalEntitlementDelete(d *schema.ResourceData, m interface{}) error {
	if d.Id() == "" {
		return nil
	}

	servicePrincipalEntitlementID := d.Id()
	id, err := uuid.Parse(servicePrincipalEntitlementID)
	if err != nil {
		return fmt.Errorf("Error parsing ServicePrincipalEntitlement ID. ServicePrincipalEntitlementID: %s. %v", servicePrincipalEntitlementID, err)
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

func resourceServicePrincipalEntitlementUpdate(d *schema.ResourceData, m interface{}) error {
	servicePrincipalEntitlementID := d.Id()
	id, err := uuid.Parse(servicePrincipalEntitlementID)
	if err != nil {
		return fmt.Errorf("Parsing ServicePrincipalEntitlement ID. ServicePrincipalEntitlementID: %s. %v", servicePrincipalEntitlementID, err)
	}

	accountLicenseType, err := converter.AccountLicenseType(d.Get("account_license_type").(string))
	if err != nil {
		return err
	}
	licensingSource, ok := d.GetOk("licensing_source")
	if !ok {
		return fmt.Errorf("Reading account licensing source for ServicePrincipalEntitlementID: %s", servicePrincipalEntitlementID)
	}

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

	if !*patchResponse.IsSuccess {
		return fmt.Errorf("Updating service principal entitlement: %s", getServicePrincipalEntitlementAPIErrorMessage(patchResponse.OperationResults))
	}
	return resourceServicePrincipalEntitlementRead(d, m)
}

func importServicePrincipalEntitlement(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	upn := d.Id()
	id, err := uuid.Parse(upn)

	if err != nil {
		return nil, fmt.Errorf("Only UUID values can used for import [%s]", upn)
	}

	clients := m.(*client.AggregatedClient)
	result, err := clients.MemberEntitleManagementClient.GetServicePrincipalEntitlement(clients.Ctx, memberentitlementmanagement.GetServicePrincipalEntitlementArgs{
		ServicePrincipalId: &id,
	})
	if err != nil {
		return nil, fmt.Errorf("Error getting the service principal entitlement with supplied id %s: %s", upn, err)
	}

	d.SetId((*result).Id.String())

	return []*schema.ResourceData{d}, nil
}

func flattenServicePrincipalEntitlement(d *schema.ResourceData, servicePrincipalEntitlement *memberentitlementmanagement.ServicePrincipalEntitlement) {
	d.SetId(servicePrincipalEntitlement.Id.String())
	d.Set("descriptor", *servicePrincipalEntitlement.ServicePrincipal.Descriptor)
	d.Set("origin", *servicePrincipalEntitlement.ServicePrincipal.Origin)
	d.Set("origin_id", *servicePrincipalEntitlement.ServicePrincipal.OriginId)
	d.Set("account_license_type", string(*servicePrincipalEntitlement.AccessLevel.AccountLicenseType))
	d.Set("licensing_source", *servicePrincipalEntitlement.AccessLevel.LicensingSource)
}

func expandServicePrincipalEntitlement(d *schema.ResourceData) (*memberentitlementmanagement.ServicePrincipalEntitlement, error) {
	origin := d.Get("origin").(string)
	originID := d.Get("origin_id").(string)

	accountLicenseType, err := converter.AccountLicenseType(d.Get("account_license_type").(string))
	if err != nil {
		return nil, err
	}
	licensingSource, err := converter.AccountLicensingSource(d.Get("licensing_source").(string))
	if err != nil {
		return nil, err
	}

	return &memberentitlementmanagement.ServicePrincipalEntitlement{
		AccessLevel: &licensing.AccessLevel{
			AccountLicenseType: accountLicenseType,
			LicensingSource:    licensingSource,
		},

		ServicePrincipal: &graph.GraphServicePrincipal{
			Origin:      &origin,
			OriginId:    &originID,
			SubjectKind: converter.String("servicePrincipal"),
		},
	}, nil
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
