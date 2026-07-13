// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"time"

	"github.com/biznetgio/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &BiznetgioProvider{}

type BiznetgioProvider struct {
	version string
}

type BiznetgioProviderModel struct {
	APIKey  types.String `tfsdk:"api_key"`
	BaseURL types.String `tfsdk:"base_url"`
	Timeout types.Int64  `tfsdk:"timeout"`
}

func (p *BiznetgioProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "biznetgio"
	resp.Version = p.version
}

func (p *BiznetgioProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "BiznetGIO API token sent as x-token header. May also be set via BIZNETGIO_API_KEY env var.",
			},
			"base_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "BiznetGIO API base URL. May also be set via BIZNETGIO_BASE_URL env var. Defaults to https://api.portal.biznetgio.com/v1",
			},
			"timeout": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "HTTP client timeout in seconds. Defaults to 30.",
			},
		},
	}
}

func (p *BiznetgioProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	apiKey := os.Getenv("BIZNETGIO_API_KEY")
	baseURL := os.Getenv("BIZNETGIO_BASE_URL")

	var data BiznetgioProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.APIKey.ValueString() != "" {
		apiKey = data.APIKey.ValueString()
	}
	if data.BaseURL.ValueString() != "" {
		baseURL = data.BaseURL.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing API Key Configuration",
			"While configuring the provider, the API key was not found in "+
				"the BIZNETGIO_API_KEY environment variable or provider "+
				"configuration block api_key attribute.",
		)
		// jangan early return, biar error kumpul semua dulu
	}
	if baseURL == "" {
		baseURL = "https://api.portal.biznetgio.com/v1"
	}
	if resp.Diagnostics.HasError() {
		return
	}

	timeout := 30 * time.Second
	if !data.Timeout.IsNull() {
		timeout = time.Duration(data.Timeout.ValueInt64()) * time.Second
	}

	client := biznetgio.New(baseURL, apiKey, timeout)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *BiznetgioProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// metal dulu
		NewBaremetalResource,
		NewBaremetalKeypairResource,
		NewBaremetalAdditionalIPResource,
		NewBaremetalAdditionalIPAssignmentResource,
		NewBaremetalElasticStorageResource,
		// gpu cekidot
		NewGpuInstanceResource,
		NewGpuKeypairResource,
		// neolite gaskeun
		NewNeoliteVMResource,
		NewNeoliteKeypairResource,
		NewNeoliteSnapshotResource,
		NewNeoliteVMFromSnapshotResource,
		NewNeoliteDiskResource,
		// pro gacor
		NewNeoliteProVmResource,
		NewNeoliteProKeypairResource,
		NewNeoliteProSnapshotResource,
		NewNeoliteProDiskResource,
		// object storage jos
		NewObjectStorageResource,
		NewObjectStorageBucketResource,
		NewObjectStorageCredentialResource,
		NewObjectStorageObjectResource,
	}
}

func (p *BiznetgioProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// metal dulu
		NewBaremetalProductsDataSource,
		NewBaremetalRebuildOSListDataSource,
		NewBaremetalOpenVPNDataSource,
		// gpu cekidot
		NewGpuProductsDataSource,
		NewGpuConsoleDataSource,
		NewGpuGraphDataSource,
		// neolite gaskeun
		NewNeoliteProductsDataSource,
		NewNeoliteOsListDataSource,
		NewNeoliteChangePackageOptionsDataSource,
		NewNeoliteStorageUpgradeOptionsDataSource,
		NewNeoliteIPAvailabilityDataSource,
		// pro gacor
		NewNeoliteProProductsDataSource,
		NewNeoliteProOsListDataSource,
		NewNeoliteProChangePackageOptionsDataSource,
		NewNeoliteProStorageUpgradeOptionsDataSource,
		NewNeoliteProIPAvailabilityDataSource,
		// object storage jos
		NewObjectStorageInstancesDataSource,
		NewObjectStorageBucketsDataSource,
		NewObjectStorageCredentialsDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &BiznetgioProvider{
			version: version,
		}
	}
}
// wip 932
