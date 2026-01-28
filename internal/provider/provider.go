package provider

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/entrauth/aztfauth"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/terraform-provider-azuredevops/internal/adovalidator"
	"github.com/microsoft/terraform-provider-azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/internal/framework"
	"github.com/microsoft/terraform-provider-azuredevops/internal/meta"
	"github.com/microsoft/terraform-provider-azuredevops/internal/services/core"
	"github.com/microsoft/terraform-provider-azuredevops/internal/services/graph"
)

var _ provider.Provider = (*Provider)(nil)

type Provider struct{}

type providerModel struct {
	OrgServiceUrl                types.String `tfsdk:"org_service_url"`
	PersonalAccessToken          types.String `tfsdk:"personal_access_token"`
	ClientID                     types.String `tfsdk:"client_id"`
	ClientIDFilePath             types.String `tfsdk:"client_id_file_path"`
	TenantID                     types.String `tfsdk:"tenant_id"`
	AuxiliaryTenantIDs           types.List   `tfsdk:"auxiliary_tenant_ids"`
	ClientCertificate            types.String `tfsdk:"client_certificate"`
	ClientCertificatePath        types.String `tfsdk:"client_certificate_path"`
	ClientCertificatePassword    types.String `tfsdk:"client_certificate_password"`
	ClientSecret                 types.String `tfsdk:"client_secret"`
	ClientSecretFilePath         types.String `tfsdk:"client_secret_file_path"`
	OIDCRequestToken             types.String `tfsdk:"oidc_request_token"`
	OIDCRequestURL               types.String `tfsdk:"oidc_request_url"`
	OIDCToken                    types.String `tfsdk:"oidc_token"`
	OIDCTokenFilePath            types.String `tfsdk:"oidc_token_file_path"`
	OIDCAzureServiceConnectionID types.String `tfsdk:"oidc_azure_service_connection_id"`
	UseOIDC                      types.Bool   `tfsdk:"use_oidc"`
	UseCLI                       types.Bool   `tfsdk:"use_cli"`
	UseMSI                       types.Bool   `tfsdk:"use_msi"`
}

func New() provider.Provider {
	return &Provider{}
}

func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	model, diags := p.getModel(ctx, req)
	resp.Diagnostics = append(resp.Diagnostics, diags...)
	if diags.HasError() {
		return
	}

	authProvider, err := p.newAuthProvider(model)
	if err != nil {
		resp.Diagnostics.AddError("failed to get auth provider", err.Error())
		return
	}

	client, err := client.New(ctx, authProvider, model.OrgServiceUrl.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to new client", err.Error())
		return
	}

	d := meta.Meta{Client: client}
	resp.DataSourceData = d
	resp.ResourceData = d
	resp.EphemeralResourceData = d
	resp.ActionData = d
	resp.ListResourceData = d
}

func (p *Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "azuredevops"
}

func (p *Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *Provider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		framework.WrapResource(core.NewProjectResource()),
		framework.WrapResource(graph.NewGroupResource()),
		framework.WrapResource(graph.NewGroupMembershipResource()),
	}
}

func (p *Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Azure DevOps Provider",
		Attributes: map[string]schema.Attribute{
			"org_service_url": schema.StringAttribute{
				MarkdownDescription: "The URL of the Azure DevOps instance. This can also be sourced from the `AZDO_ORG_SERVICE_URL` Environment Variable.",
				Optional:            true,
			},
			"personal_access_token": schema.StringAttribute{
				MarkdownDescription: "The personal access token. This can also be sourced from the `AZDO_PERSONAL_ACCESS_TOKEN` Environment Variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "The service principal client id which should be used for AAD auth. This can also be sourced from the `ARM_CLIENT_ID`, `AZURE_CLIENT_ID` Environment Variable.",
				Optional:            true,
				Validators: []validator.String{
					adovalidator.StringIsUUID(),
				},
			},
			"client_id_file_path": schema.StringAttribute{
				MarkdownDescription: "The path to a file containing the Client ID which should be used. This can also be sourced from the `ARM_CLIENT_ID_FILE_PATH` Environment Variable.",
				Optional:            true,
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The service principal tenant id which should be used for AAD auth. This can also be sourced from the `ARM_TENANT_ID` Environment Variable.",
				Optional:            true,
				Validators: []validator.String{
					adovalidator.StringIsUUID(),
				},
			},
			"auxiliary_tenant_ids": schema.ListAttribute{
				MarkdownDescription: "List of auxiliary Tenant IDs required for multi-tenancy and cross-tenant scenarios. This can also be sourced from the `ARM_AUXILIARY_TENANT_IDS` Environment Variable.",
				Optional:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtMost(3),
				},
			},
			"client_certificate_path": schema.StringAttribute{
				MarkdownDescription: "Path to a certificate to use to authenticate to the service principal. This can also be sourced from the `ARM_CLIENT_CERTIFICATE_PATH` Environment Variable.",
				Optional:            true,
			},
			"client_certificate": schema.StringAttribute{
				MarkdownDescription: "Base64 encoded certificate to use to authenticate to the service principal. This can also be sourced from the `ARM_CLIENT_CERTIFICATE` Environment Variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"client_certificate_password": schema.StringAttribute{
				MarkdownDescription: "Password for a client certificate password. This can also be sourced from the `ARM_CLIENT_CERTIFICATE_PASSWORD` Environment Variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "Client secret for authenticating to a service principal. This can also be sourced from the `ARM_CLIENT_SECRET` Environment Variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"client_secret_file_path": schema.StringAttribute{
				MarkdownDescription: "Path to a file containing a client secret for authenticating to a service principal. This can also be sourced from the `ARM_CLIENT_SECRET_FILE_PATH` Environment Variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"oidc_request_token": schema.StringAttribute{
				MarkdownDescription: "The bearer token for the request to the OIDC provider. For use when authenticating as a Service Principal using OpenID Connect. This can also be sourced from the `ARM_OIDC_REQUEST_TOKEN`, `ACTIONS_ID_TOKEN_REQUEST_TOKEN`, or `SYSTEM_ACCESSTOKEN` Environment Variables.",
				Optional:            true,
				Sensitive:           true,
			},
			"oidc_request_url": schema.StringAttribute{
				MarkdownDescription: "The URL for the OIDC provider from which to request an ID token. For use when authenticating as a Service Principal using OpenID Connect. This can also be sourced from the `ARM_OIDC_REQUEST_URL`, `ACTIONS_ID_TOKEN_REQUEST_URL`, or `SYSTEM_OIDCREQUESTURI` Environment Variables",
				Optional:            true,
			},
			"oidc_token": schema.StringAttribute{
				MarkdownDescription: "OIDC token to authenticate as a service principal. This can also be sourced from the `ARM_OIDC_TOKEN` Environment Variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"oidc_token_file_path": schema.StringAttribute{
				MarkdownDescription: "OIDC token from file to authenticate as a service principal. This can also be sourced from the `ARM_OIDC_TOKEN_FILE_PATH` or `AZURE_FEDERATED_TOKEN_FILE` Environment Variable.",
				Optional:            true,
			},
			"oidc_azure_service_connection_id": schema.StringAttribute{
				MarkdownDescription: "The Azure Pipelines Service Connection ID to use for authentication. This can also be sourced from the `ARM_ADO_PIPELINE_SERVICE_CONNECTION_ID`, `ARM_OIDC_AZURE_SERVICE_CONNECTION_ID` or `AZURESUBSCRIPTION_SERVICE_CONNECTION_ID` Environment Variable.",
				Optional:            true,
			},
			"use_oidc": schema.BoolAttribute{
				MarkdownDescription: "Use an OIDC token to authenticate to a service principal. Defaults to `false`. This can also be sourced from the `ARM_USE_OIDC` Environment Variable.",
				Optional:            true,
			},
			"use_cli": schema.BoolAttribute{
				MarkdownDescription: "Use Azure CLI to authenticate. Defaults to `true`. This can also be sourced from the `ARM_USE_CLI` Environment Variable.",
				Optional:            true,
			},
			"use_msi": schema.BoolAttribute{
				MarkdownDescription: "Use an Azure Managed Service Identity. Defaults to `false`. This can also be sourced from the `ARM_USE_MSI` Environment Variable.",
				Optional:            true,
			},
		},
	}
}

func (_ Provider) getModel(ctx context.Context, req provider.ConfigureRequest) (*providerModel, diag.Diagnostics) {
	var model providerModel
	diags := req.Config.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}

	// set the defaults from environment variables
	if model.OrgServiceUrl.IsNull() {
		if v := os.Getenv("AZDO_ORG_SERVICE_URL"); v != "" {
			model.OrgServiceUrl = types.StringValue(v)
		}
	}
	if model.PersonalAccessToken.IsNull() {
		if v := os.Getenv("AZDO_PERSONAL_ACCESS_TOKEN"); v != "" {
			model.PersonalAccessToken = types.StringValue(v)
		}
	}
	if model.ClientID.IsNull() {
		if v := os.Getenv("ARM_CLIENT_ID"); v != "" {
			model.ClientID = types.StringValue(v)
		} else if v := os.Getenv("AZURE_CLIENT_ID"); v != "" {
			model.ClientID = types.StringValue(v)
		}
	}
	if model.ClientIDFilePath.IsNull() {
		if v := os.Getenv("ARM_CLIENT_ID_FILE_PATH"); v != "" {
			model.ClientIDFilePath = types.StringValue(v)
		}
	}
	if model.TenantID.IsNull() {
		if v := os.Getenv("ARM_TENANT_ID"); v != "" {
			model.TenantID = types.StringValue(v)
		}
	}
	if model.AuxiliaryTenantIDs.IsNull() {
		if v := os.Getenv("ARM_AUXILIARY_TENANT_IDS"); v != "" {
			values := make([]attr.Value, 0)
			for v := range strings.SplitSeq(v, ";") {
				values = append(values, types.StringValue(v))
			}
			model.AuxiliaryTenantIDs = types.ListValueMust(types.StringType, values)
		}
	}
	if model.ClientCertificate.IsNull() {
		if v := os.Getenv("ARM_CLIENT_CERTIFICATE"); v != "" {
			model.ClientCertificate = types.StringValue(v)
		}
	}
	if model.ClientCertificatePath.IsNull() {
		if v := os.Getenv("ARM_CLIENT_CERTIFICATE_PATH"); v != "" {
			model.ClientCertificatePath = types.StringValue(v)
		}
	}
	if model.ClientCertificatePassword.IsNull() {
		if v := os.Getenv("ARM_CLIENT_CERTIFICATE_PASSWORD"); v != "" {
			model.ClientCertificatePassword = types.StringValue(v)
		}
	}
	if model.ClientSecret.IsNull() {
		if v := os.Getenv("ARM_CLIENT_SECRET"); v != "" {
			model.ClientSecret = types.StringValue(v)
		}
	}
	if model.ClientSecretFilePath.IsNull() {
		if v := os.Getenv("ARM_CLIENT_SECRET_FILE_PATH"); v != "" {
			model.ClientSecretFilePath = types.StringValue(v)
		}
	}
	if model.OIDCRequestToken.IsNull() {
		if v := os.Getenv("ARM_OIDC_REQUEST_TOKEN"); v != "" {
			model.OIDCRequestToken = types.StringValue(v)
		} else if v := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN"); v != "" {
			model.OIDCRequestToken = types.StringValue(v)
		} else if v := os.Getenv("SYSTEM_ACCESSTOKEN"); v != "" {
			model.OIDCRequestToken = types.StringValue(v)
		}
	}
	if model.OIDCRequestURL.IsNull() {
		if v := os.Getenv("ARM_OIDC_REQUEST_URL"); v != "" {
			model.OIDCRequestURL = types.StringValue(v)
		} else if v := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL"); v != "" {
			model.OIDCRequestURL = types.StringValue(v)
		} else if v := os.Getenv("SYSTEM_OIDCREQUESTURI"); v != "" {
			model.OIDCRequestURL = types.StringValue(v)
		}
	}
	if model.OIDCToken.IsNull() {
		if v := os.Getenv("ARM_OIDC_TOKEN"); v != "" {
			model.OIDCToken = types.StringValue(v)
		}
	}
	if model.OIDCTokenFilePath.IsNull() {
		if v := os.Getenv("ARM_OIDC_TOKEN_FILE_PATH"); v != "" {
			model.OIDCTokenFilePath = types.StringValue(v)
		} else if v := os.Getenv("AZURE_FEDERATED_TOKEN_FILE"); v != "" {
			model.OIDCTokenFilePath = types.StringValue(v)
		}
	}
	if model.OIDCAzureServiceConnectionID.IsNull() {
		if v := os.Getenv("ARM_ADO_PIPELINE_SERVICE_CONNECTION_ID"); v != "" {
			model.OIDCAzureServiceConnectionID = types.StringValue(v)
		} else if v := os.Getenv("ARM_OIDC_AZURE_SERVICE_CONNECTION_ID"); v != "" {
			model.OIDCAzureServiceConnectionID = types.StringValue(v)
		} else if v := os.Getenv("AZURESUBSCRIPTION_SERVICE_CONNECTION_ID"); v != "" {
			model.OIDCAzureServiceConnectionID = types.StringValue(v)
		}
	}
	if model.UseOIDC.IsNull() {
		if v := os.Getenv("ARM_USE_OIDC"); v != "" {
			model.UseOIDC = types.BoolValue(v == "true")
		} else {
			model.UseOIDC = types.BoolValue(false)
		}
	}
	if model.UseCLI.IsNull() {
		if v := os.Getenv("ARM_USE_CLI"); v != "" {
			model.UseCLI = types.BoolValue(v == "true")
		} else {
			model.UseCLI = types.BoolValue(true)
		}
	}
	if model.UseMSI.IsNull() {
		if v := os.Getenv("ARM_USE_MSI"); v != "" {
			model.UseMSI = types.BoolValue(v == "true")
		} else {
			model.UseMSI = types.BoolValue(false)
		}
	}

	return &model, diags
}

func (_ Provider) newAuthProvider(m *providerModel) (azuredevops.AuthProvider, error) {
	// Personal Access Token
	if !m.PersonalAccessToken.IsNull() {
		return azuredevops.NewAuthProviderPAT(m.PersonalAccessToken.ValueString()), nil
	}

	// AAD Authentication
	var auxTenants []string
	for _, tid := range m.AuxiliaryTenantIDs.Elements() {
		auxTenants = append(auxTenants, tid.(basetypes.StringValue).ValueString())
	}

	cred, err := aztfauth.NewCredential(aztfauth.Option{
		Logger:                     log.New(log.Default().Writer(), "[DEBUG] ", log.LstdFlags|log.Lmsgprefix),
		TenantId:                   m.TenantID.ValueString(),
		ClientId:                   m.ClientID.ValueString(),
		ClientIdFile:               m.ClientIDFilePath.ValueString(),
		UseClientSecret:            true,
		ClientSecret:               m.ClientSecret.ValueString(),
		ClientSecretFile:           m.ClientSecretFilePath.ValueString(),
		UseClientCert:              true,
		ClientCertBase64:           m.ClientCertificate.ValueString(),
		ClientCertPfxFile:          m.ClientCertificatePath.ValueString(),
		ClientCertPassword:         []byte(m.ClientCertificatePassword.ValueString()),
		UseOIDCToken:               m.UseOIDC.ValueBool(),
		OIDCToken:                  m.OIDCToken.ValueString(),
		UseOIDCTokenFile:           m.UseOIDC.ValueBool(),
		OIDCTokenFile:              m.OIDCTokenFilePath.ValueString(),
		UseOIDCTokenRequest:        m.UseOIDC.ValueBool(),
		OIDCRequestToken:           m.OIDCRequestToken.ValueString(),
		OIDCRequestURL:             m.OIDCRequestURL.ValueString(),
		ADOServiceConnectionId:     m.OIDCAzureServiceConnectionID.ValueString(),
		UseMSI:                     m.UseMSI.ValueBool(),
		UseAzureCLI:                m.UseCLI.ValueBool(),
		AdditionallyAllowedTenants: auxTenants,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to new credential")
	}

	AzureDevOpsAppDefaultScope := "499b84ac-1321-427f-aa17-267ca6975798/.default"
	ap := azuredevops.NewAuthProviderAAD(cred, policy.TokenRequestOptions{
		Scopes: []string{AzureDevOpsAppDefaultScope},
	})
	return ap, nil
}
