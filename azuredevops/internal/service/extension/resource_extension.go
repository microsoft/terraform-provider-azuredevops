package extension

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/extensionmanagement"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
)

func ResourceExtension() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceExtensionCreate,
		ReadContext:   resourceExtensionRead,
		UpdateContext: resourceExtensionUpdate,
		DeleteContext: resourceExtensionDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				id := strings.Split(d.Id(), "/")
				if len(id) != 2 {
					return nil, fmt.Errorf("unexpected ID format, expected <publisherName>/<extensionName>")
				}
				d.Set("publisher_id", id[0])
				d.Set("extension_id", id[1])
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"extension_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"publisher_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validation.StringIsNotWhiteSpace,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"extension_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publisher_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scope": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceExtensionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	publisherId := d.Get("publisher_id").(string)
	extensionId := d.Get("extension_id").(string)
	extension, err := clients.ExtensionManagementClient.InstallExtensionByName(ctx,
		extensionmanagement.InstallExtensionByNameArgs{
			PublisherName: &publisherId,
			ExtensionName: &extensionId,
			Version:       converter.String(d.Get("version").(string)),
		})

	if err != nil {
		return diag.Errorf(" Installing extension for Publisher: %s, Name: %s. Error: %v", publisherId, extensionId, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", *extension.PublisherId, *extension.ExtensionId))
	return resourceExtensionRead(ctx, d, m)
}

func resourceExtensionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	publisherId := d.Get("publisher_id").(string)
	extensionId := d.Get("extension_id").(string)
	extension, err := clients.ExtensionManagementClient.GetInstalledExtensionByName(ctx, extensionmanagement.GetInstalledExtensionByNameArgs{
		PublisherName: &publisherId,
		ExtensionName: &extensionId,
	})

	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf(" Get extension for Publisher: %s, name: %s. Error: %v", publisherId, extensionId, err)
	}

	if extension != nil {
		d.Set("extension_id", extension.ExtensionId)
		d.Set("publisher_id", extension.PublisherId)
		d.Set("version", extension.Version)
		d.Set("scope", extension.Scopes)
		d.Set("publisher_name", extension.PublisherName)
		d.Set("extension_name", extension.ExtensionName)
		if extension.InstallState != nil && extension.InstallState.Flags != nil {
			d.Set("disabled", false)

			flagsStr := string(*extension.InstallState.Flags)
			flags := strings.Split(flagsStr, ",")
			for _, flag := range flags {
				if flag == string(extensionmanagement.ExtensionStateFlagsValues.Disabled) {
					d.Set("disabled", true)
					break
				}
			}
		}
	}
	return nil
}

func resourceExtensionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	disableRaw := d.GetRawConfig().AsValueMap()["disabled"]
	if !disableRaw.IsNull() {
		publisherId := d.Get("publisher_id").(string)
		extensionId := d.Get("extension_id").(string)
		updateArgs := extensionmanagement.UpdateInstalledExtensionArgs{
			Extension: &extensionmanagement.InstalledExtension{
				PublisherId: &publisherId,
				ExtensionId: &extensionId,
				Version:     converter.String(d.Get("version").(string)),
				InstallState: &extensionmanagement.InstalledExtensionState{
					Flags: converter.ToPtr(extensionmanagement.ExtensionStateFlagsValues.None),
				},
			},
		}

		if !disableRaw.False() {
			updateArgs.Extension.InstallState = &extensionmanagement.InstalledExtensionState{
				Flags: converter.ToPtr(extensionmanagement.ExtensionStateFlagsValues.Disabled),
			}
		}

		_, err := clients.ExtensionManagementClient.UpdateInstalledExtension(ctx, updateArgs)
		if err != nil {
			return diag.Errorf(" Update extension for Publisher: %s, Name: %s. Error: %v", publisherId, extensionId, err)
		}
	}
	return resourceExtensionRead(ctx, d, m)
}

func resourceExtensionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	publisher := d.Get("publisher_id").(string)
	name := d.Get("extension_id").(string)

	err := clients.ExtensionManagementClient.UninstallExtensionByName(ctx, extensionmanagement.UninstallExtensionByNameArgs{
		PublisherName: &publisher,
		ExtensionName: &name,
	})

	if err != nil {
		return diag.Errorf(" Uninstalling extension for Publisher: %s, name: %s. Error: %v", publisher, name, err)
	}
	return nil
}
