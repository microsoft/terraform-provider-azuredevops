package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceFoo() *schema.Resource {
	return &schema.Resource{
		Create: resourceFooCreate,
		Read:   resourceFooRead,
		Update: resourceFooUpdate,
		Delete: resourceFooDelete,

		Schema: map[string]*schema.Schema{
			"fookey": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceFooCreate(d *schema.ResourceData, m interface{}) error {
	fookey := d.Get("fookey").(string)
	projectID := d.Get("project_id").(string)

	d.Set("project_id", projectID)
	d.SetId(fookey)
	return resourceFooRead(d, m)
}

func resourceFooRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceFooUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceFooRead(d, m)
}

func resourceFooDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
