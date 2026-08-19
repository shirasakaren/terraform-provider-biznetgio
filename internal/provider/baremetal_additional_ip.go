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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BaremetalAdditionalIPResource struct {
	client *biznetgio.Client
}

type BaremetalAdditionalIPResourceModel struct {
	ID                types.String   `tfsdk:"id"`
	AccountID         types.Int64    `tfsdk:"account_id"`
	ProductID         types.Int64    `tfsdk:"product_id"`
	Cycle             types.String   `tfsdk:"cycle"`
	Region            types.String   `tfsdk:"region"`
	Promocode         types.String   `tfsdk:"promocode"`
	PayWithCreditCard types.Bool     `tfsdk:"pay_with_credit_card"`
	Status            types.String   `tfsdk:"status"`
	IPAddress         types.String   `tfsdk:"ip_address"`
	CreatedAt         types.String   `tfsdk:"created_at"`
	Raw               types.String   `tfsdk:"raw"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
}

func NewBaremetalAdditionalIPResource() resource.Resource {
	return &BaremetalAdditionalIPResource{}
}

func (r *BaremetalAdditionalIPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_baremetal_additional_ip"
}

func (r *BaremetalAdditionalIPResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Additional (floating) IP untuk NEO Metal. Diorder independen dari server, lalu di-attach pake `biznetgio_baremetal_additional_ip_assignment`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource id = account_id IP.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"account_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Account id additional IP di BiznetGIO.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"product_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Product id dari `GET /baremetal-additional-ips/products`.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"cycle": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Siklus billing, misal `m` (monthly) atau `a` (annual).",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"region": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("wjv-1"),
				MarkdownDescription: "Region datacenter, list valid dari `GET /baremetal-additional-ips/regions`. Default `wjv-1`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
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
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Status IP (misal Active, Pending, Suspended, Terminated).",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"ip_address": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Alamat IP yang di-assign kalau ada di response.",
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
				Delete: true,
			}),
		},
	}
}

func (r *BaremetalAdditionalIPResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BaremetalAdditionalIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data BaremetalAdditionalIPResourceModel
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
	out, err := r.client.BaremetalAdditionalIP().Create(ctx, biznetgio.AdditionalIPCreateRequest{
		ProductID:        data.ProductID.ValueInt64(),
		Cycle:            data.Cycle.ValueString(),
		Region:           data.Region.ValueString(),
		Promocode:        data.Promocode.ValueString(),
		PayInvoiceWithCC: cc,
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create additional ip: %s", err))
		return
	}
	accountID := aliasInt(out, "account_id", "id")
	if accountID == 0 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Create additional ip response tidak ada account_id: %s", rawJSON(out)))
		return
	}
	data.ID = types.StringValue(strconv.FormatInt(accountID, 10))
	data.AccountID = types.Int64Value(accountID)
	data.Raw = types.StringValue(rawJSON(out))

	done, err := biznetgio.WaitForStatus(ctx, 5*time.Second,
		func(ctx context.Context) (map[string]any, error) {
			return r.client.BaremetalAdditionalIP().Get(ctx, accountID)
		},
		lowerStatus, []string{"active"}, []string{"terminated", "error", "failed"})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Additional ip %d gagal jadi active: %s", accountID, err))
		return
	}
	data.Status = types.StringValue(aliasStr(done, "status", "state"))
	data.Raw = types.StringValue(rawJSON(done))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update gak ada endpointnya — semua input RequiresReplace, method ini cuma formalitas interface.
func (r *BaremetalAdditionalIPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data BaremetalAdditionalIPResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BaremetalAdditionalIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data BaremetalAdditionalIPResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	accountID, err := strconv.ParseInt(data.ID.ValueString(), 10, 64)
	if err != nil || accountID == 0 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Invalid additional ip id: %q", data.ID.ValueString()))
		return
	}

	out, err := r.client.BaremetalAdditionalIP().Get(ctx, accountID)
	if err != nil {
		if biznetgio.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read additional ip %d: %s", accountID, err))
		return
	}

	data.AccountID = types.Int64Value(accountID)
	data.Status = types.StringValue(aliasStr(out, "status", "state"))
	if v := aliasStr(out, "ip", "public_ip", "ip_address", "ipv4"); v != "" {
		data.IPAddress = types.StringValue(v)
	}
	if v := aliasStr(out, "created_at", "date_created"); v != "" {
		data.CreatedAt = types.StringValue(v)
	}
	data.Raw = types.StringValue(rawJSON(out))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BaremetalAdditionalIPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data BaremetalAdditionalIPResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	accountID, err := strconv.ParseInt(data.ID.ValueString(), 10, 64)
	if err != nil || accountID == 0 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Invalid additional ip id: %q", data.ID.ValueString()))
		return
	}
	if err := r.client.BaremetalAdditionalIP().Delete(ctx, accountID); err != nil && !biznetgio.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete additional ip %d: %s", accountID, err))
	}
}

func (r *BaremetalAdditionalIPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
