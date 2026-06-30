package provider

import (
	"context"
	"fmt"

	"github.com/biznetgio/terraform-provider-biznetgio/internal/biznetgio"
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
