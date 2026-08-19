package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/shirasakaren/terraform-provider-biznetgio/internal/biznetgio"
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
		MarkdownDescription: "Snapshot akun NEO Lite (snapshot bayar per-buatan, punya account id sendiri). Snapshot ga bisa di-update - ganti input = recreate.",
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
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"cycle": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Siklus billing: `m` monthly, `a` annual, `q` quarterly, `s` semiannual, `b` biennial, `t` triennial, `p4`, `p5`.",
				Validators:          []validator.String{enumValidator{vals: []string{"m", "a", "q", "s", "b", "t", "p4", "p5"}}},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"pay_with_credit_card": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Bayar invoice pake kartu kredit saat order. Default true (auto-charge). Set false kalau mau ninggalin invoice unpaid di portal - resource bakal stuck Pending sampai dibayar.",
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"promocode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Kode promo saat order.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Status snapshot (Active, Pending, Suspended, Terminated).",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Delete: true,
			}),
		},
	}
}

func (r *NeoliteSnapshotResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NeoliteSnapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NeoliteSnapshotResourceModel
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
	billing, err := r.client.Neolite().SnapshotCreate(ctx, data.NeoliteAccountID.ValueInt64(), biznetgio.NeoliteSnapshotRequest{
		Cycle:            data.Cycle.ValueString(),
		Name:             data.Name.ValueString(),
		Description:      data.Description.ValueString(),
		Promocode:        data.Promocode.ValueString(),
		PayInvoiceWithCC: cc,
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create neolite snapshot: %s", err))
		return
	}
	if billing.AccountID == "" {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Create neolite snapshot response tidak ada account_id: order_id=%s", billing.OrderID))
		return
	}
	data.ID = types.StringValue(billing.AccountID)
	data.OrderID = types.StringValue(billing.OrderID)

	tflog.Info(ctx, "neolite snapshot created, menunggu active", map[string]any{"snapshot_account_id": billing.AccountID})
	acc, err := biznetgio.WaitForStatus(ctx, 5*time.Second,
		func(ctx context.Context) (biznetgio.SnapshotAccountResource, error) {
			return r.client.Neolite().AccountSnapshotGet(ctx, parseAccountID(billing.AccountID))
		},
		func(a biznetgio.SnapshotAccountResource) string { return strings.ToLower(a.Status) },
		[]string{"active"}, []string{"suspended", "terminated"})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Neolite snapshot %s gagal jadi active: %s", billing.AccountID, err))
		return
	}
	data.Status = types.StringValue(acc.Status)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NeoliteSnapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NeoliteSnapshotResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	list, err := r.client.Neolite().AccountSnapshotList(ctx, "")
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list neolite snapshots: %s", err))
		return
	}
	var found *biznetgio.SnapshotAccountResource
	for i := range list {
		if list[i].AccountID == data.ID.ValueString() {
			found = &list[i]
			break
		}
	}
	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	data.Status = types.StringValue(found.Status)
	if v := found.ExtraDetails.Name; v != "" {
		data.Name = types.StringValue(v)
	}
	if v := found.ExtraDetails.Description; v != "" {
		data.Description = types.StringValue(v)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NeoliteSnapshotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// semua input RequiresReplace - update ga akan ke-schedule sama framework.
	resp.Diagnostics.AddError("Unsupported Update", "neolite_snapshot tidak support update; ganti input buat recreate.")
}

func (r *NeoliteSnapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NeoliteSnapshotResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, 10*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	err := r.client.Neolite().SnapshotDelete(ctx, parseAccountID(data.ID.ValueString()))
	if err != nil && !biznetgio.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete neolite snapshot %s: %s", data.ID.ValueString(), err))
		return
	}
}

func (r *NeoliteSnapshotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
