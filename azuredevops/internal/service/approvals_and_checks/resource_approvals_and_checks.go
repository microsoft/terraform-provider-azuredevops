package serviceendpoint

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/pipelineschecks"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/pipelinestaskcheck"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

var branchControlCheckType = pipelineschecks.CheckType{
	Id: converter.UUID("fe1de3ee-a436-41b4-bb20-f6eb4cb879a7"),
}

var branchProtectionDefName = "evaluatebranchProtection"
var branchProtetctinoDefVersion = "0.0.1"

type TaskCheckDefinitionReference struct {
	Id      *uuid.UUID `json:"id,omitempty"`
	Name    *string    `json:"name,omitempty"`
	Version *string    `json:"version,omitempty"`
}

var evaluateBranchProtectionDef = pipelinestaskcheck.TaskCheckDefinitionReference{
	Id:      converter.UUID("86b05a0c-73e6-4f7d-b3cf-e38f3b39a75b"),
	Name:    &branchProtectionDefName,
	Version: &branchProtetctinoDefVersion,
}

// ResourceBranchControlCheck schema and implementation for branch control checks
func ResourceServiceEndpointCheckBranchControl() *schema.Resource {
	return &schema.Resource{
		Create:   resourceBranchControlCheckCreate,
		Read:     resourceBranchControlCheckRead,
		Update:   resourceBranchControlCheckUpdate,
		Delete:   resourceBranchControlCheckDelete,
		Importer: tfhelper.ImportProjectQualifiedResource(),
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"allowed_branches": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  `*`,
			},
			"verify_branch_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ignore_unknown_protection_status": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceBranchControlCheckCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	configuration, projectID, err := expandBranchControlCheck(d)
	if err != nil {
		return fmt.Errorf(" failed in expandBranchControlCheck. Error: %+v", err)
	}

	createdBranchControlCheck, err := clients.V5PipelinesChecksClient.AddCheckConfiguration(clients.Ctx, pipelineschecks.AddCheckConfigurationArgs{
		Project:       &projectID,
		Configuration: configuration,
	})
	if err != nil {
		return fmt.Errorf(" failed creating Brach Control Check, project ID: %s. Error: %+v", projectID, err)
	}

	flattenBranchControlCheck(d, createdBranchControlCheck, projectID)
	return resourceBranchControlCheckRead(d, m)
}

func resourceBranchControlCheckRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectID, branchControlCheckID, err := tfhelper.ParseProjectIDAndResourceID(d)
	if err != nil {
		return err
	}

	branchControlCheck, err := GetCheckConfiguration(&clients.V5PipelinesChecksClient, clients.Ctx, pipelineschecks.GetCheckConfigurationArgs{
		Project: &projectID,
		Id:      &branchControlCheckID,
	})

	if err != nil {
		if utils.ResponseWasNotFound(err) || strings.Contains(err.Error(), "does not exist.") {
			d.SetId("")
			return nil
		}
		return err
	}

	err = flattenBranchControlCheck(d, branchControlCheck, projectID)
	if err != nil {
		return err
	}
	return nil
}

func resourceBranchControlCheckUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	branchControlCheck, projectID, err := expandBranchControlCheck(d)
	if err != nil {
		return err
	}

	updatedBranchControlCheck, err := clients.V5PipelinesChecksClient.UpdateCheckConfiguration(clients.Ctx,
		pipelineschecks.UpdateCheckConfigurationArgs{
			Project:       &projectID,
			Configuration: branchControlCheck,
			Id:            branchControlCheck.Id,
		})

	if err != nil {
		return err
	}

	flattenBranchControlCheck(d, updatedBranchControlCheck, projectID)
	return resourceBranchControlCheckRead(d, m)
}

func resourceBranchControlCheckDelete(d *schema.ResourceData, m interface{}) error {
	if strings.EqualFold(d.Id(), "") {
		return nil
	}

	clients := m.(*client.AggregatedClient)
	projectID, BranchControlCheckID, err := tfhelper.ParseProjectIDAndResourceID(d)
	if err != nil {
		return err
	}

	err = clients.V5PipelinesChecksClient.DeleteCheckConfiguration(m.(*client.AggregatedClient).Ctx,
		pipelineschecks.DeleteCheckConfigurationArgs{
			Project: &projectID,
			Id:      &BranchControlCheckID,
		})

	return err
}

func flattenBranchControlCheck(d *schema.ResourceData, branchControlCheck *pipelineschecks.CheckConfiguration, projectID string) error {
	d.SetId(fmt.Sprintf("%d", *branchControlCheck.Id))

	// TODO verify definitionRef

	d.Set("project_id", projectID)
	d.Set("endpoint_id", branchControlCheck.Resource.Id)

	if branchControlCheck.Settings == nil {
		return fmt.Errorf("Settings nil")
	}

	if DisplayName, found := branchControlCheck.Settings.(map[string]interface{})["displayName"]; found {
		d.Set("display_name", DisplayName.(string))
	} else {
		return fmt.Errorf("displayName setting not found")
	}

	var inputs map[string]interface{}

	if inputMap, found := branchControlCheck.Settings.(map[string]interface{})["inputs"]; found {
		inputs = inputMap.(map[string]interface{})
	} else {
		return fmt.Errorf("inputs not found")
	}

	if AllowedBranches, found := inputs["allowedBranches"]; found {
		d.Set("allowed_branches", AllowedBranches)
	} else {
		return fmt.Errorf("allowedBranches input not found")
	}

	if verifyBranchProtection, found := inputs["ensureProtectionOfBranch"]; found {
		d.Set("verify_branch_protection", verifyBranchProtection)
	} else {
		return fmt.Errorf("ensureProtectionOfBranch input not found")
	}

	if ignoreUnknownProtectionStatus, found := inputs["allowUnknownStatusBranch"]; found {
		d.Set("ignore_unknown_protection_status", ignoreUnknownProtectionStatus)
	} else {
		return fmt.Errorf("allowUnknownStatusBranch input not found")
	}

	return nil
}

func expandBranchControlCheck(d *schema.ResourceData) (*pipelineschecks.CheckConfiguration, string, error) {
	projectID := d.Get("project_id").(string)
	endpointID := d.Get("endpoint_id").(string)
	displayName := d.Get("display_name").(string)
	allowedBranchs := d.Get("allowed_branches").(string)
	verifyBranchProtection := d.Get("verify_branch_protection").(bool)
	ignoreUnknownProtectionStatus := d.Get("ignore_unknown_protection_status").(bool)

	endpointType := "endpoint"
	endpointResource := pipelineschecks.Resource{
		Id:   &endpointID,
		Type: &endpointType,
	}

	branchControlInputs := map[string]string{
		"allowedBranches":          allowedBranchs,
		"ensureProtectionOfBranch": strconv.FormatBool(verifyBranchProtection),
		"allowUnknownStatusBranch": strconv.FormatBool(ignoreUnknownProtectionStatus),
	}

	branchControlCheckSettings := pipelinestaskcheck.TaskCheckConfig{
		DefinitionRef: &evaluateBranchProtectionDef,
		DisplayName:   &displayName,
		Inputs:        &branchControlInputs,
	}

	branchControlCheck := pipelineschecks.CheckConfiguration{
		Type:     &branchControlCheckType,
		Settings: &branchControlCheckSettings,
		Resource: &endpointResource,
	}

	if d.Id() != "" {
		branchControlCheckId, err := strconv.Atoi(d.Id())
		if err != nil {
			return nil, "", fmt.Errorf("Error parsing branch control ID: (%+v)", err)
		}
		branchControlCheck.Id = &branchControlCheckId
	}

	return &branchControlCheck, projectID, nil
}

// Neither the v5 or v6 client are complete solutions for managing checks & approvals.
// The v5 client is missing the expand option, but the v6 client is missing the Settings parameter
// in the configuration. Copying the v6 `GetCheckConfiguration` here but with the v5 types to reconcil
// the discrepancies.
func GetCheckConfiguration(client *pipelineschecks.Client, ctx context.Context, args pipelineschecks.GetCheckConfigurationArgs) (*pipelineschecks.CheckConfiguration, error) {
	fullClient := (*client).(*pipelineschecks.ClientImpl)
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Id == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = strconv.Itoa(*args.Id)

	queryParams := url.Values{}
	queryParams.Add("$expand", "1")

	locationId, _ := uuid.Parse("86c8381e-5aee-4cde-8ae4-25c0c7f5eaea")
	resp, err := fullClient.Client.Send(ctx, http.MethodGet, locationId, "5.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue pipelineschecks.CheckConfiguration
	err = fullClient.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}
