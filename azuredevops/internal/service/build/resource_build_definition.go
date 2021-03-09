package build

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/model"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/validate"
)

const (
	bdVariable              = "variable"
	bdVariableName          = "name"
	bdVariableValue         = "value"
	bdSecretVariableValue   = "secret_value"
	bdVariableIsSecret      = "is_secret"
	bdVariableAllowOverride = "allow_override"
)

// ResourceBuildDefinition schema and implementation for build definition resource
func ResourceBuildDefinition() *schema.Resource {
	filterSchema := map[string]*schema.Schema{
		"include": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
		},
		"exclude": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}

	branchFilter := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: filterSchema,
		},
	}

	pathFilter := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: filterSchema,
		},
	}

	return &schema.Resource{
		Create:   resourceBuildDefinitionCreate,
		Read:     resourceBuildDefinitionRead,
		Update:   resourceBuildDefinitionUpdate,
		Delete:   resourceBuildDefinitionDelete,
		Importer: tfhelper.ImportProjectQualifiedResource(),
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"revision": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"path": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      `\`,
				ValidateFunc: validate.Path,
			},
			"variable_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntAtLeast(1),
				},
			},
			bdVariable: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						bdVariableName: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
						bdVariableValue: {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						bdSecretVariableValue: {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
							Default:   "",
						},
						bdVariableIsSecret: {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						bdVariableAllowOverride: {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
			"agent_pool_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Hosted Ubuntu 1604",
			},
			"repository": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"yml_path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"repo_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"repo_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(model.RepoTypeValues.GitHub),
								string(model.RepoTypeValues.TfsGit),
								string(model.RepoTypeValues.Bitbucket),
								string(model.RepoTypeValues.GitHubEnterprise),
							}, false),
						},
						"branch_name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "master",
						},
						"service_connection_id": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"github_enterprise_url": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"report_build_status": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
			"ci_trigger": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"use_yaml": {
							Type:          schema.TypeBool,
							Optional:      true,
							Default:       false,
							ConflictsWith: []string{"ci_trigger.0.override"},
						},
						"override": {
							Type:     schema.TypeList,
							Optional: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"batch": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"branch_filter": branchFilter,
									"max_concurrent_builds_per_branch": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  1,
									},
									"path_filter": pathFilter,
									"polling_interval": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"polling_job_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"pull_request_trigger": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"use_yaml": {
							Type:          schema.TypeBool,
							Optional:      true,
							Default:       false,
							ConflictsWith: []string{"pull_request_trigger.0.override"},
						},
						"initial_branch": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Managed by Terraform",
						},
						"override": {
							Type:     schema.TypeList,
							Optional: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auto_cancel": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"branch_filter": branchFilter,
									"path_filter":   pathFilter,
								},
							},
						},
						"forks": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"share_secrets": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"comment_required": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"All", "NonTeamMembers"}, false),
						},
					},
				},
			},
			"tags": &model.TagsSchema,
		},
	}
}

func resourceBuildDefinitionCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	err := validateServiceConnectionIDExistsIfNeeded(d)
	if err != nil {
		return err
	}
	buildDefinition, projectID, err := expandBuildDefinition(d)
	if err != nil {
		return fmt.Errorf("error creating resource Build Definition: %+v", err)
	}

	createdBuildDefinition, err := createBuildDefinition(clients, buildDefinition, projectID)
	if err != nil {
		return fmt.Errorf("error creating resource Build Definition: %+v", err)
	}

	flattenBuildDefinition(d, createdBuildDefinition, projectID)
	return resourceBuildDefinitionRead(d, m)
}

func flattenBuildDefinition(d *schema.ResourceData, buildDefinition *build.BuildDefinition, projectID string) {
	d.SetId(strconv.Itoa(*buildDefinition.Id))

	d.Set("project_id", projectID)
	d.Set("name", *buildDefinition.Name)
	d.Set("path", *buildDefinition.Path)
	d.Set("repository", flattenRepository(buildDefinition))

	if buildDefinition.Queue != nil && buildDefinition.Queue.Pool != nil {
		d.Set("agent_pool_name", *buildDefinition.Queue.Pool.Name)
	}

	d.Set("path", *buildDefinition.Path)

	if buildDefinition.Tags != nil {
		d.Set("tags", *buildDefinition.Tags)
	}

	d.Set("variable_groups", flattenVariableGroups(buildDefinition))
	d.Set(bdVariable, flattenBuildVariables(d, buildDefinition))

	if buildDefinition.Triggers != nil {
		yamlCiTrigger := hasSettingsSourceType(buildDefinition.Triggers, build.DefinitionTriggerTypeValues.ContinuousIntegration, 2)
		d.Set("ci_trigger", flattenBuildDefinitionTriggers(buildDefinition.Triggers, yamlCiTrigger, build.DefinitionTriggerTypeValues.ContinuousIntegration))

		yamlPrTrigger := hasSettingsSourceType(buildDefinition.Triggers, build.DefinitionTriggerTypeValues.PullRequest, 2)
		d.Set("pull_request_trigger", flattenBuildDefinitionTriggers(buildDefinition.Triggers, yamlPrTrigger, build.DefinitionTriggerTypeValues.PullRequest))
	}

	revision := 0
	if buildDefinition.Revision != nil {
		revision = *buildDefinition.Revision
	}

	d.Set("revision", revision)
}

// Return an interface suitable for serialization into the resource state. This function ensures that
// any secrets, for which values will not be returned by the service, are not overidden with null or
// empty values
func flattenBuildVariables(d *schema.ResourceData, buildDefinition *build.BuildDefinition) interface{} {
	if buildDefinition.Variables == nil {
		return nil
	}
	variables := make([]map[string]interface{}, len(*buildDefinition.Variables))

	index := 0
	for varName, varVal := range *buildDefinition.Variables {
		var variable map[string]interface{}

		isSecret := converter.ToBool(varVal.IsSecret, false)
		variable = map[string]interface{}{
			bdVariableName:          varName,
			bdVariableValue:         converter.ToString(varVal.Value, ""),
			bdVariableIsSecret:      isSecret,
			bdVariableAllowOverride: converter.ToBool(varVal.AllowOverride, false),
		}

		//read secret variable from state if exist
		if isSecret {
			if stateVal := tfhelper.FindMapInSetWithGivenKeyValue(d, bdVariable, bdVariableName, varName); stateVal != nil {
				variable = stateVal
			}
		}
		variables[index] = variable
		index = index + 1
	}

	return variables
}

func createBuildDefinition(clients *client.AggregatedClient, buildDefinition *build.BuildDefinition, project string) (*build.BuildDefinition, error) {
	createdBuild, err := clients.BuildClient.CreateDefinition(clients.Ctx, build.CreateDefinitionArgs{
		Definition: buildDefinition,
		Project:    &project,
	})

	return createdBuild, err
}

func resourceBuildDefinitionRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectID, buildDefinitionID, err := tfhelper.ParseProjectIDAndResourceID(d)

	if err != nil {
		return err
	}

	buildDefinition, err := clients.BuildClient.GetDefinition(clients.Ctx, build.GetDefinitionArgs{
		Project:      &projectID,
		DefinitionId: &buildDefinitionID,
	})

	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	flattenBuildDefinition(d, buildDefinition, projectID)
	return nil
}

func resourceBuildDefinitionDelete(d *schema.ResourceData, m interface{}) error {
	if strings.EqualFold(d.Id(), "") {
		return nil
	}

	clients := m.(*client.AggregatedClient)
	projectID, buildDefinitionID, err := tfhelper.ParseProjectIDAndResourceID(d)
	if err != nil {
		return err
	}

	err = clients.BuildClient.DeleteDefinition(m.(*client.AggregatedClient).Ctx, build.DeleteDefinitionArgs{
		Project:      &projectID,
		DefinitionId: &buildDefinitionID,
	})

	return err
}

func resourceBuildDefinitionUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	err := validateServiceConnectionIDExistsIfNeeded(d)
	if err != nil {
		return err
	}
	buildDefinition, projectID, err := expandBuildDefinition(d)
	if err != nil {
		return err
	}

	updatedBuildDefinition, err := clients.BuildClient.UpdateDefinition(m.(*client.AggregatedClient).Ctx, build.UpdateDefinitionArgs{
		Definition:   buildDefinition,
		Project:      &projectID,
		DefinitionId: buildDefinition.Id,
	})

	if err != nil {
		return err
	}

	flattenBuildDefinition(d, updatedBuildDefinition, projectID)
	return resourceBuildDefinitionRead(d, m)
}

func flattenVariableGroups(buildDefinition *build.BuildDefinition) []int {
	if buildDefinition.VariableGroups == nil {
		return nil
	}

	variableGroups := make([]int, len(*buildDefinition.VariableGroups))

	for i, variableGroup := range *buildDefinition.VariableGroups {
		variableGroups[i] = *variableGroup.Id
	}

	return variableGroups
}

func flattenRepository(buildDefinition *build.BuildDefinition) interface{} {
	yamlFilePath := ""
	githubEnterpriseURL := ""

	// The process member can be of many types -- the only typing information
	// available from the compiler is `interface{}` so we can probe for known
	// implementations
	if processMap, ok := buildDefinition.Process.(map[string]interface{}); ok {
		yamlFilePath = processMap["yamlFilename"].(string)
	}
	if yamlProcess, ok := buildDefinition.Process.(*build.YamlProcess); ok {
		yamlFilePath = *yamlProcess.YamlFilename
	}

	// Set github_enterprise_url value from buildDefinition.Repository URL
	if strings.EqualFold(*buildDefinition.Repository.Type, string(model.RepoTypeValues.GitHubEnterprise)) {
		url, err := url.Parse(*buildDefinition.Repository.Url)
		if err != nil {
			return fmt.Errorf("Unable to parse repository URL: %+v ", err)
		}
		githubEnterpriseURL = fmt.Sprintf("%s://%s", url.Scheme, url.Host)
	}

	reportBuildStatus, err := strconv.ParseBool((*buildDefinition.Repository.Properties)["reportBuildStatus"])
	if err != nil {
		return fmt.Errorf("Unable to parse `reportBuildStatus` property: %+v ", err)
	}
	return []map[string]interface{}{{
		"yml_path":              yamlFilePath,
		"repo_id":               *buildDefinition.Repository.Id,
		"repo_type":             *buildDefinition.Repository.Type,
		"branch_name":           *buildDefinition.Repository.DefaultBranch,
		"service_connection_id": (*buildDefinition.Repository.Properties)["connectedServiceId"],
		"github_enterprise_url": githubEnterpriseURL,
		"report_build_status":   reportBuildStatus,
	}}
}

func flattenBuildDefinitionBranchOrPathFilter(m []interface{}) []interface{} {
	var include []string
	var exclude []string

	for _, v := range m {
		if v2, ok := v.(string); ok {
			if strings.HasPrefix(v2, "-") {
				exclude = append(exclude, strings.TrimPrefix(v2, "-"))
			} else if strings.HasPrefix(v2, "+") {
				include = append(include, strings.TrimPrefix(v2, "+"))
			}
		}
	}

	return []interface{}{
		map[string]interface{}{
			"include": include,
			"exclude": exclude,
		},
	}
}

func flattenBuildDefinitionContinuousIntegrationTrigger(m interface{}, isYaml bool) interface{} {
	if ms, ok := m.(map[string]interface{}); ok {
		f := map[string]interface{}{
			"use_yaml": isYaml,
		}
		if !isYaml {
			f["override"] = []map[string]interface{}{{
				"batch":                            ms["batchChanges"],
				"branch_filter":                    flattenBuildDefinitionBranchOrPathFilter(ms["branchFilters"].([]interface{})),
				"max_concurrent_builds_per_branch": ms["maxConcurrentBuildsPerBranch"],
				"polling_interval":                 ms["pollingInterval"],
				"polling_job_id":                   ms["pollingJobId"],
				"path_filter":                      flattenBuildDefinitionBranchOrPathFilter(ms["pathFilters"].([]interface{})),
			}}
		}
		return f
	}
	return nil
}

func flattenBuildDefinitionPullRequestTrigger(m interface{}, isYaml bool) interface{} {
	if ms, ok := m.(map[string]interface{}); ok {
		forks := ms["forks"].(map[string]interface{})
		isCommentRequired := ms["isCommentRequiredForPullRequest"].(bool)
		isCommentRequiredNonTeam := ms["requireCommentsForNonTeamMembersOnly"].(bool)

		var commentRequired string
		if isCommentRequired {
			commentRequired = "All"
		}
		if isCommentRequired && isCommentRequiredNonTeam {
			commentRequired = "NonTeamMembers"
		}

		branchFilters := ms["branchFilters"].([]interface{})
		var initialBranch string
		if len(branchFilters) > 0 {
			initialBranch = strings.TrimPrefix(branchFilters[0].(string), "+")
		}

		f := map[string]interface{}{
			"use_yaml":         isYaml,
			"initial_branch":   initialBranch,
			"comment_required": commentRequired,
			"forks": []map[string]interface{}{{
				"enabled":       forks["enabled"],
				"share_secrets": forks["allowSecrets"],
			}},
		}
		if !isYaml {
			f["override"] = []map[string]interface{}{{
				"auto_cancel":   ms["autoCancel"],
				"branch_filter": flattenBuildDefinitionBranchOrPathFilter(branchFilters),
				"path_filter":   flattenBuildDefinitionBranchOrPathFilter(ms["pathFilters"].([]interface{})),
			}}
		}
		return f
	}
	return nil
}

func flattenBuildDefinitionTrigger(m interface{}, isYaml bool, t build.DefinitionTriggerType) interface{} {
	if ms, ok := m.(map[string]interface{}); ok {
		if ms["triggerType"].(string) != string(t) {
			return nil
		}
		switch t {
		case build.DefinitionTriggerTypeValues.ContinuousIntegration:
			return flattenBuildDefinitionContinuousIntegrationTrigger(ms, isYaml)
		case build.DefinitionTriggerTypeValues.PullRequest:
			return flattenBuildDefinitionPullRequestTrigger(ms, isYaml)
		}
	}
	return nil
}

func hasSettingsSourceType(m *[]interface{}, t build.DefinitionTriggerType, sst int) bool {
	hasSetting := false
	for _, d := range *m {
		if ms, ok := d.(map[string]interface{}); ok {
			if strings.EqualFold(ms["triggerType"].(string), string(t)) {
				if val, ok := ms["settingsSourceType"]; ok {
					hasSetting = int(val.(float64)) == sst
				}
			}
		}
	}
	return hasSetting
}

func flattenBuildDefinitionTriggers(m *[]interface{}, isYaml bool, t build.DefinitionTriggerType) []interface{} {
	ds := make([]interface{}, 0, len(*m))
	for _, d := range *m {
		f := flattenBuildDefinitionTrigger(d, isYaml, t)
		if f != nil {
			ds = append(ds, f)
		}
	}
	return ds
}

func expandBuildDefinitionBranchOrPathFilter(d map[string]interface{}) []interface{} {
	include := tfhelper.ExpandStringSet(d["include"].(*schema.Set))
	exclude := tfhelper.ExpandStringSet(d["exclude"].(*schema.Set))
	m := make([]interface{}, len(include)+len(exclude))
	i := 0
	for _, v := range include {
		m[i] = "+" + v
		i++
	}
	for _, v := range exclude {
		m[i] = "-" + v
		i++
	}
	return m
}
func expandBuildDefinitionBranchOrPathFilterList(d []interface{}) [][]interface{} {
	vs := make([][]interface{}, 0, len(d))
	for _, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			vs = append(vs, expandBuildDefinitionBranchOrPathFilter(val))
		}
	}
	return vs
}
func expandBuildDefinitionBranchOrPathFilterSet(configured *schema.Set) []interface{} {
	d2 := expandBuildDefinitionBranchOrPathFilterList(configured.List())
	if len(d2) != 1 {
		return nil
	}
	return d2[0]
}

func expandBuildDefinitionFork(d map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"allowSecrets": d["share_secrets"].(bool),
		"enabled":      d["enabled"].(bool),
	}
}
func expandBuildDefinitionForkList(d []interface{}) []map[string]interface{} {
	vs := make([]map[string]interface{}, 0, len(d))
	for _, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			vs = append(vs, expandBuildDefinitionFork(val))
		}
	}
	return vs
}
func expandBuildDefinitionForkListFirstOrNil(d []interface{}) map[string]interface{} {
	d2 := expandBuildDefinitionForkList(d)
	if len(d2) != 1 {
		return nil
	}
	return d2[0]
}

func expandBuildDefinitionManualPullRequestTrigger(d map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"branchFilters": expandBuildDefinitionBranchOrPathFilterSet(d["branch_filter"].(*schema.Set)),
		"pathFilters":   expandBuildDefinitionBranchOrPathFilterSet(d["path_filter"].(*schema.Set)),
		"autoCancel":    d["auto_cancel"].(bool),
	}
}
func expandBuildDefinitionManualPullRequestTriggerList(d []interface{}) []map[string]interface{} {
	vs := make([]map[string]interface{}, 0, len(d))
	for _, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			vs = append(vs, expandBuildDefinitionManualPullRequestTrigger(val))
		}
	}
	return vs
}
func expandBuildDefinitionManualPullRequestTriggerListFirstOrNil(d []interface{}) map[string]interface{} {
	d2 := expandBuildDefinitionManualPullRequestTriggerList(d)
	if len(d2) != 1 {
		return nil
	}
	return d2[0]
}

func expandBuildDefinitionManualContinuousIntegrationTrigger(d map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"batchChanges":                 d["batch"].(bool),
		"branchFilters":                expandBuildDefinitionBranchOrPathFilterSet(d["branch_filter"].(*schema.Set)),
		"maxConcurrentBuildsPerBranch": d["max_concurrent_builds_per_branch"].(int),
		"pathFilters":                  expandBuildDefinitionBranchOrPathFilterSet(d["path_filter"].(*schema.Set)),
		"triggerType":                  string(build.DefinitionTriggerTypeValues.ContinuousIntegration),
		"pollingInterval":              d["polling_interval"].(int),
	}
}
func expandBuildDefinitionManualContinuousIntegrationTriggerList(d []interface{}) []map[string]interface{} {
	vs := make([]map[string]interface{}, 0, len(d))
	for _, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			vs = append(vs, expandBuildDefinitionManualContinuousIntegrationTrigger(val))
		}
	}
	return vs
}
func expandBuildDefinitionManualContinuousIntegrationTriggerListFirstOrNil(d []interface{}) map[string]interface{} {
	d2 := expandBuildDefinitionManualContinuousIntegrationTriggerList(d)
	if len(d2) != 1 {
		return nil
	}
	return d2[0]
}

func expandBuildDefinitionTrigger(d map[string]interface{}, t build.DefinitionTriggerType) interface{} {
	switch t {
	case build.DefinitionTriggerTypeValues.ContinuousIntegration:
		isYaml := d["use_yaml"].(bool)
		if isYaml {
			return map[string]interface{}{
				"batchChanges":                 false,
				"branchFilters":                []interface{}{},
				"maxConcurrentBuildsPerBranch": 1,
				"pathFilters":                  []interface{}{},
				"triggerType":                  string(t),
				"settingsSourceType":           float64(2),
			}
		}
		return expandBuildDefinitionManualContinuousIntegrationTriggerListFirstOrNil(d["override"].([]interface{}))
	case build.DefinitionTriggerTypeValues.PullRequest:
		isYaml := d["use_yaml"].(bool)
		commentRequired := d["comment_required"].(string)
		vs := map[string]interface{}{
			"forks":                                expandBuildDefinitionForkListFirstOrNil(d["forks"].([]interface{})),
			"isCommentRequiredForPullRequest":      len(commentRequired) > 0,
			"requireCommentsForNonTeamMembersOnly": commentRequired == "NonTeamMembers",
			"triggerType":                          string(t),
		}
		if isYaml {
			vs["branchFilters"] = []interface{}{
				"+" + d["initial_branch"].(string),
			}
			vs["pathFilters"] = []interface{}{}
			vs["settingsSourceType"] = float64(2)
		} else {
			override := expandBuildDefinitionManualPullRequestTriggerListFirstOrNil(d["override"].([]interface{}))
			vs["branchFilters"] = override["branchFilters"]
			vs["pathFilters"] = override["pathFilters"]
			vs["autoCancel"] = override["autoCancel"]
		}
		return vs
	}
	return nil
}
func expandBuildDefinitionTriggerList(d []interface{}, t build.DefinitionTriggerType) []interface{} {
	vs := make([]interface{}, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			vs = append(vs, expandBuildDefinitionTrigger(val, t))
		}
	}
	return vs
}

func expandVariableGroups(d *schema.ResourceData) *[]build.VariableGroup {
	variableGroupsInterface := d.Get("variable_groups").(*schema.Set).List()
	variableGroups := make([]build.VariableGroup, len(variableGroupsInterface))

	for i, variableGroup := range variableGroupsInterface {
		variableGroups[i] = *buildVariableGroup(variableGroup.(int))
	}

	return &variableGroups
}

func expandVariables(d *schema.ResourceData) (*map[string]build.BuildDefinitionVariable, error) {
	variables := d.Get(bdVariable)
	if variables == nil {
		return nil, nil
	}

	variablesList := variables.(*schema.Set).List()
	if len(variablesList) == 0 {
		return nil, nil
	}

	expandedVars := map[string]build.BuildDefinitionVariable{}
	for _, variable := range variablesList {
		varAsMap := variable.(map[string]interface{})
		varName := varAsMap[bdVariableName].(string)

		if _, ok := expandedVars[varName]; ok {
			return nil, fmt.Errorf("Unexpectedly found duplicate variable with name %s", varName)
		}

		isSecret := converter.Bool(varAsMap[bdVariableIsSecret].(bool))
		var val *string

		if *isSecret {
			val = converter.String(varAsMap[bdSecretVariableValue].(string))
		} else {
			val = converter.String(varAsMap[bdVariableValue].(string))
		}
		expandedVars[varName] = build.BuildDefinitionVariable{
			AllowOverride: converter.Bool(varAsMap[bdVariableAllowOverride].(bool)),
			IsSecret:      isSecret,
			Value:         val,
		}
	}

	return &expandedVars, nil
}

func expandBuildDefinition(d *schema.ResourceData) (*build.BuildDefinition, string, error) {
	projectID := d.Get("project_id").(string)
	repositories := d.Get("repository").([]interface{})

	// Note: If configured, this will be of length 1 based on the schema definition above.
	if len(repositories) != 1 {
		return nil, "", fmt.Errorf("Unexpectedly did not find repository metadata in the resource data")
	}

	repository := repositories[0].(map[string]interface{})

	repoID := repository["repo_id"].(string)
	repoType := model.RepoType(repository["repo_type"].(string))
	repoURL := ""
	repoAPIURL := ""

	if strings.EqualFold(string(repoType), string(model.RepoTypeValues.GitHub)) {
		repoURL = fmt.Sprintf("https://github.com/%s.git", repoID)
		repoAPIURL = fmt.Sprintf("https://api.github.com/repos/%s", repoID)
	}
	if strings.EqualFold(string(repoType), string(model.RepoTypeValues.Bitbucket)) {
		repoURL = fmt.Sprintf("https://bitbucket.org/%s.git", repoID)
		repoAPIURL = fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s", repoID)
	}
	if strings.EqualFold(string(repoType), string(model.RepoTypeValues.GitHubEnterprise)) {
		githubEnterpriseURL := repository["github_enterprise_url"].(string)
		repoURL = fmt.Sprintf("%s/%s.git", githubEnterpriseURL, repoID)
		repoAPIURL = fmt.Sprintf("%s/api/v3/repos/%s", githubEnterpriseURL, repoID)
	}

	ciTriggers := expandBuildDefinitionTriggerList(
		d.Get("ci_trigger").([]interface{}),
		build.DefinitionTriggerTypeValues.ContinuousIntegration,
	)
	pullRequestTriggers := expandBuildDefinitionTriggerList(
		d.Get("pull_request_trigger").([]interface{}),
		build.DefinitionTriggerTypeValues.PullRequest,
	)

	buildTriggers := append(ciTriggers, pullRequestTriggers...)

	// Look for the ID. This may not exist if we are within the context of a "create" operation,
	// so it is OK if it is missing.
	buildDefinitionID, err := strconv.Atoi(d.Id())
	var buildDefinitionReference *int
	if err == nil {
		buildDefinitionReference = &buildDefinitionID
	} else {
		buildDefinitionReference = nil
	}

	variables, err := expandVariables(d)
	if err != nil {
		return nil, "", fmt.Errorf("Error expanding varibles: %+v", err)
	}

	tags := tfhelper.ExpandStringList(d.Get("tags").([]interface{}))

	buildDefinition := build.BuildDefinition{
		Id:       buildDefinitionReference,
		Name:     converter.String(d.Get("name").(string)),
		Path:     converter.String(d.Get("path").(string)),
		Revision: converter.Int(d.Get("revision").(int)),
		Repository: &build.BuildRepository{
			Url:           &repoURL,
			Id:            &repoID,
			Name:          &repoID,
			DefaultBranch: converter.String(repository["branch_name"].(string)),
			Type:          converter.String(string(repoType)),
			Properties: &map[string]string{
				"connectedServiceId": repository["service_connection_id"].(string),
				"apiUrl":             repoAPIURL,
				"reportBuildStatus":  strconv.FormatBool(repository["report_build_status"].(bool)),
			},
		},
		Process: &build.YamlProcess{
			YamlFilename: converter.String(repository["yml_path"].(string)),
		},
		QueueStatus:    &build.DefinitionQueueStatusValues.Enabled,
		Type:           &build.DefinitionTypeValues.Build,
		Quality:        &build.DefinitionQualityValues.Definition,
		VariableGroups: expandVariableGroups(d),
		Variables:      variables,
		Triggers:       &buildTriggers,
		Tags:           &tags,
	}

	if agentPoolName, ok := d.GetOk("agent_pool_name"); ok {
		buildDefinition.Queue = &build.AgentPoolQueue{
			Name: converter.StringFromInterface(agentPoolName),
			Pool: &build.TaskAgentPoolReference{
				Name: converter.StringFromInterface(agentPoolName),
			},
		}
	}

	return &buildDefinition, projectID, nil
}

/**
 * certain types of build definitions require a service connection to run. This function
 * returns an error if a service connection was needed but not provided
 */
func validateServiceConnectionIDExistsIfNeeded(d *schema.ResourceData) error {
	repositories := d.Get("repository").([]interface{})
	repository := repositories[0].(map[string]interface{})

	repoType := repository["repo_type"].(string)
	serviceConnectionID := repository["service_connection_id"].(string)

	if strings.EqualFold(repoType, string(model.RepoTypeValues.Bitbucket)) && serviceConnectionID == "" {
		return errors.New("bitbucket repositories need a referenced service connection ID")
	}
	if strings.EqualFold(repoType, string(model.RepoTypeValues.GitHubEnterprise)) && serviceConnectionID == "" {
		return errors.New("GitHub Enterprise repositories need a referenced service connection ID")
	}
	return nil
}

func buildVariableGroup(id int) *build.VariableGroup {
	return &build.VariableGroup{
		Id: &id,
	}
}
