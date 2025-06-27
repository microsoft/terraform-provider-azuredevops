package taskagent

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func DataAgentPool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAgentPoolRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"pool_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_provision": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"auto_update": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceAgentPoolRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	poolName := d.Get("name").(string)

	agentPools, err := clients.TaskAgentClient.GetAgentPools(clients.Ctx, taskagent.GetAgentPoolsArgs{
		PoolName: converter.String(poolName),
	})
	if err != nil {
		return err
	}

	if len(*agentPools) > 1 {
		return fmt.Errorf("Found multiple agent pools for name: %s. Agent pools found: %+v", poolName, agentPools)
	}

	if len(*agentPools) == 0 {
		return fmt.Errorf("Unable to find agent pool with name: %s", poolName)
	}

	pool := (*agentPools)[0]

	d.SetId(strconv.Itoa(*pool.Id))
	if pool.Name != nil {
		d.Set("name", pool.Name)
	}

	if pool.PoolType != nil {
		d.Set("pool_type", pool.PoolType)
	}

	if pool.AutoProvision != nil {
		d.Set("auto_provision", pool.AutoProvision)
	}

	if pool.AutoUpdate != nil {
		d.Set("auto_update", pool.AutoUpdate)
	}
	return nil
}
