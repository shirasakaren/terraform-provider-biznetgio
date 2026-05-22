package provider

import (
	"context"
	"fmt"

	"github.com/biznetgio/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ---------- biznetgio_neolite_products ----------

type NeoliteProductsDataSource struct {
	client *biznetgio.Client
}

type NeoliteProductsDataSourceModel struct {
	ID       types.String          `tfsdk:"id"`
	Products []NeoliteProductModel `tfsdk:"products"`
}

type NeoliteProductModel struct {
	ProductID    types.Int64                  `tfsdk:"product_id"`
	Name         types.String                 `tfsdk:"name"`
	Description  types.String                 `tfsdk:"description"`
	CategoryID   types.Int64                  `tfsdk:"category_id"`
	CategoryName types.String                 `tfsdk:"category_name"`
	Options      NeoliteProductOptionsModel   `tfsdk:"options"`
	Billing      []NeoliteProductBillingModel `tfsdk:"billing"`
}

type NeoliteProductOptionsModel struct {
	Type           types.String `tfsdk:"type"`
	Cores          types.Int64  `tfsdk:"cores"`
	Memory         types.Int64  `tfsdk:"memory"`
	AllowDowngrade types.Int64  `tfsdk:"allow_downgrade"`
}

type NeoliteProductBillingModel struct {
	Label      types.String                   `tfsdk:"label"`
	Cycle      types.String                   `tfsdk:"cycle"`
	Price      types.Int64                    `tfsdk:"price"`
	Components []NeoliteProductComponentModel `tfsdk:"components"`
}

type NeoliteProductComponentModel struct {
	Label  types.String               `tfsdk:"label"`
	Field  types.String               `tfsdk:"field"`
	Prices []NeoliteProductPriceModel `tfsdk:"prices"`
}

type NeoliteProductPriceModel struct {
	QtyMin types.Int64 `tfsdk:"qty_min"`
	QtyMax types.Int64 `tfsdk:"qty_max"`
	Price  types.Int64 `tfsdk:"price"`
}

func NewNeoliteProductsDataSource() datasource.DataSource {
	return &NeoliteProductsDataSource{}
}

func (d *NeoliteProductsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_neolite_products"
}

func (d *NeoliteProductsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Daftar product NEO Lite (pakai buat isi `product_id` di `biznetgio_neolite_vm`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source id statis.",
			},
			"products": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Daftar product NEO Lite.",
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

func (d *NeoliteProductsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NeoliteProductsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data NeoliteProductsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plans, err := d.client.Neolite().ProductList(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list neolite products: %s", err))
		return
	}

	data.ID = types.StringValue("neolite-products")
	for _, p := range plans {
		model := NeoliteProductModel{
			ProductID:    types.Int64Value(p.ProductID),
			Name:         types.StringValue(p.Name),
			Description:  types.StringValue(p.Description),
			CategoryID:   types.Int64Value(p.CategoryID),
			CategoryName: types.StringValue(p.CategoryName),
			Options: NeoliteProductOptionsModel{
				Type:           types.StringValue(p.Options.Type),
				Cores:          types.Int64Value(p.Options.Cores),
				Memory:         types.Int64Value(p.Options.Memory),
				AllowDowngrade: types.Int64Value(p.Options.AllowDowngrade),
			},
		}
		for _, b := range p.Billing {
			bm := NeoliteProductBillingModel{
				Label: types.StringValue(b.Label),
				Cycle: types.StringValue(b.Cycle),
				Price: types.Int64Value(b.Price),
			}
			for _, c := range b.Components {
				cm := NeoliteProductComponentModel{
					Label: types.StringValue(c.Label),
					Field: types.StringValue(c.Field),
				}
				for _, pr := range c.Prices {
					cm.Prices = append(cm.Prices, NeoliteProductPriceModel{
						QtyMin: types.Int64Value(pr.QtyMin),
						QtyMax: types.Int64Value(pr.QtyMax),
						Price:  types.Int64Value(pr.Price),
					})
				}
				bm.Components = append(bm.Components, cm)
			}
			model.Billing = append(model.Billing, bm)
		}
		data.Products = append(data.Products, model)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ---------- biznetgio_neolite_os_list ----------

type NeoliteOSListDataSource struct {
	client *biznetgio.Client
}

type NeoliteOSListDataSourceModel struct {
	ID        types.String     `tfsdk:"id"`
	ProductID types.Int64      `tfsdk:"product_id"`
	Oss       []NeoliteOSModel `tfsdk:"oss"`
}

type NeoliteOSModel struct {
	VMID   types.Int64  `tfsdk:"vmid"`
	Node   types.String `tfsdk:"node"`
	Name   types.String `tfsdk:"name"`
	MaxMem types.Int64  `tfsdk:"maxmem"`
	MaxCPU types.Int64  `tfsdk:"maxcpu"`
}

func NewNeoliteOsListDataSource() datasource.DataSource {
	return &NeoliteOSListDataSource{}
}

func (d *NeoliteOSListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_neolite_os_list"
}

func (d *NeoliteOSListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Daftar OS yang tersedia untuk product NEO Lite (pakai buat isi `select_os` di `biznetgio_neolite_vm`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source id statis.",
			},
			"product_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Product id NEO Lite.",
			},
			"oss": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Daftar OS.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"vmid":   schema.Int64Attribute{Computed: true, MarkdownDescription: "Id OS template."},
						"node":   schema.StringAttribute{Computed: true, MarkdownDescription: "Node Proxmox."},
						"name":   schema.StringAttribute{Computed: true, MarkdownDescription: "Nama OS — isi buat `select_os`."},
						"maxmem": schema.Int64Attribute{Computed: true, MarkdownDescription: "Memory maksimal (MB)."},
						"maxcpu": schema.Int64Attribute{Computed: true, MarkdownDescription: "CPU maksimal."},
					},
				},
			},
		},
	}
}

func (d *NeoliteOSListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
