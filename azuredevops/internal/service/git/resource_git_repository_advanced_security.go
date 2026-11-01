package git

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/sdk/advancedsecurity"
)

// ResourceGitRepositoryAdvancedSecurity schema to manage the Advanced Security enablement of a git repository
func ResourceGitRepositoryAdvancedSecurity() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGitRepositoryAdvancedSecurityCreateUpdate,
		ReadContext:   resourceGitRepositoryAdvancedSecurityRead,
		UpdateContext: resourceGitRepositoryAdvancedSecurityCreateUpdate,
		DeleteContext: resourceGitRepositoryAdvancedSecurityDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceGitRepositoryAdvancedSecurityImport,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"repository_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"secret_protection_enabled": {
				Description: "Secret Protection (secret scanning)",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"block_pushes": {
				Description: "Block pushes that contain secrets",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"code_security_enabled": {
				Description: "Code Security (dependency and code scanning)",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"dependency_scanning_injection_enabled": {
				Description: "Inject dependency scanning into pipelines",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"codeql_enabled": {
				Description: "CodeQL default setup",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"autofix_enabled": {
				Description: "Copilot Autofix",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func getAdvancedSecurityClient(m interface{}) (advancedsecurity.Client, error) {
	clients := m.(*client.AggregatedClient)
	if clients.AdvancedSecurityClient == nil {
		return nil, fmt.Errorf("Advanced Security is not available for organization %s", clients.OrganizationURL)
	}
	return clients.AdvancedSecurityClient, nil
}

func resourceGitRepositoryAdvancedSecurityCreateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	advSecClient, err := getAdvancedSecurityClient(m)
	if err != nil {
		return diag.FromErr(err)
	}
	projectID := d.Get("project_id").(string)
	repositoryID := d.Get("repository_id").(string)

	err = advSecClient.UpdateRepoEnablement(ctx, advancedsecurity.UpdateRepoEnablementArgs{
		Project:    converter.String(projectID),
		Repository: converter.String(repositoryID),
		Settings:   expandRepoEnablement(rawConfigValues(d)),
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("updating Advanced Security of repository %s: %v", repositoryID, err))
	}

	d.SetId(fmt.Sprintf("%s/%s", projectID, repositoryID))
	return resourceGitRepositoryAdvancedSecurityRead(ctx, d, m)
}

func resourceGitRepositoryAdvancedSecurityRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	advSecClient, err := getAdvancedSecurityClient(m)
	if err != nil {
		return diag.FromErr(err)
	}
	projectID := d.Get("project_id").(string)
	repositoryID := d.Get("repository_id").(string)

	settings, err := advSecClient.GetRepoEnablement(ctx, advancedsecurity.GetRepoEnablementArgs{
		Project:              converter.String(projectID),
		Repository:           converter.String(repositoryID),
		IncludeAllProperties: converter.Bool(true),
	})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("reading Advanced Security of repository %s: %v", repositoryID, err))
	}

	if secret := settings.SecretProtectionFeatures; secret != nil {
		d.Set("secret_protection_enabled", converter.ToBool(secret.SecretProtectionEnabled, false))
		d.Set("block_pushes", converter.ToBool(secret.BlockPushes, false))
	}
	if code := settings.CodeSecurityFeatures; code != nil {
		d.Set("code_security_enabled", converter.ToBool(code.CodeSecurityEnabled, false))
		d.Set("dependency_scanning_injection_enabled", converter.ToBool(code.DependencyScanningInjectionEnabled, false))
		d.Set("codeql_enabled", converter.ToBool(code.CodeQLEnabled, false))
		d.Set("autofix_enabled", converter.ToBool(code.AutofixEnabled, false))
	}
	return nil
}

func resourceGitRepositoryAdvancedSecurityDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	advSecClient, err := getAdvancedSecurityClient(m)
	if err != nil {
		return diag.FromErr(err)
	}
	repositoryID := d.Get("repository_id").(string)

	err = advSecClient.UpdateRepoEnablement(ctx, advancedsecurity.UpdateRepoEnablementArgs{
		Project:    converter.String(d.Get("project_id").(string)),
		Repository: converter.String(repositoryID),
		Settings: &advancedsecurity.RepoEnablementSettings{
			SecretProtectionFeatures: &advancedsecurity.SecretProtectionFeatures{SecretProtectionEnabled: converter.Bool(false)},
			CodeSecurityFeatures:     &advancedsecurity.CodeSecurityFeatures{CodeSecurityEnabled: converter.Bool(false)},
		},
	})
	if err != nil && !utils.ResponseWasNotFound(err) {
		return diag.FromErr(fmt.Errorf("disabling Advanced Security of repository %s: %v", repositoryID, err))
	}
	d.SetId("")
	return nil
}

func resourceGitRepositoryAdvancedSecurityImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("unexpected import ID %q, expected <project_id>/<repository_id>", d.Id())
	}
	d.Set("project_id", parts[0])
	d.Set("repository_id", parts[1])
	return []*schema.ResourceData{d}, nil
}

func rawConfigValues(d *schema.ResourceData) map[string]cty.Value {
	rawConfig := d.GetRawConfig()
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		return map[string]cty.Value{}
	}
	return rawConfig.AsValueMap()
}

// expandRepoEnablement sends only the settings present in the configuration; the API leaves the others unchanged.
func expandRepoEnablement(rawConfig map[string]cty.Value) *advancedsecurity.RepoEnablementSettings {
	configured := func(name string) *bool {
		value, ok := rawConfig[name]
		if !ok || value.IsNull() {
			return nil
		}
		return converter.Bool(value.True())
	}

	settings := &advancedsecurity.RepoEnablementSettings{}
	secret := &advancedsecurity.SecretProtectionFeatures{
		SecretProtectionEnabled: configured("secret_protection_enabled"),
		BlockPushes:             configured("block_pushes"),
	}
	if secret.SecretProtectionEnabled != nil || secret.BlockPushes != nil {
		settings.SecretProtectionFeatures = secret
	}
	code := &advancedsecurity.CodeSecurityFeatures{
		CodeSecurityEnabled:                configured("code_security_enabled"),
		DependencyScanningInjectionEnabled: configured("dependency_scanning_injection_enabled"),
		CodeQLEnabled:                      configured("codeql_enabled"),
		AutofixEnabled:                     configured("autofix_enabled"),
	}
	if code.CodeSecurityEnabled != nil || code.DependencyScanningInjectionEnabled != nil || code.CodeQLEnabled != nil || code.AutofixEnabled != nil {
		settings.CodeSecurityFeatures = code
	}
	return settings
}
