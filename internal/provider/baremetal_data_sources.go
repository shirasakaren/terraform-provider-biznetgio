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

	items, err := d.client.Baremetal().ProductsList(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list baremetal products: %s", err))
		return
	}
	for _, it := range items {
		data.Products = append(data.Products, BaremetalProductModel{
			ProductID:   types.Int64Value(aliasInt(it, "product_id", "id")),
			Name:        types.StringValue(aliasStr(it, "name", "product_name", "label")),
			Description: types.StringValue(aliasStr(it, "description")),
			Raw:         types.StringValue(rawJSON(it)),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ---------- rebuild os list ----------

type BaremetalRebuildOSListDataSource struct {
	client *biznetgio.Client
}

type BaremetalRebuildOSListDataSourceModel struct {
	AccountID types.Int64  `tfsdk:"account_id"`
	OSS       types.List   `tfsdk:"oss"`
	Raw       types.String `tfsdk:"raw"`
}

func NewBaremetalRebuildOSListDataSource() datasource.DataSource {
	return &BaremetalRebuildOSListDataSource{}
}

func (d *BaremetalRebuildOSListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_baremetal_rebuild_os_list"
}

func (d *BaremetalRebuildOSListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List OS yang valid buat rebuild (`GET /baremetals/{account_id}/rebuild/oss`) — lebih sempit dari OS list create karena ngikutin layout disk instance.",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Account id baremetal yang mau di-rebuild.",
			},
			"oss": schema.ListAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List identifier OS (misal `centos7-base`).",
			},
			"raw": schema.StringAttribute{
				Sensitive:           true,
				Computed:            true,
				MarkdownDescription: "Full JSON response.",
			},
		},
	}
}

func (d *BaremetalRebuildOSListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BaremetalRebuildOSListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data BaremetalRebuildOSListDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	items, err := d.client.Baremetal().RebuildOSList(ctx, data.AccountID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list rebuild oss: %s", err))
		return
	}

	var oss []string
	for _, it := range items {
		// item bisa object {name: "..."} atau map lain — ambil defensif
		if v := aliasStr(it, "name", "os", "label", "value"); v != "" {
			oss = append(oss, v)
		}
	}
	listValue, diags := types.ListValueFrom(ctx, types.StringType, oss)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.OSS = listValue
	data.Raw = types.StringValue(rawJSON(map[string]any{"oss": items}))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ---------- openvpn ----------

type BaremetalOpenVPNDataSource struct {
	client *biznetgio.Client
}

type BaremetalOpenVPNDataSourceModel struct {
	Config types.String `tfsdk:"config"`
	Raw    types.String `tfsdk:"raw"`
}

func NewBaremetalOpenVPNDataSource() datasource.DataSource {
	return &BaremetalOpenVPNDataSource{}
}

func (d *BaremetalOpenVPNDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_baremetal_openvpn"
}

func (d *BaremetalOpenVPNDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Konfigurasi OpenVPN akses OOB ke bare-metal (`GET /baremetals/openvpn`, scope seluruh customer).",
		Attributes: map[string]schema.Attribute{
