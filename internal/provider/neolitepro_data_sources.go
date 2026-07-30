package provider

import (
	"context"
	"fmt"

	"github.com/biznetgio/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ---------- biznetgio_neolite_pro_products ----------

type NeoliteProProductsDataSource struct {
	client *biznetgio.Client
}

type NeoliteProProductsDataSourceModel struct {
	ID       types.String          `tfsdk:"id"`
	Products []NeoliteProductModel `tfsdk:"products"`
}

func NewNeoliteProProductsDataSource() datasource.DataSource {
	return &NeoliteProProductsDataSource{}
}

func (d *NeoliteProProductsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_neolite_pro_products"
}

func (d *NeoliteProProductsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Daftar product NEO Lite Pro (pakai buat isi `product_id` di `biznetgio_neolite_pro_vm`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source id statis.",
			},
			"products": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Daftar product NEO Lite Pro.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"product_id":    schema.Int64Attribute{Computed: true, MarkdownDescription: "Id product."},
						"name":          schema.StringAttribute{Computed: true, MarkdownDescription: "Nama product."},
						"description":   schema.StringAttribute{Computed: true, MarkdownDescription: "Deskripsi product."},
						"category_id":   schema.Int64Attribute{Computed: true, MarkdownDescription: "Id kategori."},
						"category_name": schema.StringAttribute{Computed: true, MarkdownDescription: "Nama kategori."},
						"options": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Opsi product.",
							Attributes: map[string]schema.Attribute{
								"type":            schema.StringAttribute{Computed: true, MarkdownDescription: "Tipe option."},
								"cores":           schema.Int64Attribute{Computed: true, MarkdownDescription: "Jumlah core."},
								"memory":          schema.Int64Attribute{Computed: true, MarkdownDescription: "Memory (MB)."},
								"allow_downgrade": schema.Int64Attribute{Computed: true, MarkdownDescription: "1 kalau downgrade dibolehin."},
							},
						},
						"billing": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Daftar harga billing.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"label": schema.StringAttribute{Computed: true, MarkdownDescription: "Label billing."},
									"cycle": schema.StringAttribute{Computed: true, MarkdownDescription: "Siklus billing."},
									"price": schema.Int64Attribute{Computed: true, MarkdownDescription: "Harga."},
									"components": schema.ListNestedAttribute{
										Computed:            true,
										MarkdownDescription: "Komponen harga.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"label": schema.StringAttribute{Computed: true, MarkdownDescription: "Label komponen."},
												"field": schema.StringAttribute{Computed: true, MarkdownDescription: "Field komponen."},
												"prices": schema.ListNestedAttribute{
													Computed:            true,
													MarkdownDescription: "Tier harga per kuantitas.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"qty_min": schema.Int64Attribute{Computed: true, MarkdownDescription: "Kuantitas minimal."},
															"qty_max": schema.Int64Attribute{Computed: true, MarkdownDescription: "Kuantitas maksimal."},
															"price":   schema.Int64Attribute{Computed: true, MarkdownDescription: "Harga tier."},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// wip 1158
