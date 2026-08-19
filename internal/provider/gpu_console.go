package provider

import (
	"context"
	"fmt"

	"github.com/biznetgio/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GpuConsoleDataSource struct {
	client *biznetgio.Client
}

type GpuConsoleDataSourceModel struct {
	AccountID types.Int64  `tfsdk:"account_id"`
	URL       types.String `tfsdk:"url"`
	Raw       types.String `tfsdk:"raw"`
}

func NewGpuConsoleDataSource() datasource.DataSource {
	return &GpuConsoleDataSource{}
}

func (d *GpuConsoleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gpu_console"
}

func (d *GpuConsoleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Akses console NEO GPU (`POST /neo-gpus/accounts/{account_id}/console-access`). Side-effecting: setiap pemanggilan bisa mint session/credential baru, jangan dipakai dalam plan diff.",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Account id GPU yang mau diakses consolenya.",
			},
			"url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "URL console (alias url/console_url/access_url/href).",
			},
			"raw": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "Full JSON response.",
			},
		},
	}
}

func (d *GpuConsoleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*biznetgio.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *biznetgio.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *GpuConsoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GpuConsoleDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, err := d.client.GPU().ConsoleAccess(ctx, data.AccountID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get gpu console: %s", err))
		return
	}
	data.URL = types.StringValue(aliasStr(out, "url", "console_url", "access_url", "href"))
	data.Raw = types.StringValue(rawJSON(out))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
