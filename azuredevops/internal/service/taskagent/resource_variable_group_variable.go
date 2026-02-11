package taskagent

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

const VariableGroupVariable = "azuredevops_variable_group_variable"

var forEachLock = new(sync.Mutex)

func ResourceVariableGroupVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceVariableGroupVariableCreateUpdate,
		Read:   resourceVariableGroupVariableRead,
		Update: resourceVariableGroupVariableCreateUpdate,
		Delete: resourceVariableGroupVariableDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"variable_group_id": {
				// TODO: Ideally this shall be an int, but the existing group id is a string.
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"value": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"value", "secret_value"},
			},
			"secret_value": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ExactlyOneOf: []string{"value", "secret_value"},
			},
		},
	}
}

func resourceVariableGroupVariableCreateUpdate(d *schema.ResourceData, m interface{}) error {

	forEachLock.Lock()
	defer forEachLock.Unlock()

	clients := m.(*client.AggregatedClient)

	projectId := d.Get("project_id").(string)
	variableGroupId, err := strconv.Atoi(d.Get("variable_group_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `variable_group_id` as an integer: %v", err)
	}

	resp, err := clients.TaskAgentClient.GetVariableGroup(
		clients.Ctx,
		taskagent.GetVariableGroupArgs{
			GroupId: &variableGroupId,
			Project: &projectId,
		},
	)
	if err != nil {
		return fmt.Errorf("Looking up variable group given ID (%v) and project ID (%v): %+v", variableGroupId, projectId, err)
	}
	if resp.Variables == nil {
		return fmt.Errorf("unexpected null existing variables")
	}
	vars := *resp.Variables
	name := d.Get("name").(string)
	id := fmt.Sprintf("%s/%d/%s", projectId, variableGroupId, name)

	// Existence check
	if d.IsNewResource() {
		if _, ok := vars[name]; ok {
			return tfhelper.ImportAsExistsError(VariableGroupVariable, id)
		}
	}

	// Upsert the variable
	params := taskagent.VariableGroupParameters{
		Description:                    resp.Description,
		Name:                           resp.Name,
		ProviderData:                   resp.ProviderData,
		Type:                           resp.Type,
		VariableGroupProjectReferences: resp.VariableGroupProjectReferences,
	}

	var (
		value    string
		isSecret bool
	)
	cfgMap := d.GetRawConfig().AsValueMap()
	if val := cfgMap["value"]; !val.IsNull() {
		value = val.AsString()
	} else if val := cfgMap["secret_value"]; !val.IsNull() {
		value = val.AsString()
		isSecret = true
	}
	vars[name] = map[string]any{
		"value":    value,
		"isSecret": isSecret,
	}
	params.Variables = &vars

	if _, err := updateVariableGroup(clients, &params, &variableGroupId); err != nil {
		return err
	}

	d.SetId(id)

	return resourceVariableGroupVariableRead(d, m)
}

func resourceVariableGroupVariableRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectId, variableGroupId, varName, err := ResourceVariableGroupVariableParseId(d.Id())
	if err != nil {
		return err
	}

	resp, err := clients.TaskAgentClient.GetVariableGroup(
		clients.Ctx,
		taskagent.GetVariableGroupArgs{
			GroupId: &variableGroupId,
			Project: &projectId,
		},
	)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Looking up variable group given ID (%v) and project ID (%v): %+v", variableGroupId, projectId, err)
	}

	if resp.Variables == nil {
		d.SetId("")
		return nil
	}

	vars := *resp.Variables

	varVal, ok := vars[varName]
	if !ok {
		d.SetId("")
		return nil
	}

	var (
		value    string
		isSecret bool
	)
	if varMap, ok := varVal.(map[string]any); ok {
		if v, ok := varMap["value"].(string); ok {
			value = v
		}
		if v, ok := varMap["isSecret"].(bool); ok {
			isSecret = v
		}
	}
	if isSecret {
		// Azure doesn't store the secret value, read it from state
		value = d.Get("secret_value").(string)
	}

	d.Set("project_id", projectId)
	d.Set("variable_group_id", strconv.Itoa(variableGroupId))
	d.Set("name", varName)
	if isSecret {
		d.Set("secret_value", value)
	} else {
		d.Set("value", value)
	}
	return nil
}

func resourceVariableGroupVariableDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectId := d.Get("project_id").(string)
	variableGroupId, err := strconv.Atoi(d.Get("variable_group_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `variable_group_id` as an integer: %v", err)
	}

	resp, err := clients.TaskAgentClient.GetVariableGroup(
		clients.Ctx,
		taskagent.GetVariableGroupArgs{
			GroupId: &variableGroupId,
			Project: &projectId,
		},
	)
	if err != nil {
		return fmt.Errorf("Looking up variable group given ID (%v) and project ID (%v): %+v", variableGroupId, projectId, err)
	}
	if resp.Variables == nil {
		return fmt.Errorf("unexpected null existing variables")
	}
	vars := *resp.Variables

	name := d.Get("name").(string)
	if _, ok := vars[name]; !ok {
		// If the var doesn't exist, just return
		return nil
	}

	// Delete the variable
	delete(vars, name)
	params := taskagent.VariableGroupParameters{
		Description:                    resp.Description,
		Name:                           resp.Name,
		ProviderData:                   resp.ProviderData,
		Type:                           resp.Type,
		VariableGroupProjectReferences: resp.VariableGroupProjectReferences,
		Variables:                      &vars,
	}

	if _, err := updateVariableGroup(clients, &params, &variableGroupId); err != nil {
		return err
	}
	return nil
}

func ResourceVariableGroupVariableParseId(id string) (string, int, string, error) {
	segs := strings.SplitN(id, "/", 3)
	if len(segs) != 3 {
		return "", 0, "", fmt.Errorf("invalid resource id, expect length=3, got=%d", len(segs))
	}
	projectId, variableGroupIdStr, varName := segs[0], segs[1], segs[2]
	variableGroupId, err := strconv.Atoi(variableGroupIdStr)
	if err != nil {
		return "", 0, "", fmt.Errorf("converting the variable group id as integer: %v", err)
	}
	return projectId, variableGroupId, varName, nil
}
