// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/biznetgio/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &GpuProductsDataSource{}

type GpuProductsDataSource struct {
	client *biznetgio.Client
}

type GpuProductsDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Products types.List   `tfsdk:"products"`
}

type GpuProductModel struct {
	ProductID    types.Int64  `tfsdk:"product_id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	CategoryName types.String `tfsdk:"category_name"`
	Raw          types.String `tfsdk:"raw"`
	Flavors      types.List   `tfsdk:"flavors"`
}

type GpuFlavorModel struct {
	FlavorID types.Int64  `tfsdk:"flavor_id"`
	Name     types.String `tfsdk:"name"`
	Raw      types.String `tfsdk:"raw"`
}

func gpuProductAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"product_id":    types.Int64Type,
		"name":          types.StringType,
		"description":   types.StringType,
		"category_name": types.StringType,
		"raw":           types.StringType,
		"flavors":       types.ListType{ElemType: types.ObjectType{AttrTypes: gpuFlavorAttrTypes()}},
	}
}

func gpuFlavorAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"flavor_id": types.Int64Type,
		"name":      types.StringType,
		"raw":       types.StringType,
	}
}

func (d *GpuProductsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "biznetgio_gpu_products"
}

func (d *GpuProductsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "GPU product catalog, including available flavors per product.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Static identifier of the data source.",
			},
			"products": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"product_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Product id, orderable in `biznetgio_gpu_instance`.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Product name.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Product description.",
						},
						"category_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Product category.",
						},
						"raw": schema.StringAttribute{
							Sensitive:           true,
							Computed:            true,
							MarkdownDescription: "Raw JSON of the product, for anything not modeled yet.",
						},
						"flavors": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"flavor_id": schema.Int64Attribute{
										Computed:            true,
										MarkdownDescription: "Flavor id.",
									},
									"name": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "Flavor name.",
									},
									"raw": schema.StringAttribute{
										Sensitive:           true,
										Computed:            true,
										MarkdownDescription: "Raw JSON of the flavor.",
									},
								},
							},
							MarkdownDescription: "Available flavors of the product.",
						},
					},
				},
				MarkdownDescription: "List of GPU products.",
			},
		},
	}
}

func (d *GpuProductsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *GpuProductsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GpuProductsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	products, err := d.client.GPU().ProductsList(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to list gpu products: %s", err))
		return
	}

	models := make([]GpuProductModel, 0, len(products))
	for _, pm := range products {
		pid, _ := gpuInt64(pm, "product_id", "id")
		prod := GpuProductModel{
			ProductID:    types.Int64Value(pid),
			Name:         gpuStringValue(pm, "name", "product_name", "label"),
			Description:  gpuStringValue(pm, "description"),
			CategoryName: gpuStringValue(pm, "category_name", "category"),
			Raw:          types.StringValue(string(redactJSON(pm))),
			Flavors:      types.ListNull(types.ObjectType{AttrTypes: gpuFlavorAttrTypes()}),
		}
		if pid > 0 {
			flavors, ferr := d.client.GPU().ProductFlavors(ctx, pid)
			if ferr != nil {
				tflog.Debug(ctx, "gagal ambil flavors, dilewatin aja", map[string]any{"product_id": pid, "error": ferr.Error()})
			} else {
				fl := make([]GpuFlavorModel, 0, len(flavors))
				for _, fm := range flavors {
					fid, _ := gpuInt64(fm, "flavor_id", "id", "product_id")
					fl = append(fl, GpuFlavorModel{
						FlavorID: types.Int64Value(fid),
						Name:     gpuStringValue(fm, "name", "flavor_name", "label"),
						Raw:      types.StringValue(string(redactJSON(fm))),
					})
				}
				flavorsVal, fdiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: gpuFlavorAttrTypes()}, fl)
				resp.Diagnostics.Append(fdiags...)
				if resp.Diagnostics.HasError() {
					return
				}
				prod.Flavors = flavorsVal
			}
		}
		models = append(models, prod)
	}

	productsVal, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: gpuProductAttrTypes()}, models)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ID = types.StringValue("gpu-products")
	data.Products = productsVal
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func gpuStringValue(m map[string]any, keys ...string) types.String {
	if v, ok := gpuString(m, keys...); ok {
		return types.StringValue(v)
	}
	return types.StringNull()
}

func NewGpuProductsDataSource() datasource.DataSource {
	return &GpuProductsDataSource{}
}
