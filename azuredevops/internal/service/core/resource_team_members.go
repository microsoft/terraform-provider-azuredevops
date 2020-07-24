package core

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func ResourceTeamMembers() *schema.Resource {
	return &schema.Resource{
		Create: resourceTeamMembersCreate,
		Read:   resourceTeamMembersRead,
		Update: resourceTeamMembersUpdate,
		Delete: resourceTeamMembersDelete,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"team_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"members": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
				Computed:   true,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Set:        schema.HashString,
			},
		},
	}
}

func resourceTeamMembersCreate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("Not implemented")
}

func resourceTeamMembersRead(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("Not implemented")
}

func resourceTeamMembersUpdate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("Not implemented")
}

func resourceTeamMembersDelete(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("Not implemented")
}
