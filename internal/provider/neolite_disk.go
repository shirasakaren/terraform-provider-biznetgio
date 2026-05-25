package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/biznetgio/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"regexp"
)

type NeoliteDiskResource struct {
	client *biznetgio.Client
}

type NeoliteDiskResourceModel struct {
	ID                types.String   `tfsdk:"id"`
	OrderID           types.String   `tfsdk:"order_id"`
	ProductID         types.Int64    `tfsdk:"product_id"`
	Cycle             types.String   `tfsdk:"cycle"`
	NeoliteAccountID  types.Int64    `tfsdk:"neolite_account_id"`
	ServiceName       types.String   `tfsdk:"service_name"`
	Promocode         types.String   `tfsdk:"promocode"`
	PayWithCreditCard types.Bool     `tfsdk:"pay_with_credit_card"`
	Size              types.Int64    `tfsdk:"size"`
	Status            types.String   `tfsdk:"status"`
	Raw               types.String   `tfsdk:"raw"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
}

func NewNeoliteDiskResource() resource.Resource {
	return &NeoliteDiskResource{}
}

func (r *NeoliteDiskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_neolite_disk"
}

func (r *NeoliteDiskResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Disk tambahan NEO Lite (extra disk, punya account id & product catalog sendiri). Upgrade disk pake `additional_size` (increment).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource id = account_id disk di BiznetGIO.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"order_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Order id dari response create.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"product_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Product id disk dari endpoint `GET /neolites/disks/products`.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"cycle": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Siklus billing: `m` monthly, `a` annual, `q` quarterly, `s` semiannual, `b` biennial, `t` triennial, `p4`, `p5`.",
				Validators:          []validator.String{enumValidator{vals: []string{"m", "a", "q", "s", "b", "t", "p4", "p5"}}},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"neolite_account_id": schema.Int64Attribute{
				Required:            true,
