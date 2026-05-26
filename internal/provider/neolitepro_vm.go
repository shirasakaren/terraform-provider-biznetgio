package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/biznetgio/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"regexp"
)

type NeoliteProVmResource struct {
	client *biznetgio.Client
}

type NeoliteProVmResourceModel struct {
	ID                types.String            `tfsdk:"id"`
	OrderID           types.String            `tfsdk:"order_id"`
	SSHAndConsoleUser types.String            `tfsdk:"ssh_and_console_user"`
	ConsolePassword   types.String            `tfsdk:"console_password"`
	VMName            types.String            `tfsdk:"vm_name"`
	Description       types.String            `tfsdk:"description"`
	ProductID         types.Int64             `tfsdk:"product_id"`
	SelectOS          types.String            `tfsdk:"select_os"`
	KeypairID         types.Int64             `tfsdk:"keypair_id"`
	Cycle             types.String            `tfsdk:"cycle"`
	PayWithCreditCard types.Bool              `tfsdk:"pay_with_credit_card"`
	Promocode         types.String            `tfsdk:"promocode"`
	PowerState        types.String            `tfsdk:"power_state"`
	RebuildOS         types.String            `tfsdk:"rebuild_os"`
	DiskSize          types.Int64             `tfsdk:"disk_size"`
	Status            types.String            `tfsdk:"status"`
	Uptime            types.Int64             `tfsdk:"uptime"`
	MaxDisk           types.Int64             `tfsdk:"maxdisk"`
	MaxMem            types.Int64             `tfsdk:"maxmem"`
	Mem               types.Int64             `tfsdk:"mem"`
	CPUs              types.Int64             `tfsdk:"cpus"`
	CIUser            types.String            `tfsdk:"ciuser"`
	CIPassword        types.String            `tfsdk:"cipassword"`
	OSName            types.String            `tfsdk:"osname"`
	Region            types.String            `tfsdk:"region"`
	RegionLabel       types.String            `tfsdk:"region_label"`
	NextDue           types.String            `tfsdk:"next_due"`
	RecurringAmount   types.Int64             `tfsdk:"recurring_amount"`
	BillingCycle      types.String            `tfsdk:"billingcycle"`
	ProductName       types.String            `tfsdk:"product_name"`
	LastInvoice       NeoliteLastInvoiceModel `tfsdk:"last_invoice"`
	Raw               types.String            `tfsdk:"raw"`
	Timeouts          timeouts.Value          `tfsdk:"timeouts"`
}

func NewNeoliteProVmResource() resource.Resource {
	return &NeoliteProVmResource{}
}

func (r *NeoliteProVmResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_neolite_pro_vm"
}

func (r *NeoliteProVmResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NEO Lite Pro virtual machine. Order VM, update nama/keypair/package/storage/power, dan rebuild OS.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource id = account_id dari BiznetGIO.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"order_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Order id dari response create.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"ssh_and_console_user": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "User SSH & console yang dipasang saat create. Panjang 6-32, hanya huruf/angka/titik/dash.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(6, 32),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9\-\.]*$`), "hanya huruf, angka, titik, atau dash"),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"console_password": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "Password console saat create. Write-only: ga pernah di-refetch dari API. Minimal 8 karakter alphanumeric, harus ada huruf kecil, huruf besar, dan angka.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(8),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9]+$`), "hanya alphanumeric"),
					stringvalidator.RegexMatches(regexp.MustCompile(`[a-z]`), "harus ada huruf kecil"),
					stringvalidator.RegexMatches(regexp.MustCompile(`[A-Z]`), "harus ada huruf besar"),
					stringvalidator.RegexMatches(regexp.MustCompile(`\d`), "harus ada angka"),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"vm_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("server-name"),
				MarkdownDescription: "Nama VM. Default `server-name`. Bisa diubah via change-vm-name. Panjang 6-16, hanya huruf/angka/titik/dash.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(6, 16),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9\-\.]*$`), "hanya huruf, angka, titik, atau dash"),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Deskripsi VM. Create-only: ga ada endpoint update, jadi nilai dari API yang dipakai.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"product_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Product id dari data source `biznetgio_neolite_pro_products` atau portal. Bisa diubah via change-package.",
			},
			"select_os": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "OS yang dipasang saat create, dari data source `biznetgio_neolite_pro_os_list`. Ganti OS = pakai `rebuild_os`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"keypair_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Id keypair dari `biznetgio_neolite_pro_keypair`. Bisa diganti via change-keypair.",
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
				MarkdownDescription: "Bayar invoice pake kartu kredit saat order. Default true (auto-charge). Set false kalau mau ninggalin invoice unpaid di portal — resource bakal stuck Pending sampai dibayar. Bisa diubah in-place via change-package/storage.",
			},
			"promocode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Kode promo saat order.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"power_state": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Power state VM: `start`, `stop`, `suspend`, `resume`, atau `shutdown`. Update cuma mengirim action kalau nilainya berubah.",
				Validators:          []validator.String{enumValidator{vals: []string{"start", "stop", "suspend", "resume", "shutdown"}}},
			},
			"rebuild_os": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Kalau berubah, VM di-rebuild (wipe OS) pake OS baru via endpoint rebuild. List OS valid ada di data source `biznetgio_neolite_pro_os_list`.",
			},
			"disk_size": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
// wip 386
