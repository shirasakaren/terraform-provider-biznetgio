package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/biznetgio/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"regexp"
)

type NeoliteSnapshotResource struct {
	client *biznetgio.Client
}

type NeoliteSnapshotResourceModel struct {
	ID                types.String   `tfsdk:"id"`
	OrderID           types.String   `tfsdk:"order_id"`
	NeoliteAccountID  types.Int64    `tfsdk:"neolite_account_id"`
	Name              types.String   `tfsdk:"name"`
	Description       types.String   `tfsdk:"description"`
	Cycle             types.String   `tfsdk:"cycle"`
	PayWithCreditCard types.Bool     `tfsdk:"pay_with_credit_card"`
	Promocode         types.String   `tfsdk:"promocode"`
	Status            types.String   `tfsdk:"status"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
}

func NewNeoliteSnapshotResource() resource.Resource {
	return &NeoliteSnapshotResource{}
}

func (r *NeoliteSnapshotResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_neolite_snapshot"
}

func (r *NeoliteSnapshotResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Snapshot akun NEO Lite (snapshot bayar per-buatan, punya account id sendiri). Snapshot ga bisa di-update — ganti input = recreate.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource id = account_id snapshot di BiznetGIO.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"order_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Order id dari response create.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"neolite_account_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Account id VM NEO Lite yang di-snapshot.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("snapshot-name"),
				MarkdownDescription: "Nama snapshot. Default `snapshot-name`. Panjang 6-16, hanya huruf/angka/titik/dash.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(6, 16),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9\-\.]*$`), "hanya huruf, angka, titik, atau dash"),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Deskripsi snapshot.",
