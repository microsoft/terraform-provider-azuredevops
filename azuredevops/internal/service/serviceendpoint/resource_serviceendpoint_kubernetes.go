package serviceendpoint

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
	"gopkg.in/yaml.v3"
)

const (
	resourceAttrAuthType            = "authorization_type"
	resourceAttrAPIURL              = "apiserver_url"
	resourceBlockAzSubscription     = "azure_subscription"
	resourceBlockKubeconfig         = "kubeconfig"
	resourceBlockServiceAccount     = "service_account"
	serviceEndpointDataAttrAuthType = "authorizationType"
)

func makeSchemaAzureSubscription(r *schema.Resource) {
	r.Schema[resourceBlockAzSubscription] = &schema.Schema{
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
	}
}

func makeSchemaKubeconfig(r *schema.Resource) {
	resourceElemSchema := &schema.Resource{
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
	}
	r.Schema[resourceBlockKubeconfig] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MinItems:    1,
		MaxItems:    1,
		Description: "'Kubeconfig'-type of configuration",
		Elem:        resourceElemSchema,
	}
}

func makeSchemaServiceAccount(r *schema.Resource) {
	resourceElemSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ca_cert": {
				Type:         schema.TypeString,
				Sensitive:    true,
				Optional:     true,
				Description:  "Service account certificate",
				ValidateFunc: validation.StringIsNotWhiteSpace,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
			},
			"token": {
				Type:         schema.TypeString,
				Sensitive:    true,
				Optional:     true,
				Description:  "Token",
				ValidateFunc: validation.StringIsNotWhiteSpace,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
			},
		},
	}
	makeProtectedSchema(resourceElemSchema, "ca_cert", "AZDO_KUBERNETES_SERVICE_CONNECTION_SERVICE_ACCOUNT_CERT", "Secret cert")
	makeProtectedSchema(resourceElemSchema, "token", "AZDO_KUBERNETES_SERVICE_CONNECTION_SERVICE_ACCOUNT_TOKEN", "Secret token")
	r.Schema[resourceBlockServiceAccount] = &schema.Schema{
		Type:        schema.TypeList,
		MaxItems:    1,
		Optional:    true,
		Description: "'ServiceAccount'-type of configuration",
		Elem:        resourceElemSchema,
	}
}

// ResourceServiceEndpointKubernetes schema and implementation for kubernetes service endpoint resource
func ResourceServiceEndpointKubernetes() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointKubernetes, expandServiceEndpointKubernetes)
	r.Schema[resourceAttrAPIURL] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Description:  "URL to Kubernete's API-Server",
		ValidateFunc: validation.IsURLWithHTTPorHTTPS,
	}
	r.Schema[resourceAttrAuthType] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Description:  "Type of credentials to use",
		ValidateFunc: validation.StringInSlice([]string{"AzureSubscription", "Kubeconfig", "ServiceAccount"}, false),
	}
	makeSchemaAzureSubscription(r)
	makeSchemaKubeconfig(r)
	makeSchemaServiceAccount(r)

	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointKubernetes(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Type = converter.String("kubernetes")
	serviceEndpoint.Url = converter.String(d.Get(resourceAttrAPIURL).(string))

	switch d.Get(resourceAttrAuthType).(string) {
	case "AzureSubscription":
		configurationRaw := d.Get(resourceBlockAzSubscription).(*schema.Set).List()
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
		configurationRaw := d.Get(resourceBlockKubeconfig).([]interface{})
		configuration := configurationRaw[0].(map[string]interface{})

		clusterContextInput := configuration["cluster_context"].(string)
		if clusterContextInput == "" {
			kubeConfigYAML := configuration["kube_config"].(string)
			var kubeConfigYAMLUnmarshalled map[string]interface{}
			err := yaml.Unmarshal([]byte(kubeConfigYAML), &kubeConfigYAMLUnmarshalled)
			if err != nil {
				errResult := fmt.Errorf("kube_config contains an invalid YAML: %s", err)
				return nil, nil, errResult
			}
			clusterContextInputList := kubeConfigYAMLUnmarshalled["contexts"].([]interface{})[0].(map[string]interface{})
			clusterContextInput = clusterContextInputList["name"].(string)
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
		configurationRaw := d.Get(resourceBlockServiceAccount).([]interface{})
		configuration := configurationRaw[0].(map[string]interface{})

		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"apiToken":                  configuration["token"].(string),
				"serviceAccountCertificate": configuration["ca_cert"].(string),
			},
			Scheme: converter.String("Token"),
		}

		serviceEndpoint.Data = &map[string]string{
			"authorizationType": "ServiceAccount",
		}
	}

	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointKubernetes(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	d.Set(resourceAttrAuthType, (*serviceEndpoint.Data)[serviceEndpointDataAttrAuthType])
	d.Set(resourceAttrAPIURL, (*serviceEndpoint.Url))

	switch (*serviceEndpoint.Data)[serviceEndpointDataAttrAuthType] {
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
		clusterAdmin, _ := strconv.ParseBool((*serviceEndpoint.Data)["clusterAdmin"])
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

		d.Set(resourceBlockAzSubscription, configItemList)
	case "Kubeconfig":
		var kubeconfig map[string]interface{}
		kubeconfigSet := d.Get("kubeconfig").([]interface{})

		configuration := kubeconfigSet[0].(map[string]interface{})
		acceptUntrustedCerts, _ := strconv.ParseBool((*serviceEndpoint.Data)["acceptUntrustedCerts"])
		kubeconfig = map[string]interface{}{
			"kube_config":            configuration["kube_config"].(string),
			"cluster_context":        (*serviceEndpoint.Authorization.Parameters)["clusterContext"],
			"accept_untrusted_certs": acceptUntrustedCerts,
		}

		kubeconfigList := make([]map[string]interface{}, 1)
		kubeconfigList[0] = kubeconfig
		d.Set(resourceBlockKubeconfig, kubeconfigList)
	case "ServiceAccount":
		var serviceAccount map[string]interface{}
		serviceAccountSet := d.Get("service_account").([]interface{})

		if len(serviceAccountSet) == 0 {
			newHashToken, hashKeyToken := tfhelper.HelpFlattenSecretNested(d, resourceBlockServiceAccount, nil, "token")
			newHashCert, hashKeyCert := tfhelper.HelpFlattenSecretNested(d, resourceBlockServiceAccount, nil, "ca_cert")
			serviceAccount = map[string]interface{}{
				"token":      "",
				"ca_cert":    "",
				hashKeyToken: newHashToken,
				hashKeyCert:  newHashCert,
			}
		} else {
			configuration := serviceAccountSet[0].(map[string]interface{})
			newHashToken, hashKeyToken := tfhelper.HelpFlattenSecretNested(d, resourceBlockServiceAccount, configuration, "token")
			newHashCert, hashKeyCert := tfhelper.HelpFlattenSecretNested(d, resourceBlockServiceAccount, configuration, "ca_cert")
			serviceAccount = map[string]interface{}{
				"token":      configuration["token"].(string),
				"ca_cert":    configuration["ca_cert"].(string),
				hashKeyToken: newHashToken,
				hashKeyCert:  newHashCert,
			}
		}

		serviceAccountList := make([]map[string]interface{}, 1)
		serviceAccountList[0] = serviceAccount
		d.Set(resourceBlockServiceAccount, serviceAccountList)
	}
}
