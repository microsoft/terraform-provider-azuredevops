package model

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// TagsSchema list of tags
var TagsSchema = schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
}
