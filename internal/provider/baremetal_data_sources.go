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
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"product_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Product id buat `biznetgio_baremetal.product_id`.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Nama product (alias name/product_name/label).",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Deskripsi product.",
						},
						"raw": schema.StringAttribute{
							Sensitive:           true,
							Computed:            true,
							MarkdownDescription: "Full JSON product.",
						},
					},
				},
			},
		},
	}
}

func (d *BaremetalProductsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BaremetalProductsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data BaremetalProductsDataSourceModel
// wip 828
