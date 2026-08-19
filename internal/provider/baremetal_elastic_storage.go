package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/shirasakaren/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BaremetalElasticStorageResource struct {
	client *biznetgio.Client
}

type BaremetalElasticStorageResourceModel struct {
	ID                types.String   `tfsdk:"id"`
	AccountID         types.Int64    `tfsdk:"account_id"`
	ProductID         types.Int64    `tfsdk:"product_id"`
	Cycle             types.String   `tfsdk:"cycle"`
	StorageName       types.String   `tfsdk:"storage_name"`
	MetalAccountID    types.Int64    `tfsdk:"metal_account_id"`
	Size              types.Int64    `tfsdk:"size"`
	Promocode         types.String   `tfsdk:"promocode"`
	PayWithCreditCard types.Bool     `tfsdk:"pay_with_credit_card"`
	Status            types.String   `tfsdk:"status"`
	CreatedAt         types.String   `tfsdk:"created_at"`
	Raw               types.String   `tfsdk:"raw"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
}

func NewBaremetalElasticStorageResource() resource.Resource {
	return &BaremetalElasticStorageResource{}
}

func (r *BaremetalElasticStorageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_baremetal_elastic_storage"
}

func (r *BaremetalElasticStorageResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NEO Elastic Storage — block storage yang ke-bind permanen ke satu baremetal saat create (gak bisa pindah server).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource id = account_id storage.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"account_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Account id elastic storage di BiznetGIO.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"product_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Product id dari `GET /baremetal-neo-elastic-storages/products`. Kalau berubah, provider panggil change-package (`POST .../{account_id}`).",
			},
			"cycle": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Siklus billing, misal `m` (monthly) atau `a` (annual).",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"storage_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Nama storage.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"metal_account_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Account id baremetal target. Cuma bisa di-set pas create (no re-attach endpoint).",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"size": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(100),
				MarkdownDescription: "Ukuran storage dalam GB. Default 100. Kalau berubah, provider panggil upgrade (`PUT .../{account_id}`) — grow-only, set lebih kecil bakal ditolak API.",
			},
			"promocode": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Kode promo.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"pay_with_credit_card": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Bayar invoice pake kartu kredit saat order/upgrade. Default true (auto-charge). Set false kalau mau ninggalin invoice unpaid di portal — resource bakal stuck Pending sampai dibayar.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Status storage (misal Active, Pending, Suspended, Terminated).",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Tanggal dibuat (alias created_at/date_created).",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"raw": schema.StringAttribute{
				Sensitive:           true,
				Computed:            true,
				MarkdownDescription: "Full JSON response terakhir dari API.",
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

func (r *BaremetalElasticStorageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BaremetalElasticStorageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data BaremetalElasticStorageResourceModel
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
	out, err := r.client.BaremetalElasticStorage().Create(ctx, biznetgio.ElasticStorageCreateRequest{
		ProductID:        data.ProductID.ValueInt64(),
		Cycle:            data.Cycle.ValueString(),
		StorageName:      data.StorageName.ValueString(),
		MetalAccountID:   data.MetalAccountID.ValueInt64(),
		Size:             data.Size.ValueInt64(),
		Promocode:        data.Promocode.ValueString(),
		PayInvoiceWithCC: cc,
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create elastic storage: %s", err))
		return
	}
	accountID := aliasInt(out, "account_id", "id")
	if accountID == 0 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Create elastic storage response tidak ada account_id: %s", rawJSON(out)))
		return
	}
	data.ID = types.StringValue(strconv.FormatInt(accountID, 10))
	data.AccountID = types.Int64Value(accountID)
	data.Raw = types.StringValue(rawJSON(out))

	done, err := biznetgio.WaitForStatus(ctx, 5*time.Second,
		func(ctx context.Context) (map[string]any, error) {
			return r.client.BaremetalElasticStorage().Get(ctx, accountID)
		},
		lowerStatus, []string{"active"}, []string{"terminated", "error", "failed"})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Elastic storage %d gagal jadi active: %s", accountID, err))
		return
	}
	data.Status = types.StringValue(aliasStr(done, "status", "state"))
	data.Raw = types.StringValue(rawJSON(done))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BaremetalElasticStorageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data BaremetalElasticStorageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	accountID, err := strconv.ParseInt(data.ID.ValueString(), 10, 64)
	if err != nil || accountID == 0 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Invalid elastic storage id: %q", data.ID.ValueString()))
		return
	}

	out, err := r.client.BaremetalElasticStorage().Get(ctx, accountID)
	if err != nil {
		if biznetgio.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read elastic storage %d: %s", accountID, err))
		return
	}

	data.AccountID = types.Int64Value(accountID)
	data.Status = types.StringValue(aliasStr(out, "status", "state"))
	if v := aliasInt(out, "size"); v != 0 {
		data.Size = types.Int64Value(v)
	}
	if v := aliasStr(out, "created_at", "date_created"); v != "" {
		data.CreatedAt = types.StringValue(v)
	}
	data.Raw = types.StringValue(rawJSON(out))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BaremetalElasticStorageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state BaremetalElasticStorageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	accountID, err := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if err != nil || accountID == 0 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Invalid elastic storage id: %q", state.ID.ValueString()))
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	cc := "yes"
	if !plan.PayWithCreditCard.ValueBool() {
		cc = "no"
	}

	if !plan.Size.Equal(state.Size) {
		if _, err := r.client.BaremetalElasticStorage().Upgrade(ctx, accountID, biznetgio.UpgradeElasticStorageRequest{
			Size:             plan.Size.ValueInt64(),
			PayInvoiceWithCC: cc,
		}); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to upgrade elastic storage size: %s", err))
			return
		}
	}
	if !plan.ProductID.Equal(state.ProductID) {
		if _, err := r.client.BaremetalElasticStorage().ChangePackage(ctx, accountID, biznetgio.ChangeElasticStoragePackageRequest{
			NewProductID:     plan.ProductID.ValueInt64(),
			PayInvoiceWithCC: cc,
		}); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to change elastic storage package: %s", err))
			return
		}
	}

	if _, err := biznetgio.WaitForStatus(ctx, 5*time.Second,
		func(ctx context.Context) (map[string]any, error) {
			return r.client.BaremetalElasticStorage().Get(ctx, accountID)
		},
		lowerStatus, []string{"active"}, []string{"terminated", "error", "failed"}); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Elastic storage %d gagal active setelah update: %s", accountID, err))
		return
	}

	plan.ID = state.ID
	plan.AccountID = state.AccountID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BaremetalElasticStorageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data BaremetalElasticStorageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	accountID, err := strconv.ParseInt(data.ID.ValueString(), 10, 64)
	if err != nil || accountID == 0 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Invalid elastic storage id: %q", data.ID.ValueString()))
		return
	}
	if err := r.client.BaremetalElasticStorage().Delete(ctx, accountID); err != nil && !biznetgio.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete elastic storage %d: %s", accountID, err))
	}
}

func (r *BaremetalElasticStorageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
