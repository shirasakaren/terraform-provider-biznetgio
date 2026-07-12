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
