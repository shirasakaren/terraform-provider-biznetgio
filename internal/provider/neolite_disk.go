package provider

import (
	"context"
	"fmt"
	"time"

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
	"github.com/shirasakaren/terraform-provider-biznetgio/internal/biznetgio"
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
				MarkdownDescription: "Account id VM NEO Lite tempat disk dipasang.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"service_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("service-name"),
				MarkdownDescription: "Nama layanan disk. Default `service-name`. Panjang 6-16, hanya huruf/angka/titik/dash.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(6, 16),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9\-\.]*$`), "hanya huruf, angka, titik, atau dash"),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
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
				MarkdownDescription: "Bayar invoice pake kartu kredit saat order. Default true (auto-charge). Set false kalau mau ninggalin invoice unpaid di portal - resource bakal stuck Pending sampai dibayar. Bisa diubah in-place saat upgrade disk.",
			},
			"size": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(60),
				MarkdownDescription: "Ukuran disk (GB). Default 60, minimal 60. Cuma bisa naik, bukan turun.",
				Validators:          []validator.Int64{int64validator.AtLeast(60)},
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Status disk (Active, Pending, Suspended, Terminated).",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"raw": schema.StringAttribute{
				Sensitive:           true,
				Computed:            true,
				MarkdownDescription: "Full JSON response disk terakhir dari API, buat akses field yang belum dimodel.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
		},
	}
}

func (r *NeoliteDiskResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NeoliteDiskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NeoliteDiskResourceModel
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
	out, err := r.client.Neolite().DiskCreate(ctx, biznetgio.NeoliteDiskCreateRequest{
		ProductID:        data.ProductID.ValueInt64(),
		Cycle:            data.Cycle.ValueString(),
		NeoliteAccountID: data.NeoliteAccountID.ValueInt64(),
		ServiceName:      data.ServiceName.ValueString(),
		Promocode:        data.Promocode.ValueString(),
		PayInvoiceWithCC: cc,
		Size:             data.Size.ValueInt64(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create neolite disk: %s", err))
		return
	}
	diskID := aliasInt(out, "account_id", "id")
	if diskID == 0 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Create neolite disk response tidak ada account_id: %s", rawJSON(out)))
		return
	}
	data.ID = types.StringValue(fmt.Sprintf("%d", diskID))
	data.OrderID = types.StringValue(aliasStr(out, "order_id", "orderid"))
	data.Raw = types.StringValue(rawJSON(out))

	tflog.Info(ctx, "neolite disk created, menunggu active", map[string]any{"disk_account_id": diskID})
	done, err := biznetgio.WaitForStatus(ctx, 5*time.Second,
		func(ctx context.Context) (map[string]any, error) {
			return r.client.Neolite().DiskGet(ctx, diskID)
		},
		lowerStatus, []string{"active"}, []string{"suspended", "terminated"})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Neolite disk %d gagal jadi active: %s", diskID, err))
		return
	}
	data.Status = types.StringValue(aliasStr(done, "status", "state"))
	data.Raw = types.StringValue(rawJSON(done))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NeoliteDiskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NeoliteDiskResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, err := r.client.Neolite().DiskGet(ctx, parseAccountID(data.ID.ValueString()))
	if err != nil {
		if biznetgio.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read neolite disk %s: %s", data.ID.ValueString(), err))
		return
	}

	data.Status = types.StringValue(aliasStr(out, "status", "state"))
	if v := aliasStr(out, "service_name", "name", "label"); v != "" {
		data.ServiceName = types.StringValue(v)
	}
	if v := aliasInt(out, "size", "disk_size"); v > 0 {
		data.Size = types.Int64Value(v)
	}
	data.Raw = types.StringValue(rawJSON(out))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NeoliteDiskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state NeoliteDiskResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	diskID := parseAccountID(state.ID.ValueString())
	newSize := plan.Size.ValueInt64()
	oldSize := state.Size.ValueInt64()
	if newSize == oldSize {
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}
	if newSize < oldSize {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Neolite disk cuma bisa di-upgrade: %d -> %d", oldSize, newSize))
		return
	}

	// upgrade pakai additional_size INCREMENT, bukan target absolute.
	_, err := r.client.Neolite().DiskUpgrade(ctx, diskID, biznetgio.NeoliteDiskUpgradeRequest{
		AdditionalSize:   newSize - oldSize,
		PayInvoiceWithCC: ccYesNo(plan.PayWithCreditCard),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to upgrade neolite disk %s: %s", state.ID.ValueString(), err))
		return
	}
	state.Size = plan.Size

	tflog.Info(ctx, "neolite disk upgraded, menunggu active", map[string]any{"disk_account_id": diskID})
	done, err := biznetgio.WaitForStatus(ctx, 5*time.Second,
		func(ctx context.Context) (map[string]any, error) {
			return r.client.Neolite().DiskGet(ctx, diskID)
		},
		lowerStatus, []string{"active"}, []string{"suspended", "terminated"})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Neolite disk %s gagal balik active: %s", state.ID.ValueString(), err))
		return
	}
	state.Status = types.StringValue(aliasStr(done, "status", "state"))
	state.Raw = types.StringValue(rawJSON(done))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NeoliteDiskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NeoliteDiskResourceModel
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

	err := r.client.Neolite().DiskDelete(ctx, parseAccountID(data.ID.ValueString()))
	if err != nil && !biznetgio.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete neolite disk %s: %s", data.ID.ValueString(), err))
		return
	}
}

func (r *NeoliteDiskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
