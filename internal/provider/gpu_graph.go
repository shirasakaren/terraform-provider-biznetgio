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

