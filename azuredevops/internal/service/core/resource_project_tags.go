package core

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func ResourceProjectTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectTagsCreate,
		ReadContext:   resourceProjectTagsRead,
		UpdateContext: resourceProjectTagsUpdate,
		DeleteContext: resourceProjectTagsDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
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
			"tags": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func resourceProjectTagsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID, err := uuid.Parse(d.Get("project_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	tags := expandProjectTags(d.Get("tags").(*schema.Set).List())
	err = clients.CoreClient.SetProjectProperties(clients.Ctx, core.SetProjectPropertiesArgs{
		PatchDocument: tags,
		ProjectId:     &projectID,
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("Creating Project Tags. Project ID: %s, Error: %+v", projectID.String(), err))
	}

	d.SetId(projectID.String())
	return resourceProjectTagsRead(clients.Ctx, d, m)
}

func resourceProjectTagsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID, err := uuid.Parse(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	tags, err := clients.CoreClient.GetProjectProperties(clients.Ctx, core.GetProjectPropertiesArgs{
		ProjectId: &projectID,
		Keys:      &[]string{"Microsoft.TeamFoundation.Project.Tag.*"},
	})

	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("Gettting Project tags. Projecct ID: %s. Error: %+v", projectID, err))
	}

	if tags == nil || len(*tags) == 0 {
		d.SetId("")
		return nil
	}

	d.Set("project_id", d.Id())
	d.Set("tags", flattenProjectTags(*tags))
	return nil
}

func resourceProjectTagsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectId, err := uuid.Parse(d.Get("project_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	tagsLocalRaw := d.Get("tags").(*schema.Set).List()
	tagsLocal := make([]string, 0)
	for _, tag := range tagsLocalRaw {
		if tag.(string) != "" {
			tagsLocal = append(tagsLocal, tag.(string))
		}
	}

	// Get current tags
	resp, err := clients.CoreClient.GetProjectProperties(clients.Ctx, core.GetProjectPropertiesArgs{
		ProjectId: &projectId,
		Keys:      &[]string{"Microsoft.TeamFoundation.Project.Tag.*"},
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("Gettting Project tags. Projecct ID: %s. Error: %+v", projectId, err))
	}
	tagsExists := flattenProjectTags(*resp)

	// Filter tags that will be added or removed
	tagsAdd := sliceDifference(tagsExists, tagsLocal)
	tagsRemove := sliceDifference(tagsLocal, tagsExists)

	// All tags needs update
	allTagsOp := expandProjectTags(tagsRemove)
	for i := 0; i < len(*allTagsOp); i++ {
		(*allTagsOp)[i].Op = &webapi.OperationValues.Remove
		(*allTagsOp)[i].Value = nil
	}
	*allTagsOp = append(*allTagsOp, *expandProjectTags(tagsAdd)...)

	err = clients.CoreClient.SetProjectProperties(clients.Ctx, core.SetProjectPropertiesArgs{
		PatchDocument: allTagsOp,
		ProjectId:     &projectId,
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("Updating Project Tags. Project ID: %s, Error: %+v", projectId.String(), err))
	}
	return resourceProjectTagsRead(clients.Ctx, d, m)
}

func resourceProjectTagsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectId, err := uuid.Parse(d.Get("project_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	tagsRemoveOp := expandProjectTags(d.Get("tags").(*schema.Set).List())
	for i := 0; i < len(*tagsRemoveOp); i++ {
		(*tagsRemoveOp)[i].Op = &webapi.OperationValues.Remove
		(*tagsRemoveOp)[i].Value = nil
	}

	err = clients.CoreClient.SetProjectProperties(clients.Ctx, core.SetProjectPropertiesArgs{
		PatchDocument: tagsRemoveOp,
		ProjectId:     &projectId,
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("Deleting Project Tags. ProjectID: %s, Error: %+v", projectId, err))
	}
	return nil
}

func expandProjectTags(input []interface{}) *[]webapi.JsonPatchOperation {
	tags := make([]webapi.JsonPatchOperation, 0)
	for _, v := range input {
		tags = append(tags, webapi.JsonPatchOperation{
			From:  converter.String(""),
			Op:    converter.ToPtr(webapi.OperationValues.Add),
			Path:  converter.String(fmt.Sprintf("/Microsoft.TeamFoundation.Project.Tag.%s", v.(string))),
			Value: converter.String("true"),
		})
	}
	return &tags
}

func flattenProjectTags(input []core.ProjectProperty) []string {
	tags := make([]string, 0)
	for _, v := range input {
		tags = append(tags, (*v.Name)[37:])
	}
	return tags
}

func sliceDifference(input1, input2 []string) []interface{} {
	result := make([]interface{}, 0)
	inputMap := make(map[string]string, 0)
	for _, v := range input1 {
		inputMap[v] = ""
	}

	for _, v := range input2 {
		if _, ok := inputMap[v]; !ok {
			result = append(result, v)
		}
	}
	return result
}
