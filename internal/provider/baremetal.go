package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/biznetgio/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type BaremetalResource struct {
	client *biznetgio.Client
}

type BaremetalResourceModel struct {
	ID                types.String   `tfsdk:"id"`
	AccountID         types.Int64    `tfsdk:"account_id"`
	ProductID         types.Int64    `tfsdk:"product_id"`
	Cycle             types.String   `tfsdk:"cycle"`
	SelectOS          types.String   `tfsdk:"select_os"`
	KeypairID         types.Int64    `tfsdk:"keypair_id"`
	Label             types.String   `tfsdk:"label"`
	PublicIP          types.Int64    `tfsdk:"public_ip"`
	Promocode         types.String   `tfsdk:"promocode"`
	PayWithCreditCard types.Bool     `tfsdk:"pay_with_credit_card"`
	PowerState        types.String   `tfsdk:"power_state"`
	ResetTrigger      types.String   `tfsdk:"reset_trigger"`
	RebuildOS         types.String   `tfsdk:"rebuild_os"`
	Status            types.String   `tfsdk:"status"`
	OrderID           types.String   `tfsdk:"order_id"`
	IPAddress         types.String   `tfsdk:"ip_address"`
	CreatedAt         types.String   `tfsdk:"created_at"`
	Raw               types.String   `tfsdk:"raw"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
}

func NewBaremetalResource() resource.Resource {
	return &BaremetalResource{}
}

func (r *BaremetalResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_baremetal"
}

func (r *BaremetalResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NEO Metal bare-metal server. Order instance, atur label, power state, reset, dan rebuild OS.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource id = account_id dari BiznetGIO.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"account_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Account id baremetal di BiznetGIO.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"product_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Product id dari data source `biznetgio_baremetal_products` atau portal.",
			},
			"cycle": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Siklus billing, misal `m` (monthly) atau `a` (annual).",
			},
			"select_os": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("ubuntu-22"),
				MarkdownDescription: "OS yang dipasang saat create, dari `GET /baremetals/products/{product_id}/oss`. Default `ubuntu-22`.",
