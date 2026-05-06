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
				Required:            true,
				MarkdownDescription: "Nama VM hasil restore. Panjang 6-16, hanya huruf/angka/titik/dash.",
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
				MarkdownDescription: "Deskripsi VM.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"ssh_and_console_user": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "User SSH & console yang dipasang saat create.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"console_password": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "Password console saat create. Write-only: ga pernah di-refetch dari API.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"promocode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Kode promo saat order.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"pay_with_credit_card": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Bayar invoice pake kartu kredit saat order. Default true (auto-charge). Set false kalau mau ninggalin invoice unpaid di portal — resource bakal stuck Pending sampai dibayar.",
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Status VM (Active, Pending, Suspended, Terminated).",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Delete: true,
			}),
		},
	}
}

func (r *NeoliteVMFromSnapshotResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*biznetgio.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *biznetgio.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *NeoliteVMFromSnapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NeoliteVMFromSnapshotResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	cc := "yes"
	if !data.PayWithCreditCard.ValueBool() {
		cc = "no"
	}
	billing, err := r.client.Neolite().SnapshotRestoreWith(ctx, data.SnapshotID.ValueInt64(), biznetgio.NeoliteFromSnapshotRequest{
		ProductID:         data.ProductID.ValueInt64(),
		Cycle:             data.Cycle.ValueString(),
		KeypairID:         data.KeypairID.ValueInt64(),
		Name:              data.Name.ValueString(),
		Description:       data.Description.ValueString(),
		SSHAndConsoleUser: data.SSHAndConsoleUser.ValueString(),
		ConsolePassword:   data.ConsolePassword.ValueString(),
		Promocode:         data.Promocode.ValueString(),
		PayInvoiceWithCC:  cc,
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create neolite vm from snapshot: %s", err))
		return
	}
	if billing.AccountID == "" {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Create neolite vm from snapshot response tidak ada account_id: order_id=%s", billing.OrderID))
		return
	}
	data.ID = types.StringValue(billing.AccountID)
	data.OrderID = types.StringValue(billing.OrderID)

	tflog.Info(ctx, "neolite vm from snapshot created, menunggu active", map[string]any{"account_id": billing.AccountID})
	acc, err := biznetgio.WaitForStatus(ctx, 5*time.Second,
		func(ctx context.Context) (biznetgio.AccountResource, error) {
			return r.client.Neolite().AccountGet(ctx, parseAccountID(billing.AccountID))
		},
		accountStatus, []string{"active"}, []string{"suspended", "terminated"})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Neolite vm from snapshot %s gagal jadi active: %s", billing.AccountID, err))
		return
	}
	data.Status = types.StringValue(acc.Status)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NeoliteVMFromSnapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NeoliteVMFromSnapshotResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	acc, err := r.client.Neolite().AccountGet(ctx, parseAccountID(data.ID.ValueString()))
	if err != nil {
		if biznetgio.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
// wip 39
