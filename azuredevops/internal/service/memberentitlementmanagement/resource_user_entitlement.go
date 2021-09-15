package memberentitlementmanagement

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/licensing"
	"github.com/microsoft/azure-devops-go-api/azuredevops/memberentitlementmanagement"
	"github.com/microsoft/azure-devops-go-api/azuredevops/webapi"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
)

var (
	configurationKeys = []string{
		"origin_id",
		"origin",
		"principal_name",
	}
)

// ResourceUserEntitlement schema and implementation for user entitlement resource
func ResourceUserEntitlement() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserEntitlementCreate,
		Read:   resourceUserEntitlementRead,
		Delete: resourceUserEntitlementDelete,
		Update: resourceUserEntitlementUpdate,
		Importer: &schema.ResourceImporter{
			State: importUserEntitlement,
		},
		Schema: map[string]*schema.Schema{
			"principal_name": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				ConflictsWith:    []string{"origin_id", "origin"},
				AtLeastOneOf:     configurationKeys,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validation.StringIsNotWhiteSpace,
			},
			"origin_id": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"principal_name"},
				AtLeastOneOf:  configurationKeys,
				ValidateFunc:  validation.StringIsNotWhiteSpace,
			},
			"origin": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"principal_name"},
				AtLeastOneOf:  configurationKeys,
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

func resourceUserEntitlementCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	userEntitlement, err := expandUserEntitlement(d)
	if err != nil {
		return fmt.Errorf("Creating user entitlement: %v", err)
	}

	addedUserEntitlement, err := addUserEntitlement(clients, userEntitlement)
	if err != nil {
		return fmt.Errorf("Creating user entitlement: %v", err)
	}

	flattenUserEntitlement(d, addedUserEntitlement)
	return resourceUserEntitlementRead(d, m)
}

func resourceUserEntitlementRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	userEntitlementID := d.Id()
	id, err := uuid.Parse(userEntitlementID)
	if err != nil {
		return fmt.Errorf("Error parsing UserEntitlementID: %s. %v", userEntitlementID, err)
	}

	userEntitlement, err := readUserEntitlement(clients, &id)

	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading user entitlement: %v", err)
	}

	flattenUserEntitlement(d, userEntitlement)
	return nil
}

func expandUserEntitlement(d *schema.ResourceData) (*memberentitlementmanagement.UserEntitlement, error) {
	origin := d.Get("origin").(string)
	originID := d.Get("origin_id").(string)
	principalName := d.Get("principal_name").(string)

	if len(originID) > 0 && len(principalName) > 0 {
		return nil, fmt.Errorf("Both origin_id and principal_name set. You can not use both: origin_id: %s principal_name %s", originID, principalName)
	}

	if len(originID) == 0 && len(principalName) == 0 {
		return nil, fmt.Errorf("Neither origin_id and principal_name set. Use origin_id or principal_name")
	}

	if len(originID) > 0 && len(origin) == 0 {
		return nil, fmt.Errorf("Origin_id requires an origin to be set")
	}

	accountLicenseType, err := converter.AccountLicenseType(d.Get("account_license_type").(string))
	if err != nil {
		return nil, err
	}
	licensingSource, err := converter.AccountLicensingSource(d.Get("licensing_source").(string))
	if err != nil {
		return nil, err
	}

	return &memberentitlementmanagement.UserEntitlement{

		AccessLevel: &licensing.AccessLevel{
			AccountLicenseType: accountLicenseType,
			LicensingSource:    licensingSource,
		},

		// TODO check if it works in both case for GitHub and AzureDevOps
		User: &graph.GraphUser{
			Origin:        &origin,
			OriginId:      &originID,
			PrincipalName: &principalName,
			SubjectKind:   converter.String("user"),
		},
	}, nil
}

func flattenUserEntitlement(d *schema.ResourceData, userEntitlement *memberentitlementmanagement.UserEntitlement) {
	d.SetId(userEntitlement.Id.String())
	d.Set("descriptor", *userEntitlement.User.Descriptor)
	d.Set("origin", *userEntitlement.User.Origin)
	if userEntitlement.User.OriginId != nil {
		d.Set("origin_id", *userEntitlement.User.OriginId)
	}
	d.Set("principal_name", *userEntitlement.User.PrincipalName)
	d.Set("account_license_type", string(*userEntitlement.AccessLevel.AccountLicenseType))
	d.Set("licensing_source", *userEntitlement.AccessLevel.LicensingSource)
}

func addUserEntitlement(clients *client.AggregatedClient, userEntitlement *memberentitlementmanagement.UserEntitlement) (*memberentitlementmanagement.UserEntitlement, error) {
	userEntitlementsPostResponse, err := clients.MemberEntitleManagementClient.AddUserEntitlement(clients.Ctx, memberentitlementmanagement.AddUserEntitlementArgs{
		UserEntitlement: userEntitlement,
	})

	if err != nil {
		return nil, err
	}

	if !*userEntitlementsPostResponse.IsSuccess {
		opResults := []memberentitlementmanagement.UserEntitlementOperationResult{}
		if userEntitlementsPostResponse.OperationResult != nil {
			opResults = append(opResults, *userEntitlementsPostResponse.OperationResult)
		}
		return nil, fmt.Errorf("Adding user entitlement: %s", getAPIErrorMessage(&opResults))
	}

	return userEntitlementsPostResponse.UserEntitlement, nil
}

func readUserEntitlement(clients *client.AggregatedClient, id *uuid.UUID) (*memberentitlementmanagement.UserEntitlement, error) {
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

	clients := m.(*client.AggregatedClient)

	err = clients.MemberEntitleManagementClient.DeleteUserEntitlement(m.(*client.AggregatedClient).Ctx, memberentitlementmanagement.DeleteUserEntitlementArgs{
		UserId: &id,
	})

	if err != nil {
		return fmt.Errorf("Deleting user entitlement: %v", err)
	}

	return nil
}

func resourceUserEntitlementUpdate(d *schema.ResourceData, m interface{}) error {
	userEntitlementID := d.Id()
	id, err := uuid.Parse(userEntitlementID)
	if err != nil {
		return fmt.Errorf("Parsing UserEntitlement ID. UserEntitlementID: %s. %v", userEntitlementID, err)
	}

	accountLicenseType, err := converter.AccountLicenseType(d.Get("account_license_type").(string))
	if err != nil {
		return err
	}
	licensingSource, ok := d.GetOk("licensing_source")
	if !ok {
		return fmt.Errorf("Reading account licensing source for UserEntitlementID: %s", userEntitlementID)
	}

	clients := m.(*client.AggregatedClient)

	patchResponse, err := clients.MemberEntitleManagementClient.UpdateUserEntitlement(clients.Ctx,
		memberentitlementmanagement.UpdateUserEntitlementArgs{
			UserId: &id,
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
		return fmt.Errorf("Updating user entitlement: %v", err)
	}

	if !*patchResponse.IsSuccess {
		return fmt.Errorf("Updating user entitlement: %s", getAPIErrorMessage(patchResponse.OperationResults))
	}
	return resourceUserEntitlementRead(d, m)
}

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func importUserEntitlement(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	_, err := uuid.Parse(d.Id())
	if err != nil {
		upn := d.Id()
		if !emailRegexp.MatchString(upn) {
			return nil, fmt.Errorf("Only UUID and UPN values can used for import [%s]", upn)
		}

		id := ""
		clients := m.(*client.AggregatedClient)
		skip := 0
		page := 100

		for {
			memberList, err := clients.MemberEntitleManagementClient.GetUserEntitlements(clients.Ctx, memberentitlementmanagement.GetUserEntitlementsArgs{
				Top:  &page,
				Skip: &skip,
				// Using sort, so that we can predictably skip over items that we have read already
				// FIXME: should be sorted by dateCreated because LastAccessDate is not stable,
				// but currently the API fails with an error ("Value cannot be null.\r\nParameter name: keySelector")
				SortOption: converter.String("LastAccessDate Desc"),
			})
			if err != nil {
				return nil, err
			}
			for _, item := range *(*memberList).Members {
				if strings.EqualFold(*item.User.PrincipalName, upn) {
					id = item.Id.String()
					goto Found
				}
			}
			if len(*(*memberList).Members) < page {
				break
			}
			skip += page
		}

		return nil, fmt.Errorf("No entitlement found for [%s]", upn)

	Found:
		d.SetId(id)
	}
	return []*schema.ResourceData{d}, nil
}

func getAPIErrorMessage(operationResults *[]memberentitlementmanagement.UserEntitlementOperationResult) string {
	errMsg := "Unknown API error"
	if operationResults != nil && len(*operationResults) > 0 {
		errMsg = linq.From(*operationResults).
			Where(func(elem interface{}) bool {
				ueo := elem.(memberentitlementmanagement.UserEntitlementOperationResult)
				return !*ueo.IsSuccess
			}).
			SelectMany(func(elem interface{}) linq.Query {
				ueo := elem.(memberentitlementmanagement.UserEntitlementOperationResult)
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
