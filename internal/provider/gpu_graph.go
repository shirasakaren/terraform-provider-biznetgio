package provider

import (
	"context"
	"fmt"

	"github.com/shirasakaren/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GpuGraphDataSource struct {
	client *biznetgio.Client
}

type GpuGraphDataSourceModel struct {
	AccountID types.Int64  `tfsdk:"account_id"`
	Timeframe types.String `tfsdk:"timeframe"`
	Graph     types.String `tfsdk:"graph"`
	Raw       types.String `tfsdk:"raw"`
}

func NewGpuGraphDataSource() datasource.DataSource {
	return &GpuGraphDataSource{}
}

func (d *GpuGraphDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gpu_graph"
}

func (d *GpuGraphDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data monitoring graph NEO GPU (`GET /neo-gpus/accounts/{account_id}/graph-monitor`).",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Account id GPU yang mau diambil graph-nya.",
			},
			"timeframe": schema.StringAttribute{
				Optional:            true,
				Validators:          []validator.String{stringvalidator.OneOf("hour", "day", "week", "month", "year")},
				MarkdownDescription: "Timeframe graph: `hour`, `day`, `week`, `month`, atau `year`. Default `hour`.",
			},
			"graph": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Raw graph payload sebagai string (alias graph/data/payload/content).",
			},
			"raw": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "Full JSON response.",
			},
		},
	}
}

func (d *GpuGraphDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *GpuGraphDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GpuGraphDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeframe := data.Timeframe.ValueString()
	if timeframe == "" {
		timeframe = "hour"
	}
	out, err := d.client.GPU().GraphMonitor(ctx, data.AccountID.ValueInt64(), timeframe)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get gpu graph: %s", err))
		return
	}
	if v := aliasStr(out, "graph", "data", "payload", "content"); v != "" {
		data.Graph = types.StringValue(v)
	} else {
		data.Graph = types.StringValue(rawJSON(out))
	}
	data.Raw = types.StringValue(rawJSON(out))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
