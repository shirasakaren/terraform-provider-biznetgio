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

// wip 561
