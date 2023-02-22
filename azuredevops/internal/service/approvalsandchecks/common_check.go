package approvalsandchecks

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/pipelineschecks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

var taskCheckType = pipelineschecks.CheckType{
	Id: converter.UUID("fe1de3ee-a436-41b4-bb20-f6eb4cb879a7"),
}

// NOTE: In theory the API should accept "agentpool" as well, but the API client requires a project ID
// so it doesn't seem to work and the website UI doesn't have it available
var targetResourceTypes = []string{"endpoint", "environment", "queue", "repository", "securefile", "variablegroup"}

type flatFunc func(d *schema.ResourceData, check *pipelineschecks.CheckConfiguration, projectID string) error
type expandFunc func(d *schema.ResourceData) (*pipelineschecks.CheckConfiguration, string, error)

// genBaseCheckResource creates a Resource with the common parts
// that all checks require.
func genBaseCheckResource(f flatFunc, e expandFunc) *schema.Resource {
	return &schema.Resource{
		Create: genCheckCreateFunc(f, e),
		Read:   genCheckReadFunc(f),
		Update: genCheckUpdateFunc(f, e),
		Delete: genCheckDeleteFunc(e),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Importer: nil,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"target_resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"target_resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(targetResourceTypes, false),
			},
			"display_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Managed by Terraform",
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

// doBaseExpansion performs the expansion for the 'base' attributes that are defined in the schema, above
func doBaseExpansion(d *schema.ResourceData, inputs map[string]interface{}, definitionRef interface{}) (*pipelineschecks.CheckConfiguration, string, error) {
	projectID := d.Get("project_id").(string)

	taskCheck := pipelineschecks.CheckConfiguration{
		Type: &taskCheckType,
		Settings: map[string]interface{}{
			"definitionRef": definitionRef,
			"displayName":   d.Get("display_name").(string),
			"inputs":        inputs,
		},
		Resource: &pipelineschecks.Resource{
			Id:   converter.String(d.Get("target_resource_id").(string)),
			Type: converter.String(d.Get("target_resource_type").(string)),
		},
	}

	if d.Id() != "" {
		taskCheckId, err := strconv.Atoi(d.Id())
		if err != nil {
			return nil, "", fmt.Errorf("Error parsing task check ID: (%+v)", err)
		}
		taskCheck.Id = &taskCheckId
	}

	return &taskCheck, projectID, nil
}

// doBaseFlattening performs the flattening for the 'base' attributes that are defined in the schema, above
func doBaseFlattening(d *schema.ResourceData, check *pipelineschecks.CheckConfiguration, projectID string, definitionId string, definitionVersion string) error {
	d.SetId(fmt.Sprintf("%d", *check.Id))

	d.Set("project_id", projectID)
	d.Set("target_resource_id", check.Resource.Id)
	d.Set("target_resource_type", check.Resource.Type)

	if check.Settings == nil {
		return fmt.Errorf("Settings nil")
	}

	var definitionRef map[string]interface{}

	if definitionRefMap, found := check.Settings.(map[string]interface{})["definitionRef"]; found {
		definitionRef = definitionRefMap.(map[string]interface{})
	} else {
		return fmt.Errorf("definitionRef not found")
	}

	if id, found := definitionRef["id"]; found {
		if !strings.EqualFold(id.(string), definitionId) {
			return fmt.Errorf("invalid definitionRef id")
		}
	} else {
		return fmt.Errorf("definitionRef id not found")
	}

	if version, found := definitionRef["version"]; found {
		if version != definitionVersion {
			return fmt.Errorf("unsupported definitionRef version")
		}
	} else {
		return fmt.Errorf("unsupported definitionRef version")
	}

	if DisplayName, found := check.Settings.(map[string]interface{})["displayName"]; found {
		d.Set("display_name", DisplayName.(string))
	} else {
		return fmt.Errorf("displayName setting not found")
	}

	return nil
}

func genCheckCreateFunc(flatFunc flatFunc, expandFunc expandFunc) func(d *schema.ResourceData, m interface{}) error {
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*client.AggregatedClient)
		configuration, projectID, err := expandFunc(d)
		if err != nil {
			return fmt.Errorf(" failed in expandFunc. Error: %+v", err)
		}

		createdCheck, err := clients.V5PipelinesChecksClient.AddCheckConfiguration(clients.Ctx, pipelineschecks.AddCheckConfigurationArgs{
			Project:       &projectID,
			Configuration: configuration,
		})
		if err != nil {
			return fmt.Errorf(" failed creating check, project ID: %s. Error: %+v", projectID, err)
		}

		flatFunc(d, createdCheck, projectID)
		return genCheckReadFunc(flatFunc)(d, m)
	}
}

func genCheckReadFunc(flatFunc flatFunc) func(d *schema.ResourceData, m interface{}) error {
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*client.AggregatedClient)
		projectID, taskCheckId, err := tfhelper.ParseProjectIDAndResourceID(d)
		if err != nil {
			return err
		}

		taskCheck, err := clients.V5PipelinesChecksClientExtras.GetCheckConfiguration(clients.Ctx, pipelineschecks.GetCheckConfigurationArgs{
			Project: &projectID,
			Id:      &taskCheckId,
		})

		if err != nil {
			if utils.ResponseWasNotFound(err) || strings.Contains(err.Error(), "does not exist.") {
				d.SetId("")
				return nil
			}
			return err
		}

		flatFunc(d, taskCheck, projectID)
		return nil
	}
}

func genCheckUpdateFunc(flatFunc flatFunc, expandFunc expandFunc) schema.UpdateFunc { //nolint:staticcheck
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*client.AggregatedClient)
		taskCheck, projectID, err := expandFunc(d)
		if err != nil {
			return err
		}

		updatedBusinessHours, err := clients.V5PipelinesChecksClient.UpdateCheckConfiguration(clients.Ctx,
			pipelineschecks.UpdateCheckConfigurationArgs{
				Project:       &projectID,
				Configuration: taskCheck,
				Id:            taskCheck.Id,
			})

		if err != nil {
			return err
		}

		flatFunc(d, updatedBusinessHours, projectID)
		return genCheckReadFunc(flatFunc)(d, m)
	}
}

func genCheckDeleteFunc(expandFunc expandFunc) schema.DeleteFunc { //nolint:staticcheck
	return func(d *schema.ResourceData, m interface{}) error {
		if strings.EqualFold(d.Id(), "") {
			return nil
		}

		clients := m.(*client.AggregatedClient)
		projectID, BusinessHoursID, err := tfhelper.ParseProjectIDAndResourceID(d)
		if err != nil {
			return err
		}

		return clients.V5PipelinesChecksClient.DeleteCheckConfiguration(m.(*client.AggregatedClient).Ctx,
			pipelineschecks.DeleteCheckConfigurationArgs{
				Project: &projectID,
				Id:      &BusinessHoursID,
			})
	}
}
