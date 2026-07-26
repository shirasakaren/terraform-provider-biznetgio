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
