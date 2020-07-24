package core

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func ResourceTeamAdministrators() *schema.Resource {
	return &schema.Resource{
		Create: resourceTeamAdministratorsCreate,
		Read:   resourceTeamAdministratorsRead,
		Update: resourceTeamAdministratorsUpdate,
		Delete: resourceTeamAdministratorsDelete,
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
			"administrators": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
				Required:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Set:        schema.HashString,
			},
		},
	}
}

func resourceTeamAdministratorsCreate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("Not implemented")
}

func resourceTeamAdministratorsRead(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("Not implemented")
}

func resourceTeamAdministratorsUpdate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("Not implemented")
}

func resourceTeamAdministratorsDelete(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("Not implemented")
}
