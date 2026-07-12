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

type NeoliteProDiskResource struct {
	client *biznetgio.Client
}

type NeoliteProDiskResourceModel struct {
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

func NewNeoliteProDiskResource() resource.Resource {
	return &NeoliteProDiskResource{}
}

func (r *NeoliteProDiskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_neolite_pro_disk"
}

func (r *NeoliteProDiskResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Disk tambahan NEO Lite Pro (extra disk, punya account id & product catalog sendiri). Upgrade disk pake `additional_size` (increment).",
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
				MarkdownDescription: "Product id disk dari endpoint `GET /neolite-pros/disks/products`.",
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
				MarkdownDescription: "Account id VM NEO Lite Pro tempat disk dipasang.",
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
				MarkdownDescription: "Bayar invoice pake kartu kredit saat order. Default true (auto-charge). Set false kalau mau ninggalin invoice unpaid di portal — resource bakal stuck Pending sampai dibayar. Bisa diubah in-place saat upgrade disk.",
			},
			"size": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(30),
				MarkdownDescription: "Ukuran disk (GB). Default 30, minimal 30 (minimal product disk pro). Cuma bisa naik, bukan turun.",
				Validators:          []validator.Int64{int64validator.AtLeast(30)},
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

