package serviceendpoint

import (
	"fmt"
	"maps"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
	"gopkg.in/yaml.v3"
)

// ResourceServiceEndpointKubernetes schema and implementation for kubernetes service endpoint resource
func ResourceServiceEndpointKubernetes() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointKubernetesCreate,
		Read:   resourceServiceEndpointKubernetesRead,
		Update: resourceServiceEndpointKubernetesUpdate,
		Delete: resourceServiceEndpointKubernetesDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Importer: tfhelper.ImportProjectQualifiedResourceUUID(),
		Schema:   baseSchema(),
	}

	maps.Copy(r.Schema, map[string]*schema.Schema{
		"apiserver_url": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "URL to Kubernete's API-Server",
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"authorization_type": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Type of credentials to use",
			ValidateFunc: validation.StringInSlice([]string{"AzureSubscription", "Kubeconfig", "ServiceAccount"}, false),
		},

		"azure_subscription": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "'AzureSubscription'-type of configuration",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"azure_environment": {
						Type:         schema.TypeString,
						Optional:     true,
						Default:      "AzureCloud",
						Description:  "type of azure cloud: AzureCloud",
						ValidateFunc: validation.StringInSlice([]string{"AzureCloud"}, false),
					},
					"cluster_name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "name of aks-resource",
					},
					"subscription_id": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "id of azure subscription",
						ValidateFunc: validation.IsUUID,
					},
					"subscription_name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "name of azure subscription",
					},
					"tenant_id": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "id of aad-tenant",
						ValidateFunc: validation.IsUUID,
					},
					"resourcegroup_id": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "id of resourcegroup",
					},
					"namespace": {
						Type:        schema.TypeString,
						Optional:    true,
						Default:     "default",
						Description: "accessed namespace",
					},
					"cluster_admin": {
						Type:        schema.TypeBool,
						Optional:    true,
						ForceNew:    true,
						Default:     false,
						Description: "Enable Cluster Admin",
					},
				},
			},
		},

		"kubeconfig": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Description: "'Kubeconfig'-type of configuration",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kube_config": {
						Type:        schema.TypeString,
						Required:    true,
						Sensitive:   true,
						DefaultFunc: schema.EnvDefaultFunc("AZDO_KUBERNETES_SERVICE_CONNECTION_KUBECONFIG", nil),
						Description: "Content of the kubeconfig file. The configuration information in your kubeconfig file allows Kubernetes clients to talk to your Kubernetes API servers. This file is used by kubectl and all supported Kubernetes clients.",
					},
					"cluster_context": {
						Type:        schema.TypeString,
						Optional:    true,
						Computed:    true,
						Description: "Context of your cluster",
					},
					"accept_untrusted_certs": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     true,
						Description: "Enable this if your authentication uses untrusted certificates",
					},
				},
			},
		},

		"service_account": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "'ServiceAccount'-type of configuration",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ca_cert": {
						Type:         schema.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						DefaultFunc:  schema.EnvDefaultFunc("AZDO_KUBERNETES_SERVICE_CONNECTION_SERVICE_ACCOUNT_CERT", nil),
						Description:  "Secret cert",
					},
					"token": {
						Type:         schema.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						DefaultFunc:  schema.EnvDefaultFunc("AZDO_KUBERNETES_SERVICE_CONNECTION_SERVICE_ACCOUNT_TOKEN", nil),
						Description:  "Secret token",
					},
					"accept_untrusted_certs": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},
	})

	return r
}

func resourceServiceEndpointKubernetesCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointKubernetes(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointKubernetesRead(d, m)
}

func resourceServiceEndpointKubernetesRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	getArgs, err := serviceEndpointGetArgs(d)
	if err != nil {
		return err
	}

	serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(clients.Ctx, *getArgs)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf(" looking up service endpoint given ID (%v) and project ID (%v): %v", getArgs.EndpointId, getArgs.Project, err)
	}

	if err = checkServiceConnection(serviceEndpoint); err != nil {
		return err
	}

	doBaseFlattening(d, serviceEndpoint)
	if err = flattenServiceEndpointKubernetes(d, serviceEndpoint); err != nil {
		return err
	}
	return nil
}

func resourceServiceEndpointKubernetesUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointKubernetes(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	if _, err = updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return fmt.Errorf(" updating service endpoint in Azure DevOps: %+v", err)
	}

	return resourceServiceEndpointKubernetesRead(d, m)
}

func resourceServiceEndpointKubernetesDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointKubernetes(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete))
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointKubernetes(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpoint := doBaseExpansion(d)
	serviceEndpoint.Type = converter.String("kubernetes")
	serviceEndpoint.Url = converter.String(d.Get("apiserver_url").(string))

	switch d.Get("authorization_type").(string) {
	case "AzureSubscription":
		configurationRaw := d.Get("azure_subscription").(*schema.Set).List()
		configuration := configurationRaw[0].(map[string]interface{})
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"azureEnvironment": configuration["azure_environment"].(string),
				"azureTenantId":    configuration["tenant_id"].(string),
			},
			Scheme: converter.String("Kubernetes"),
		}

		clusterID := fmt.Sprintf("/subscriptions/%s/resourcegroups/%s/providers/Microsoft.ContainerService/managedClusters/%s", configuration["subscription_id"].(string), configuration["resourcegroup_id"].(string), configuration["cluster_name"].(string))
		serviceEndpoint.Data = &map[string]string{
			"authorizationType":     "AzureSubscription",
			"azureSubscriptionId":   configuration["subscription_id"].(string),
			"azureSubscriptionName": configuration["subscription_name"].(string),
			"clusterId":             clusterID,
			"namespace":             configuration["namespace"].(string),
			"clusterAdmin":          strconv.FormatBool(configuration["cluster_admin"].(bool)),
		}
	case "Kubeconfig":
		configurationRaw := d.Get("kubeconfig").([]interface{})
		configuration := configurationRaw[0].(map[string]interface{})

		clusterContextInput := configuration["cluster_context"].(string)
		if clusterContextInput == "" {
			kubeConfigYAML := configuration["kube_config"].(string)
			var kubeConfigYAMLUnmarshalled map[string]interface{}
			err := yaml.Unmarshal([]byte(kubeConfigYAML), &kubeConfigYAMLUnmarshalled)
			if err != nil {
				errResult := fmt.Errorf("kube_config contains an invalid YAML: %s", err)
				return nil, errResult
			}
			if v, ok := kubeConfigYAMLUnmarshalled["contexts"]; ok {
				if rawConfig, ok := v.([]interface{}); ok && len(rawConfig) > 0 {
					clusterContextInputList := rawConfig[0].(map[string]interface{})
					if name, exist := clusterContextInputList["name"]; exist {
						clusterContextInput = name.(string)
					}
				}
			}
		}

		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"clusterContext": clusterContextInput,
				"kubeconfig":     configuration["kube_config"].(string),
			},
			Scheme: converter.String("Kubernetes"),
		}

		serviceEndpoint.Data = &map[string]string{
			"authorizationType":    "Kubeconfig",
			"acceptUntrustedCerts": fmt.Sprintf("%v", configuration["accept_untrusted_certs"].(bool)),
		}
	case "ServiceAccount":
		configurationRaw := d.Get("service_account").([]interface{})
		configuration := configurationRaw[0].(map[string]interface{})

		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"apiToken":                  configuration["token"].(string),
				"serviceAccountCertificate": configuration["ca_cert"].(string),
			},
			Scheme: converter.String("Token"),
		}

		serviceEndpoint.Data = &map[string]string{
			"acceptUntrustedCerts": strconv.FormatBool(configuration["accept_untrusted_certs"].(bool)),
			"authorizationType":    "ServiceAccount",
		}
	}

	return serviceEndpoint, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointKubernetes(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) error {
	if serviceEndpoint.Data != nil {
		if v, ok := (*serviceEndpoint.Data)["authorizationType"]; ok {
			d.Set("authorization_type", v)
		}
	}

	if serviceEndpoint.Url != nil {
		d.Set("apiserver_url", *serviceEndpoint.Url)
	}

	switch (*serviceEndpoint.Data)["authorizationType"] {
	case "AzureSubscription":
		clusterIDSplit := strings.Split((*serviceEndpoint.Data)["clusterId"], "/")
		var clusterNameIndex int
		var resourceGroupIDIndex int
		for k, v := range clusterIDSplit {
			if v == "resourcegroups" {
				resourceGroupIDIndex = k + 1
			}
			if v == "managedClusters" {
				clusterNameIndex = k + 1
			}
		}
		clusterAdmin, err := strconv.ParseBool((*serviceEndpoint.Data)["clusterAdmin"])
		if err != nil {
			return fmt.Errorf(" Parsing `cluster_admin` value. Error: %+v", err)
		}
		configItems := map[string]interface{}{
			"azure_environment": (*serviceEndpoint.Authorization.Parameters)["azureEnvironment"],
			"tenant_id":         (*serviceEndpoint.Authorization.Parameters)["azureTenantId"],
			"subscription_id":   (*serviceEndpoint.Data)["azureSubscriptionId"],
			"subscription_name": (*serviceEndpoint.Data)["azureSubscriptionName"],
			"cluster_name":      clusterIDSplit[clusterNameIndex],
			"resourcegroup_id":  clusterIDSplit[resourceGroupIDIndex],
			"namespace":         (*serviceEndpoint.Data)["namespace"],
			"cluster_admin":     clusterAdmin,
		}
		configItemList := make([]map[string]interface{}, 1)
		configItemList[0] = configItems

		d.Set("azure_subscription", configItemList)
	case "Kubeconfig":
		var kubeconfig map[string]interface{}
		kubeconfigSet := d.Get("kubeconfig").([]interface{})
		configuration := kubeconfigSet[0].(map[string]interface{})

		if len(configuration) > 0 {
			kubeconfig = map[string]interface{}{}

			if v, ok := configuration["kube_config"]; ok {
				kubeconfig["kube_config"] = v.(string)
			}

			if serviceEndpoint.Data != nil {
				if v, ok := (*serviceEndpoint.Data)["acceptUntrustedCerts"]; ok {
					acceptUntrustedCerts, err := strconv.ParseBool(v)
					if err != nil {
						return fmt.Errorf(" failed to parse `accept_untrusted_certs`: %+v ", err)
					}
					kubeconfig["accept_untrusted_certs"] = acceptUntrustedCerts
				}
			}

			if serviceEndpoint.Authorization != nil && serviceEndpoint.Authorization.Parameters != nil {
				if v, ok := (*serviceEndpoint.Authorization.Parameters)["clusterContext"]; ok {
					kubeconfig["cluster_context"] = v
				}
			}

			kubeconfigList := make([]map[string]interface{}, 1)
			kubeconfigList[0] = kubeconfig
			d.Set("kubeconfig", kubeconfigList)
		}
	case "ServiceAccount":
		var serviceAccount map[string]interface{}
		serviceAccountSet := d.Get("service_account").([]interface{})

		if len(serviceAccountSet) == 0 {
			serviceAccount = map[string]interface{}{
				"token":   "",
				"ca_cert": "",
			}
		} else {
			configuration := serviceAccountSet[0].(map[string]interface{})
			serviceAccount = map[string]interface{}{
				"token":   configuration["token"].(string),
				"ca_cert": configuration["ca_cert"].(string),
			}
			if v, ok := (*serviceEndpoint.Data)["acceptUntrustedCerts"]; ok {
				acceptUntrustedCerts, err := strconv.ParseBool(v)
				if err != nil {
					return fmt.Errorf(" Pparse `accept_untrusted_certs`. Error: %+v ", err)
				}
				serviceAccount["accept_untrusted_certs"] = acceptUntrustedCerts
			}
		}

		serviceAccountList := make([]map[string]interface{}, 1)
		serviceAccountList[0] = serviceAccount
		d.Set("service_account", serviceAccountList)
	}
	return nil
}
