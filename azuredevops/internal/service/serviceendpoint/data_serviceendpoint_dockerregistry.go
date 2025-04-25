package serviceendpoint

import (
	"maps"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataResourceServiceEndpointDockerRegistry() *schema.Resource {
	r := &schema.Resource{
		Read: dataResourceServiceEndpointDockerRegistryRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: dataSourceGenBaseSchema(),
	}

	maps.Copy(r.Schema, map[string]*schema.Schema{
		"docker_registry": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"docker_username": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"docker_password": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"docker_email": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"registry_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
	})
	return r
}

func dataResourceServiceEndpointDockerRegistryRead(d *schema.ResourceData, m interface{}) error {
	serviceEndpoint, err := dataSourceGetBaseServiceEndpoint(d, m)
	if err != nil {
		return err
	}

	if serviceEndpoint != nil && serviceEndpoint.Id != nil {
		if err = checkServiceConnection(serviceEndpoint); err != nil {
			return err
		}

		doBaseFlattening(d, serviceEndpoint)
		if serviceEndpoint.Authorization != nil {
			if serviceEndpoint.Authorization.Parameters != nil {
				if v, ok := (*serviceEndpoint.Authorization.Parameters)["registry"]; ok {
					d.Set("docker_registry", v)
				}
				if v, ok := (*serviceEndpoint.Authorization.Parameters)["email"]; ok {
					d.Set("docker_email", v)
				}
				if v, ok := (*serviceEndpoint.Authorization.Parameters)["username"]; ok {
					d.Set("docker_username", v)
				}
			}
		}
		if serviceEndpoint.Data != nil {
			if v, ok := (*serviceEndpoint.Data)["registrytype"]; ok {
				d.Set("registry_type", v)
			}
		}
	}
	return nil
}
