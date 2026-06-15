package provider

import (
	"context"
	"fmt"
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

type NeoliteVMFromSnapshotResource struct {
	client *biznetgio.Client
}

type NeoliteVMFromSnapshotResourceModel struct {
	ID                types.String   `tfsdk:"id"`
	OrderID           types.String   `tfsdk:"order_id"`
	SnapshotID        types.Int64    `tfsdk:"snapshot_id"`
	ProductID         types.Int64    `tfsdk:"product_id"`
	Cycle             types.String   `tfsdk:"cycle"`
	KeypairID         types.Int64    `tfsdk:"keypair_id"`
	Name              types.String   `tfsdk:"name"`
	Description       types.String   `tfsdk:"description"`
	SSHAndConsoleUser types.String   `tfsdk:"ssh_and_console_user"`
	ConsolePassword   types.String   `tfsdk:"console_password"`
	Promocode         types.String   `tfsdk:"promocode"`
	PayWithCreditCard types.Bool     `tfsdk:"pay_with_credit_card"`
	Status            types.String   `tfsdk:"status"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
}

func NewNeoliteVMFromSnapshotResource() resource.Resource {
	return &NeoliteVMFromSnapshotResource{}
}

func (r *NeoliteVMFromSnapshotResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_neolite_vm_from_snapshot"
}

func (r *NeoliteVMFromSnapshotResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "VM NEO Lite baru yang di-restore dari snapshot (endpoint restore-with). Hasilnya VM biasa: delete resource = delete VM. Semua input create-only.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource id = account_id VM hasil restore di BiznetGIO.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"order_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Order id dari response create.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"snapshot_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Account id snapshot sumber, dari `biznetgio_neolite_snapshot`.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"product_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Product id dari data source `biznetgio_neolite_products` atau portal.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"cycle": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Siklus billing: `m` monthly, `a` annual, `q` quarterly, `s` semiannual, `b` biennial, `t` triennial, `p4`, `p5`.",
				Validators:          []validator.String{enumValidator{vals: []string{"m", "a", "q", "s", "b", "t", "p4", "p5"}}},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"keypair_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Id keypair dari `biznetgio_neolite_keypair`.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
