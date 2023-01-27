package serviceendpoint

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/pipelineschecks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

var taskCheckType = pipelineschecks.CheckType{
	Id: converter.UUID("fe1de3ee-a436-41b4-bb20-f6eb4cb879a7"),
}

var evaluateBranchProtectionDefVersion = "0.0.1"
var evaluateBranchProtectionDefId = "86b05a0c-73e6-4f7d-b3cf-e38f3b39a75b"

var evaluateBranchProtectionDef = map[string]string{
	"id":      evaluateBranchProtectionDefId,
	"name":    "evaluatebranchProtection",
	"version": evaluateBranchProtectionDefVersion,
}

// ResourceBranchControlCheck schema and implementation for build definition resource
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
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "endpoint",
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

	branchControlCheck, err := clients.V5PipelinesChecksClientExtras.GetCheckConfiguration(clients.Ctx, pipelineschecks.GetCheckConfigurationArgs{
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

	d.Set("project_id", projectID)
	d.Set("endpoint_id", branchControlCheck.Resource.Id)
	d.Set("resource_type", branchControlCheck.Resource.Type)

	if branchControlCheck.Settings == nil {
		return fmt.Errorf("Settings nil")
	}

	var definitionRef map[string]interface{}

	if definitionRefMap, found := branchControlCheck.Settings.(map[string]interface{})["definitionRef"]; found {
		definitionRef = definitionRefMap.(map[string]interface{})
	} else {
		return fmt.Errorf("definitionRef not found")
	}

	if id, found := definitionRef["id"]; found {
		if !strings.EqualFold(id.(string), evaluateBranchProtectionDefId) {
			return fmt.Errorf("invalid definitionRef id, not a branch control")
		}
	} else {
		return fmt.Errorf("definitionRef id not found")
	}

	if version, found := definitionRef["version"]; found {
		if version != evaluateBranchProtectionDefVersion {
			return fmt.Errorf("unsupported definitionRef version")
		}
	} else {
		return fmt.Errorf("unsupported definitionRef version")
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
		value, err := strconv.ParseBool(verifyBranchProtection.(string))
		if err != nil {
			return err
		}
		d.Set("verify_branch_protection", value)
	} else {
		return fmt.Errorf("ensureProtectionOfBranch input not found")
	}

	if ignoreUnknownProtectionStatus, found := inputs["allowUnknownStatusBranch"]; found {
		value, err := strconv.ParseBool(ignoreUnknownProtectionStatus.(string))
		if err != nil {
			return err
		}
		d.Set("ignore_unknown_protection_status", value)
	} else {
		return fmt.Errorf("allowUnknownStatusBranch input not found")
	}

	return nil
}

func expandBranchControlCheck(d *schema.ResourceData) (*pipelineschecks.CheckConfiguration, string, error) {
	projectID := d.Get("project_id").(string)
	endpointID := d.Get("endpoint_id").(string)
	displayName := d.Get("display_name").(string)
	allowedBranches := d.Get("allowed_branches").(string)
	verifyBranchProtection := d.Get("verify_branch_protection").(bool)
	ignoreUnknownProtectionStatus := d.Get("ignore_unknown_protection_status").(bool)

	endpointType := "endpoint"
	endpointResource := pipelineschecks.Resource{
		Id:   &endpointID,
		Type: &endpointType,
	}

	branchControlInputs := map[string]string{
		"allowedBranches":          allowedBranches,
		"ensureProtectionOfBranch": strconv.FormatBool(verifyBranchProtection),
		"allowUnknownStatusBranch": strconv.FormatBool(ignoreUnknownProtectionStatus),
	}

	branchControlCheckSettings := map[string]interface{}{
		"definitionRef": evaluateBranchProtectionDef,
		"displayName":   displayName,
		"inputs":        branchControlInputs,
	}

	branchControlCheck := pipelineschecks.CheckConfiguration{
		Type:     &taskCheckType,
		Settings: branchControlCheckSettings,
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
