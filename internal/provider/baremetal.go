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
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"keypair_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Id keypair dari `biznetgio_baremetal_keypair`. Keypair pool baremetal terpisah dari neolite/gpu.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"label": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Nama tampilan server. Satu-satunya field yang bisa diupdate via `PUT /baremetals/{account_id}`.",
			},
			"public_ip": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
				MarkdownDescription: "Jumlah public ip yang diminta (1 = dengan public ip). Enum `Public_IP_options`.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
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
				MarkdownDescription: "Bayar invoice pake kartu kredit saat order. Default true (auto-charge). Set false kalau mau ninggalin invoice unpaid di portal — resource bakal stuck Pending sampai dibayar.",
			},
			"power_state": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Power state server: `on` atau `off`. Update hanya mengirim `PUT .../state/{state}` kalau nilainya berubah.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"reset_trigger": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Trigger one-shot reset/reboot: ganti nilainya buat re-trigger. Reset cuma aksi, bukan state stabil.",
			},
			"rebuild_os": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Kalau berubah, instance di-rebuild (wipe OS) pake `PUT /baremetals/{account_id}/rebuild`. List OS valid ada di data source `biznetgio_baremetal_rebuild_os_list`.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Status terakhir dari API (misal Active, Pending, Suspended, Terminated).",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"order_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Order id dari response create.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"ip_address": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Public ip address server kalau ada di response.",
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
				MarkdownDescription: "Full JSON response terakhir dari API, buat akses field yang belum dimodel.",
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

func (r *BaremetalResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BaremetalResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data BaremetalResourceModel
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
	out, err := r.client.Baremetal().Create(ctx, biznetgio.BaremetalCreateRequest{
		ProductID:        data.ProductID.ValueInt64(),
		Cycle:            data.Cycle.ValueString(),
		SelectOS:         data.SelectOS.ValueString(),
		KeypairID:        data.KeypairID.ValueInt64(),
		Label:            data.Label.ValueString(),
		PublicIP:         data.PublicIP.ValueInt64(),
		Promocode:        data.Promocode.ValueString(),
		PayInvoiceWithCC: cc,
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create baremetal: %s", err))
		return
	}
