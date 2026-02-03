package framework

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-azuredevops/internal/meta"
)

var _ datasource.DataSource = datasourceWrapper{}
var _ datasource.DataSourceWithConfigure = datasourceWrapper{}
var _ datasource.DataSourceWithConfigValidators = datasourceWrapper{}
var _ datasource.DataSourceWithValidateConfig = datasourceWrapper{}
var _ DataSourceWithTimeout = datasourceWrapper{}

type datasourceWrapper struct {
	DataSource
}

func WrapDataSource(in DataSource) func() datasource.DataSource {
	return func() datasource.DataSource {
		return datasourceWrapper{DataSource: in}
	}
}

func (d datasourceWrapper) Timeout() DataSourceTimeout {
	if r, ok := d.DataSource.(DataSourceWithTimeout); ok {
		return r.Timeout()
	}
	return DataSourceTimeout{
		Read: 5 * time.Minute,
	}
}

func (d datasourceWrapper) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	d.DataSource.Metadata(ctx, req, resp)
	resp.TypeName = d.DataSource.ResourceType()
}

func (d datasourceWrapper) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	d.DataSource.Schema(ctx, req, resp)
	timeout := d.Timeout()
	resp.Schema.Attributes["timeouts"] = timeouts.AttributesWithOpts(ctx, timeouts.Opts{
		ReadDescription: fmt.Sprintf("(Defaults to %s) Used when reading this data source.", timeout.Read),
	})
}

func (d datasourceWrapper) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer func() {
		d.logDiags(ctx, resp.Diagnostics)
	}()

	ctx = tflog.NewSubsystem(ctx, d.DataSource.ResourceType())
	ctx = tflog.SubsystemSetField(ctx, d.DataSource.ResourceType(), "operation", "Read")

	var timeout timeouts.Value
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("timeouts"), &timeout)...)
	if resp.Diagnostics.HasError() {
		return
	}

	duration, diags := timeout.Read(ctx, d.Timeout().Read)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	tflog.SubsystemInfo(ctx, d.DataSource.ResourceType(), "Start to read the data source")
	d.DataSource.Read(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.SubsystemInfo(ctx, d.DataSource.ResourceType(), "Finish to read the data source")
}

func (d datasourceWrapper) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	defer func() {
		d.logDiags(ctx, resp.Diagnostics)
	}()

	ctx = tflog.NewSubsystem(ctx, d.DataSource.ResourceType())
	ctx = tflog.SubsystemSetField(ctx, d.DataSource.ResourceType(), "operation", "Configure")

	if req.ProviderData == nil {
		return
	}
	d.DataSource.SetMeta(req.ProviderData.(meta.Meta))
}

// ConfigValidators implements datasource.ResourceWithConfigValidators.
func (d datasourceWrapper) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	ctx = tflog.NewSubsystem(ctx, d.DataSource.ResourceType())
	ctx = tflog.SubsystemSetField(ctx, d.DataSource.ResourceType(), "operation", "ConfigValidators")

	if dd, ok := d.DataSource.(datasource.DataSourceWithConfigValidators); ok {
		return dd.ConfigValidators(ctx)
	}
	return nil
}

// ValidateConfig implements datasource.ResourceWithValidateConfig.
func (d datasourceWrapper) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	defer func() {
		d.logDiags(ctx, resp.Diagnostics)
	}()

	ctx = tflog.NewSubsystem(ctx, d.DataSource.ResourceType())
	ctx = tflog.SubsystemSetField(ctx, d.DataSource.ResourceType(), "operation", "ValidateConfig")

	if dd, ok := d.DataSource.(datasource.DataSourceWithValidateConfig); ok {
		dd.ValidateConfig(ctx, req, resp)
		return
	}
}

func (d datasourceWrapper) logDiags(ctx context.Context, diags diag.Diagnostics) {
	for _, warning := range diags.Warnings() {
		tflog.SubsystemWarn(ctx, d.DataSource.ResourceType(), fmt.Sprintf("%s: %s", warning.Summary(), warning.Detail()))
	}
	for _, err := range diags.Errors() {
		tflog.SubsystemError(ctx, d.DataSource.ResourceType(), fmt.Sprintf("%s: %s", err.Summary(), err.Detail()))
	}
}
