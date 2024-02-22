package approvalsandchecks

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/pipelineschecksextras"
)

// NOTE: In theory the API should accept "agentpool" as well, but the API client requires a project ID
// so it doesn't seem to work and the website UI doesn't have it available
var targetResourceTypes = []string{"endpoint", "environment", "queue", "repository", "securefile", "variablegroup"}

type flatFunc func(d *schema.ResourceData, check *pipelineschecksextras.CheckConfiguration, projectID string) error
type expandFunc func(d *schema.ResourceData) (*pipelineschecksextras.CheckConfiguration, string, error)

type approvalAndCheckTypes struct {
	ExtendsCheck     *pipelineschecksextras.CheckType
	Approval         *pipelineschecksextras.CheckType
	BranchProtection *pipelineschecksextras.CheckType
	BusinessHours    *pipelineschecksextras.CheckType
	TaskCheck        *pipelineschecksextras.CheckType
	ExclusiveLock    *pipelineschecksextras.CheckType
}

var approvalAndCheckType = approvalAndCheckTypes{
	ExtendsCheck: &pipelineschecksextras.CheckType{
		Id: converter.UUID("4020e66e-b0f3-47e1-bc88-48f3cc59b5f3"),
	},
	Approval: &pipelineschecksextras.CheckType{
		Id:   converter.UUID("8c6f20a7-a545-4486-9777-f762fafe0d4d"),
		Name: converter.ToPtr("Approval"),
	},
	TaskCheck: &pipelineschecksextras.CheckType{
		Id: converter.UUID("fe1de3ee-a436-41b4-bb20-f6eb4cb879a7"),
	},
	BranchProtection: &pipelineschecksextras.CheckType{
		Id: converter.UUID("fe1de3ee-a436-41b4-bb20-f6eb4cb879a7"),
	},
	BusinessHours: &pipelineschecksextras.CheckType{
		Id: converter.UUID("fe1de3ee-a436-41b4-bb20-f6eb4cb879a7"),
	},
	ExclusiveLock: &pipelineschecksextras.CheckType{
		Id: converter.UUID("2ef31ad6-baa0-403a-8b45-2cbc9b4e5563"),
	},
}

// genBaseCheckResource creates a Resource with the common parts
// that all checks require.
func genBaseCheckResource(f flatFunc, e expandFunc) *schema.Resource {
	return &schema.Resource{
		Create: genCheckCreateFunc(f, e),
		Read:   genCheckReadFunc(f),
		Update: genCheckUpdateFunc(f, e),
		Delete: genCheckDeleteFunc(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
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
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

// doBaseExpansion performs the expansion for the 'base' attributes that are defined in the schema, above
func doBaseExpansion(d *schema.ResourceData, checkType *pipelineschecksextras.CheckType, settings map[string]interface{}, timeout *int) (*pipelineschecksextras.CheckConfiguration, string, error) {
	projectID := d.Get("project_id").(string)

	taskCheck := pipelineschecksextras.CheckConfiguration{
		Type:     checkType,
		Settings: settings,
		Resource: &pipelineschecksextras.Resource{
			Id:   converter.String(d.Get("target_resource_id").(string)),
			Type: converter.String(d.Get("target_resource_type").(string)),
		},
		Version: converter.Int(d.Get("version").(int)),
	}

	if timeout != nil {
		taskCheck.Timeout = timeout
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
func doBaseFlattening(d *schema.ResourceData, check *pipelineschecksextras.CheckConfiguration, projectID string) error {
	d.SetId(fmt.Sprintf("%d", *check.Id))

	d.Set("project_id", projectID)

	if check.Resource == nil {
		return fmt.Errorf("Resource nil")
	}

	d.Set("target_resource_id", check.Resource.Id)
	d.Set("target_resource_type", check.Resource.Type)
	d.Set("version", check.Version)

	return nil
}

func genCheckCreateFunc(flatFunc flatFunc, expandFunc expandFunc) func(d *schema.ResourceData, m interface{}) error {
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*client.AggregatedClient)
		configuration, projectID, err := expandFunc(d)
		if err != nil {
			return fmt.Errorf(" failed in expandFunc. Error: %+v", err)
		}

		createdCheck, err := clients.PipelinesChecksClientExtras.AddCheckConfiguration(clients.Ctx, pipelineschecksextras.AddCheckConfigurationArgs{
			Project:       &projectID,
			Configuration: configuration,
		})
		if err != nil {
			return fmt.Errorf(" failed creating check, project ID: %s. Error: %+v", projectID, err)
		}

		err = flatFunc(d, createdCheck, projectID)
		if err != nil {
			return err
		}
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

		taskCheck, err := clients.PipelinesChecksClientExtras.GetCheckConfiguration(clients.Ctx, pipelineschecksextras.GetCheckConfigurationArgs{
			Project: &projectID,
			Id:      &taskCheckId,
			Expand:  converter.ToPtr(pipelineschecksextras.CheckConfigurationExpandParameterValues.Settings),
		})

		if err != nil {
			if utils.ResponseWasNotFound(err) || strings.Contains(err.Error(), "does not exist.") {
				d.SetId("")
				return nil
			}
			return err
		}

		return flatFunc(d, taskCheck, projectID)
	}
}

func genCheckUpdateFunc(flatFunc flatFunc, expandFunc expandFunc) schema.UpdateFunc { //nolint:staticcheck
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*client.AggregatedClient)
		taskCheck, projectID, err := expandFunc(d)
		if err != nil {
			return err
		}

		updatedBusinessHours, err := clients.PipelinesChecksClientExtras.UpdateCheckConfiguration(clients.Ctx,
			pipelineschecksextras.UpdateCheckConfigurationArgs{
				Project:       &projectID,
				Configuration: taskCheck,
				Id:            taskCheck.Id,
			})

		if err != nil {
			return err
		}

		err = flatFunc(d, updatedBusinessHours, projectID)
		if err != nil {
			return err
		}
		return genCheckReadFunc(flatFunc)(d, m)
	}
}

func genCheckDeleteFunc() schema.DeleteFunc { //nolint:staticcheck
	return func(d *schema.ResourceData, m interface{}) error {
		if strings.EqualFold(d.Id(), "") {
			return nil
		}

		clients := m.(*client.AggregatedClient)
		projectID, BusinessHoursID, err := tfhelper.ParseProjectIDAndResourceID(d)
		if err != nil {
			return err
		}

		return clients.PipelinesChecksClientExtras.DeleteCheckConfiguration(m.(*client.AggregatedClient).Ctx,
			pipelineschecksextras.DeleteCheckConfigurationArgs{
				Project: &projectID,
				Id:      &BusinessHoursID,
			})
	}
}
