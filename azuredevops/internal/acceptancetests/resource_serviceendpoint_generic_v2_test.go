package acceptancetests

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// Tests basic functionality of the Generic Service Endpoint V2 resource
func TestAccServiceEndpointGenericV2_Basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	serviceEndpointType := "github"

	resourceType := "azuredevops_serviceendpoint_generic_v2"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclServiceEndpointGenericV2TokenBasic(projectName, serviceEndpointName, serviceEndpointType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "type", serviceEndpointType),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://github.com"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "id"),
				),
			},
			{
				ResourceName:            tfSvcEpNode,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"authorization"},
			},
		},
	})
}

// Tests if the Generic Service Endpoint V2 can be updated with a different server_url
func TestAccServiceEndpointGenericV2_Update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	serviceEndpointType := "github"
	serverUrlInitial := "https://github.com"
	serverUrlUpdated := "https://api.github.com"

	resourceType := "azuredevops_serviceendpoint_generic_v2"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclServiceEndpointGenericV2TokenCustomUrl(projectName, serviceEndpointName, serviceEndpointType, serverUrlInitial),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "type", serviceEndpointType),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", serverUrlInitial),
				),
			},
			{
				Config: hclServiceEndpointGenericV2TokenCustomUrl(projectName, serviceEndpointName, serviceEndpointType, serverUrlUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "type", serviceEndpointType),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", serverUrlUpdated),
				),
			},
		},
	})
}

// Tests username/password authentication for Generic Service Endpoint V2
func TestAccServiceEndpointGenericV2_UsernamePassword(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	serviceEndpointType := "generic"
	username := "testuser"
	password := "testpass"

	resourceType := "azuredevops_serviceendpoint_generic_v2"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclServiceEndpointGenericV2UsernamePassword(projectName, serviceEndpointName, serviceEndpointType, username, password),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "type", serviceEndpointType),
				),
			},
			{
				ResourceName:            tfSvcEpNode,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"authorization"},
			},
		},
	})
}

// Tests certificate authentication for Generic Service Endpoint V2
func TestAccServiceEndpointGenericV2_Certificate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	serviceEndpointType := "generic"
	certificate := "-----BEGIN CERTIFICATE-----\nMIIFvTCCA6WgAwIBAgIUNTIwYTU2MGU4MjFjNDBiMWI4YzczZDAwDQYJKoZIhvcNAQELBQAwbjEL\nMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRYwFAYDVQQHDA1TYW4gRnJhbmNpc2NvMRMwEQYDVQQK\nDApFeGFtcGxlIENvMRMwEQYDVQQLDApFeGFtcGxlIE9VMRAwDgYDVQQDDAdUZXN0IENBMB4XDTIz\nMDcwNzA3NDY1MVoXDTMzMDcwNDA3NDY1MVowbjELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRYw\nFAYDVQQHDA1TYW4gRnJhbmNpc2NvMRMwEQYDVQQKDApFeGFtcGxlIENvMRMwEQYDVQQLDApFeGFt\ncGxlIE9VMRAwDgYDVQQDDAdUZXN0IENBMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA\nxRkI7bP7lAzImKagUVte9XYeGV70GvV4oMLPmWkBMJVGNHlvA+4VbvHQcJXXw6S0EUbxGQT5tDgp\nu9t+YYzQPLQiUJp2YNvz93dWVHCsEQbSVXeSxChK9zSKNTGoqJQTZF71Cy9qoN87QqHHy8TtNgKf\nHS7+Vz/xJ8HBZUdJ8XLtPm/xKDROXWGxSplAjFBJ8OGGJuZOkKVQmji4kJVMQh8KQEkQYJJVGPGt\nIPxlbFXM/1HakLNMcOSFmfK7v0J4+yi1YXF3Gyej5UGrKQozZ6PO9PRoPwOwy7jaBKDMQjWF5+jK\nCBKBGNzg0ENSwBTPfN6wO/mrRbGGRavZ56e0j6TH0CbgTCIpKhG3ivMTKHMUbBrvsLx0QHGGMsMl\nNpbXt/mGBbKSIr5ZwOLtn5OqKvVFSEHCEi5yrpotf0qedv3d6HCLH+mPkI7zdrIo+JiKDVPTt99y\nXXF6Kq75Ah9bYywbHYXJAMRs5CjLI4iXOiygUSR5MPv2zPMrXbD9vdPymjMYCheSlP/HcQiwJAJ0\nZ/silaNI3GfKCJDPbZb0+kdKXmdI55VIh9SWo6hTkXVmXjrMobp97jFnUOyKOmXLHHMaFTdVGJXW\nscvIJLmZN0XJpJqWSKK4ypypu5RFhFRWDiF1qyFmY0O/pIizUZcnvDMCAwEAAaNTMFEwHQYDVR0O\nBBYEFHDGFfIJYgPjs4brmJwjQT+D9VT+MB8GA1UdIwQYMBaAFHDGFfIJYgPjs4brmJwjQT+D9VT+\nMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggIBALIAlXwsPTvBXBbZ0A6NoJgJ5V5r\n2E/EcTNhfJzHEVdO+bQSxJTWNfA7HLJ52H21vSdBEe1S7wFuFjU9WEXfYZcOAxsAXfQN8l1I+cAd\nUtwuGlEcVr8nWuyhOlQl9DJ9gUPiqiCz0MTI7/RzT7qwzw5rjcLTSl0gQGwR3QTzWZQBIfz3uIwD\nDt0wC3JEErgKw+AWoMHMZbj7MjJj8lBfv9GwLNpfzEOFpNh5eFteUiY5Ir+wQwVZkCHdLEPMC+jI\ne3ZQIVxKIk0HOvKQvTJl6PRWFSUzCr85WKR9rG4GjnI5LQRFMaYcmDQXQQE+WR5UbWA3d9HsnnZO\neDlTBKm6QdkWQbRGPgOQIXZXnsJnGvMZ2zHQUBtE0Am4hZy2VWpdxd7xFYUIfWBlsN1RUY3aSYq4\neCPnKvIFmWijPMM8jvpTD/rSZ6W0T+5YLVn3DVr0HBhL+HbNtHl1pWhCjdG+signllPZQpNisXf6\nLbJI9MYsDZD6D4Eli0l6nScvMQX0Nd0jfF0GISKCJlcEtWlEDJX9zC7BVKSkLUnGTgP5Ldgr9cRN\nJ8H1wHDyJI2emZbFzXnLvl/pxpGEPQAYijW5qzaNnw2FOcxdKJUuP0LfljZwoJ9h+/xNRZ2jVFrF\ne9h5JK+B4walV6XIT6nD\n-----END CERTIFICATE-----"
	certPassword := "certpass123"

	resourceType := "azuredevops_serviceendpoint_generic_v2"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclServiceEndpointGenericV2Certificate(projectName, serviceEndpointName, serviceEndpointType, certificate, certPassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "type", serviceEndpointType),
				),
			},
			{
				ResourceName:            tfSvcEpNode,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"authorization"},
			},
		},
	})
}

// Tests additional data configuration for Generic Service Endpoint V2
func TestAccServiceEndpointGenericV2_WithData(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	serviceEndpointType := "generic"

	resourceType := "azuredevops_serviceendpoint_generic_v2"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclServiceEndpointGenericV2WithData(projectName, serviceEndpointName, serviceEndpointType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "data.environment", "test"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "data.region", "us-west-1"),
				),
			},
			{
				ResourceName:            tfSvcEpNode,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"authorization"},
			},
		},
	})
}

// Tests resource validation for type in Generic Service Endpoint V2
func TestAccServiceEndpointGenericV2_ValidateServiceEndpointType(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	serviceEndpointType := "invalidtype" // This should fail validation

	resourceType := "azuredevops_serviceendpoint_generic_v2"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config:      hclServiceEndpointGenericV2TokenBasic(projectName, serviceEndpointName, serviceEndpointType),
				ExpectError: validateServiceEndpointTypeError(serviceEndpointType),
			},
		},
	})
}

// Helper function to validate service endpoint type error
func validateServiceEndpointTypeError(serviceEndpointType string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("service endpoint type '%s' is not available", serviceEndpointType))
}

// Helper function to generate HCL for a generic service endpoint with token auth
func hclServiceEndpointGenericV2TokenBasic(projectName string, serviceEndpointName string, serviceEndpointType string) string {
	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_generic_v2" "test" {
  project_id            = azuredevops_project.project.id
  name = "%s"
  description           = "Managed by Terraform"
  type = "%s"
  server_url            = "https://github.com"

  authorization {
    scheme = "Token"
    parameters = {
      token      = "test-token"
      token_type = "PersonalAccessToken"
    }
  }
}`, projectResource, serviceEndpointName, serviceEndpointType)
}

// Helper function to generate HCL for a generic service endpoint with token auth and custom URL
func hclServiceEndpointGenericV2TokenCustomUrl(projectName string, serviceEndpointName string, serviceEndpointType string, serverUrl string) string {
	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_generic_v2" "test" {
  project_id            = azuredevops_project.project.id
  name = "%s"
  description           = "Managed by Terraform"
  type = "%s"
  server_url            = "%s"

  authorization {
    scheme = "Token"
    parameters = {
      token      = "test-token"
      token_type = "PersonalAccessToken"
    }
  }
}`, projectResource, serviceEndpointName, serviceEndpointType, serverUrl)
}

// Helper function to generate HCL for a generic service endpoint with username/password auth
func hclServiceEndpointGenericV2UsernamePassword(projectName string, serviceEndpointName string, serviceEndpointType string, username string, password string) string {
	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_generic_v2" "test" {
  project_id            = azuredevops_project.project.id
  name = "%s"
  description           = "Managed by Terraform"
  type = "%s"
  server_url            = "https://example.com"

  authorization {
    scheme = "UsernamePassword"
    parameters = {
      username = "%s"
      password = "%s"
    }
  }
}`, projectResource, serviceEndpointName, serviceEndpointType, username, password)
}

// Helper function to generate HCL for a generic service endpoint with certificate auth
func hclServiceEndpointGenericV2Certificate(projectName string, serviceEndpointName string, serviceEndpointType string, certificate string, certPassword string) string {
	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_generic_v2" "test" {
  project_id            = azuredevops_project.project.id
  name = "%s"
  description           = "Managed by Terraform"
  type = "%s"
  server_url            = "https://example.com"

  authorization {
    scheme = "Certificate"
    parameters = {
      certificate          = <<EOF
%s
EOF
      certificate_password = "%s"
    }
  }
}`, projectResource, serviceEndpointName, serviceEndpointType, certificate, certPassword)
}

// Helper function to generate HCL for a generic service endpoint with additional data
func hclServiceEndpointGenericV2WithData(projectName string, serviceEndpointName string, serviceEndpointType string) string {
	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_generic_v2" "test" {
  project_id            = azuredevops_project.project.id
  name = "%s"
  description           = "Managed by Terraform"
  type = "%s"
  server_url            = "https://example.com"

  authorization {
    scheme = "Token"
    parameters = {
      token      = "test-token"
      token_type = "PersonalAccessToken"
    }
  }

  data = {
    environment = "test"
    region      = "us-west-1"
  }
}`, projectResource, serviceEndpointName, serviceEndpointType)
}

// checkServiceEndpointDestroyed verifies that all service endpoints with the specified type have been destroyed
func checkServiceEndpointDestroyed(resourceType string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

		// verify that every service endpoint referenced in the state does not exist in AzDO
		for _, resource := range s.RootModule().Resources {
			if resource.Type != resourceType {
				continue
			}

			endpointID := resource.Primary.ID
			projectID := resource.Primary.Attributes["project_id"]

			// Ensure the service endpoint does not exist
			_, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(
				clients.Ctx,
				serviceendpoint.GetServiceEndpointDetailsArgs{
					EndpointId: converter.String(endpointID),
					Project:    converter.String(projectID),
				},
			)

			if err == nil {
				return fmt.Errorf("Service Endpoint ID %s still exists", endpointID)
			}
		}

		return nil
	}
}
