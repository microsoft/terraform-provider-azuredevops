package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
)

func resourceServiceEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceServiceEndpointCreate,
		Read:   resourceServiceEndpointRead,
		Update: resourceServiceEndpointUpdate,
		Delete: resourceServiceEndpointDelete,

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"service_endpoint_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"service_endpoint_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"service_endpoint_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"service_endpoint_owner": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"github_service_endpoint_pat": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AZDO_GITHUB_SERVICE_CONNECTION_PAT", nil),
				Description: "The GitHub personal access token which should be used.",
			},
		},
	}
}

type serviceEndpointValues struct {
	projectID                   string
	serviceEndpointName         string
	serviceEndpointType         string
	serviceEndpointURL          string
	serviceEndpointOwner        string
	githubServiceEndpointPAT    string
	githubServiceEndpointScheme string
	serviceEndpointID           string
}

func resourceServiceEndpointCreate(d *schema.ResourceData, m interface{}) error {
	// instantiate client
	clients := m.(*aggregatedClient)

	values := serviceEndpointValues{
		projectID:                   d.Get("project_id").(string),
		serviceEndpointName:         d.Get("service_endpoint_name").(string),
		serviceEndpointType:         d.Get("service_endpoint_type").(string),
		serviceEndpointURL:          d.Get("service_endpoint_url").(string),
		serviceEndpointOwner:        d.Get("service_endpoint_owner").(string),
		githubServiceEndpointPAT:    d.Get("github_service_endpoint_pat").(string),
		githubServiceEndpointScheme: "PersonalAccessToken",
	}

	if _, err := serviceEndpointCreate(clients, &values); err != nil {
		return fmt.Errorf("Error creating service endpoint: %+v", err)
	}

	serviceEndpointID, err := lookupServiceEndpointID(clients, values.projectID, values.serviceEndpointName)

	if err != nil {
		return fmt.Errorf("Error retrieving service endpoint: %+v", err)
	}

	values.serviceEndpointID = serviceEndpointID

	d.Set("service_endpoint_id", values.serviceEndpointID)
	d.SetId(values.serviceEndpointID)

	// read service endpoint and return
	return resourceServiceEndpointRead(d, m)
}

func resourceServiceEndpointRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServiceEndpointUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServiceEndpointDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func serviceEndpointCreate(clients *aggregatedClient, values *serviceEndpointValues) (string, error) {
	serviceEndpointRef, err := clients.ServiceEndpointClient.CreateServiceEndpoint(clients.ctx, serviceendpoint.CreateServiceEndpointArgs{
		Endpoint: &serviceendpoint.ServiceEndpoint{
			Name:  &values.serviceEndpointName,
			Type:  &values.serviceEndpointType,
			Url:   &values.serviceEndpointURL,
			Owner: &values.serviceEndpointOwner,
			Authorization: &serviceendpoint.EndpointAuthorization{
				Parameters: &map[string]string{
					"accessToken": values.githubServiceEndpointPAT,
				},
				Scheme: &values.githubServiceEndpointScheme,
			},
		},
		Project: &values.projectID,
	})
	if err != nil {
		return "", err
	}
	return uuid.UUID.String(*serviceEndpointRef.Id), nil
}

func lookupServiceEndpointID(clients *aggregatedClient, projectID string, serviceEndpointName string) (string, error) {
	serviceEndpoints, err := clients.ServiceEndpointClient.GetServiceEndpoints(clients.ctx, serviceendpoint.GetServiceEndpointsArgs{
		Project: &projectID,
	})
	if err != nil {
		return "", err
	}

	for _, serviceEndpoint := range *serviceEndpoints {
		if *serviceEndpoint.Name == serviceEndpointName {
			return serviceEndpoint.Id.String(), nil
		}
	}

	return "", fmt.Errorf("No service endpoint found")
}
