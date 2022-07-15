package serviceendpoint

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func DataServiceEndpointAzureRM() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServiceEndpointAzureRMRead, // Implement below
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"service_endpoint_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"service_endpoint_name", "id"},
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"id": {
				Description:  "The ID of the serviceendpoint",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"service_endpoint_name", "id"},
				ValidateFunc: validation.IsUUID,
			},

			"authorization": {
				Type:         schema.TypeMap,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"azurerm_management_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"azurerm_management_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"azurerm_subscription_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"azurerm_subscription_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"resource_group": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"azurerm_spn_tenantid": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceServiceEndpointAzureRMRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	var projectID *uuid.UUID
	projectIDString := d.Get("project_id").(string)
	parsedProjectID, err := uuid.Parse(projectIDString)
	if err != nil {
		return fmt.Errorf("Error parsing projectID from the Terraform data source declaration: %v", err)
	}

	projectID = &parsedProjectID

	if serviceEndpointIDString, ok := d.GetOk("id"); ok {
		var serviceEndpointID *uuid.UUID
		parsedServiceEndpointID, err := uuid.Parse(serviceEndpointIDString.(string))
		if err != nil {
			return fmt.Errorf("Error parsing serviceEndpointID from the Terraform data source declaration: %v", err)
		}
		serviceEndpointID = &parsedServiceEndpointID

		serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(
			clients.Ctx,
			serviceendpoint.GetServiceEndpointDetailsArgs{
				EndpointId: serviceEndpointID,
				Project:    converter.String(projectID.String()),
			},
		)
		if err != nil {
			if utils.ResponseWasNotFound(err) {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error looking up service endpoint with ID (%v) and projectID (%v): %v", serviceEndpointID, projectID, err)
		}

		if serviceEndpoint.Id == nil {
			d.SetId("")
		} else {
			// prevent credential leak into data source
			(*serviceEndpoint.Data)["creationMode"] = ""
			flattenServiceEndpointAzureRM(d, serviceEndpoint, projectID)
		}
		return nil
	}

	if serviceEndpointName, ok := d.GetOk("service_endpoint_name"); ok {
		// get service endpointdetails by name
		serviceEndpoint, err := getServiceEndpointByNameAndProject(clients, serviceEndpointName.(string), projectID)
		if err != nil {
			if utils.ResponseWasNotFound(err) {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error looking up service endpoint with name (%v) and projectID (%v): %v", serviceEndpointName, projectID, err)
		}

		if serviceEndpoint.Id == nil {
			d.SetId("")
		} else {
			// prevent credential leak into data source
			(*serviceEndpoint.Data)["creationMode"] = ""
			flattenServiceEndpointAzureRM(d, serviceEndpoint, projectID)
		}
		return nil
	}
	return nil
}

func getServiceEndpointByNameAndProject(clients *client.AggregatedClient, serviceEndpointName string, projectID *uuid.UUID) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpointNameList := &[]string{serviceEndpointName}

	serviceEndpoints, err := clients.ServiceEndpointClient.GetServiceEndpointsByNames(
		clients.Ctx,
		serviceendpoint.GetServiceEndpointsByNamesArgs{
			Project:       converter.String(projectID.String()),
			EndpointNames: serviceEndpointNameList,
		},
	)
	if err != nil {
		return nil, err
	}

	if len(*serviceEndpoints) > 1 {
		return nil, fmt.Errorf("%v returns more than one serviceEndpoint!", serviceEndpointName)
	}

	return &(*serviceEndpoints)[0], nil
}
