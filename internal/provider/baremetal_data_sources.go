package provider

import (
	"context"
	"fmt"

	"github.com/biznetgio/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ---------- products ----------

type BaremetalProductsDataSource struct {
	client *biznetgio.Client
}

type BaremetalProductsDataSourceModel struct {
	Products []BaremetalProductModel `tfsdk:"products"`
}

type BaremetalProductModel struct {
	ProductID   types.Int64  `tfsdk:"product_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Raw         types.String `tfsdk:"raw"`
}

func NewBaremetalProductsDataSource() datasource.DataSource {
	return &BaremetalProductsDataSource{}
}

func (d *BaremetalProductsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_baremetal_products"
}

func (d *BaremetalProductsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List product NEO Metal dari `GET /baremetals/products`.",
		Attributes: map[string]schema.Attribute{
			"products": schema.ListNestedAttribute{
				Computed: true,
