package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

type NeoliteVMResource struct {
	client *biznetgio.Client
}

type NeoliteVMResourceModel struct {
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
	MigrateToPro      types.String            `tfsdk:"migrate_to_pro"`
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

type NeoliteLastInvoiceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	PaidID      types.Int64  `tfsdk:"paid_id"`
	Status      types.String `tfsdk:"status"`
	Date        types.String `tfsdk:"date"`
	Duedate     types.String `tfsdk:"duedate"`
	Paybefore   types.String `tfsdk:"paybefore"`
	Datepaid    types.String `tfsdk:"datepaid"`
	InvoiceType types.String `tfsdk:"invoice_type"`
}

func NewNeoliteVMResource() resource.Resource {
	return &NeoliteVMResource{}
}

func (r *NeoliteVMResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_neolite_vm"
}

func (r *NeoliteVMResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NEO Lite virtual machine. Order VM, update nama/keypair/package/storage/power, rebuild OS, atau migrate ke NEO Lite Pro.",
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
				MarkdownDescription: "Product id dari data source `biznetgio_neolite_products` atau portal.",
			},
			"select_os": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "OS yang dipasang saat create, dari data source `biznetgio_neolite_os_list`. Ganti OS = pakai `rebuild_os`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"keypair_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Id keypair dari `biznetgio_neolite_keypair`. Bisa diganti via change-keypair.",
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
				MarkdownDescription: "Kalau berubah, VM di-rebuild (wipe OS) pake OS baru via endpoint rebuild. List OS valid ada di data source `biznetgio_neolite_os_list`.",
			},
			"migrate_to_pro": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Trigger one-shot migrate ke NEO Lite Pro: isi `neolitepro_product_id` target (lihat data source `biznetgio_neolite_change_package_options` family). Ganti nilainya buat re-trigger.",
			},
			"disk_size": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Ukuran disk target (GB, absolute — bukan tambahan). Cuma bisa naik, bukan turun.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Status akun terakhir dari API (Active, Pending, Suspended, Terminated).",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"uptime": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Uptime VM dalam detik.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"maxdisk": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Ukuran disk maksimal VM (GB).",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"maxmem": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Memory maksimal VM (MB).",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"mem": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Memory yang dipakai VM (MB).",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"cpus": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Jumlah CPU VM.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"ciuser": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Cloud-init user VM.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"cipassword": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "Cloud-init password VM (sensitive).",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"osname": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Nama OS yang jalan di VM.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"region": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Region VM.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"region_label": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Label region VM.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"next_due": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Tanggal tagihan berikutnya.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"recurring_amount": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Nominal recurring per siklus.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"billingcycle": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Siklus billing aktif.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"product_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Nama product aktif.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"last_invoice": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Invoice terakhir VM.",
				Attributes: map[string]schema.Attribute{
					"id":           schema.Int64Attribute{Computed: true, MarkdownDescription: "Id invoice."},
					"paid_id":      schema.Int64Attribute{Computed: true, MarkdownDescription: "Id pembayaran invoice."},
					"status":       schema.StringAttribute{Computed: true, MarkdownDescription: "Status invoice."},
					"date":         schema.StringAttribute{Computed: true, MarkdownDescription: "Tanggal invoice."},
					"duedate":      schema.StringAttribute{Computed: true, MarkdownDescription: "Tanggal jatuh tempo."},
					"paybefore":    schema.StringAttribute{Computed: true, MarkdownDescription: "Batas bayar."},
					"datepaid":     schema.StringAttribute{Computed: true, MarkdownDescription: "Tanggal dibayar."},
					"invoice_type": schema.StringAttribute{Computed: true, MarkdownDescription: "Tipe invoice."},
				},
			},
			"raw": schema.StringAttribute{
				Sensitive:           true,
				Computed:            true,
				MarkdownDescription: "Full JSON response akun terakhir dari API, buat akses field yang belum dimodel (cipassword di-mask).",
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

func (r *NeoliteVMResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NeoliteVMResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NeoliteVMResourceModel
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
	billing, err := r.client.Neolite().VMCreate(ctx, biznetgio.NeoliteCreateRequest{
		ProductID:         data.ProductID.ValueInt64(),
		Cycle:             data.Cycle.ValueString(),
		SelectOS:          data.SelectOS.ValueString(),
		KeypairID:         data.KeypairID.ValueInt64(),
		VMName:            data.VMName.ValueString(),
		Description:       data.Description.ValueString(),
		SSHAndConsoleUser: data.SSHAndConsoleUser.ValueString(),
		ConsolePassword:   data.ConsolePassword.ValueString(),
		Promocode:         data.Promocode.ValueString(),
		PayInvoiceWithCC:  cc,
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create neolite vm: %s", err))
		return
	}
	if billing.AccountID == "" {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Create neolite vm response tidak ada account_id: order_id=%s", billing.OrderID))
		return
	}
	data.ID = types.StringValue(billing.AccountID)
	data.OrderID = types.StringValue(billing.OrderID)

	tflog.Info(ctx, "neolite vm created, menunggu active", map[string]any{"account_id": billing.AccountID})
	acc, err := biznetgio.WaitForStatus(ctx, 5*time.Second,
		func(ctx context.Context) (biznetgio.AccountResource, error) {
			return r.client.Neolite().AccountGet(ctx, parseAccountID(billing.AccountID))
		},
		accountStatus, []string{"active"}, []string{"suspended", "terminated"})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Neolite vm %s gagal jadi active: %s", billing.AccountID, err))
		return
	}
	data.Status = types.StringValue(acc.Status)

	if err := r.refresh(ctx, &data); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read neolite vm %s: %s", billing.AccountID, err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NeoliteVMResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NeoliteVMResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.refresh(ctx, &data)
	if err != nil {
		if biznetgio.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read neolite vm %s: %s", data.ID.ValueString(), err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NeoliteVMResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state NeoliteVMResourceModel
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

	accountID := parseAccountID(state.ID.ValueString())

	if !plan.VMName.Equal(state.VMName) {
		if err := r.client.Neolite().VMChangeName(ctx, accountID, plan.VMName.ValueString()); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to change neolite vm name: %s", err))
			return
		}
		state.VMName = plan.VMName
	}
	if !plan.KeypairID.Equal(state.KeypairID) {
		if err := r.client.Neolite().VMChangeKeypair(ctx, accountID, plan.KeypairID.ValueInt64()); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to change neolite vm keypair: %s", err))
			return
		}
		state.KeypairID = plan.KeypairID
	}

	needsPoll := false
	if !plan.ProductID.Equal(state.ProductID) {
		if _, err := r.client.Neolite().VMChangePackage(ctx, accountID, biznetgio.NeoliteChangePackageRequest{
			NewProductID:     plan.ProductID.ValueInt64(),
			PayInvoiceWithCC: ccYesNo(plan.PayWithCreditCard),
		}); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to change neolite vm package: %s", err))
			return
		}
		state.ProductID = plan.ProductID
		needsPoll = true
	}
	if !plan.DiskSize.Equal(state.DiskSize) {
		newSize := plan.DiskSize.ValueInt64()
		oldSize := state.DiskSize.ValueInt64()
		if newSize < oldSize {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Neolite vm storage cuma bisa di-upgrade: %d -> %d", oldSize, newSize))
			return
		}
		if _, err := r.client.Neolite().VMChangeStorage(ctx, accountID, biznetgio.NeoliteUpgradeStorageRequest{
			DiskSize:         newSize,
			PayInvoiceWithCC: ccYesNo(plan.PayWithCreditCard),
		}); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to change neolite vm storage: %s", err))
			return
		}
		state.DiskSize = plan.DiskSize
		needsPoll = true
	}
	if !plan.PowerState.IsUnknown() && !plan.PowerState.Equal(state.PowerState) {
		ps := plan.PowerState.ValueString()
		if err := r.client.Neolite().VMState(ctx, accountID, ps); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to set neolite vm power state %q: %s", ps, err))
			return
		}
		state.PowerState = plan.PowerState
	}
	if !plan.RebuildOS.IsUnknown() && !plan.RebuildOS.Equal(state.RebuildOS) {
		if err := r.client.Neolite().VMRebuild(ctx, accountID, plan.RebuildOS.ValueString()); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to rebuild neolite vm: %s", err))
			return
		}
		state.RebuildOS = plan.RebuildOS
		needsPoll = true
	}
	if !plan.MigrateToPro.IsUnknown() && !plan.MigrateToPro.Equal(state.MigrateToPro) {
		proProductID, err := strconv.ParseInt(plan.MigrateToPro.ValueString(), 10, 64)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("migrate_to_pro harus berisi neolitepro_product_id numeric: %q", plan.MigrateToPro.ValueString()))
			return
		}
		if _, err := r.client.Neolite().MigrateToPro(ctx, accountID, biznetgio.MigrateToProRequest{
			NeoliteProProductID: proProductID,
			PayInvoiceWithCC:    ccYesNo(plan.PayWithCreditCard),
		}); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to migrate neolite vm to pro: %s", err))
			return
		}
		state.MigrateToPro = plan.MigrateToPro
		needsPoll = true
	}

	// description ga punya endpoint update — state tetap bawa nilai server dari Read/refresh.

	if needsPoll {
		tflog.Info(ctx, "neolite vm action dikirim, menunggu active", map[string]any{"account_id": state.ID.ValueString()})
		acc, err := biznetgio.WaitForStatus(ctx, 5*time.Second,
			func(ctx context.Context) (biznetgio.AccountResource, error) {
				return r.client.Neolite().AccountGet(ctx, accountID)
			},
			accountStatus, []string{"active"}, []string{"suspended", "terminated"})
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Neolite vm %s gagal balik active: %s", state.ID.ValueString(), err))
			return
		}
		state.Status = types.StringValue(acc.Status)
	}

	if err := r.refresh(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read neolite vm %s: %s", state.ID.ValueString(), err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NeoliteVMResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NeoliteVMResourceModel
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

	err := r.client.Neolite().VMDelete(ctx, parseAccountID(data.ID.ValueString()))
	if err != nil && !biznetgio.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete neolite vm %s: %s", data.ID.ValueString(), err))
		return
	}
}

func (r *NeoliteVMResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// refresh isi semua computed field dari AccountGet + VMDetails.
func (r *NeoliteVMResource) refresh(ctx context.Context, data *NeoliteVMResourceModel) error {
	accountID := parseAccountID(data.ID.ValueString())

	acc, err := r.client.Neolite().AccountGet(ctx, accountID)
	if err != nil {
		return err
	}

	data.Status = types.StringValue(acc.Status)
	data.BillingCycle = types.StringValue(acc.Billingcycle)
	data.NextDue = types.StringValue(acc.NextDue)
	data.RecurringAmount = types.Int64Value(acc.RecurringAmount)
	data.ProductID = types.Int64Value(acc.ProductID)
	data.ProductName = types.StringValue(acc.ProductName)
	data.Description = types.StringValue(acc.Description)
	data.Region = types.StringValue(acc.ExtraDetails.Region)
	data.RegionLabel = types.StringValue(acc.ExtraDetails.RegionLabel)
	data.CIUser = types.StringValue(acc.ExtraDetails.CIUser)
	data.CIPassword = types.StringValue(acc.ExtraDetails.CIPassword)
	data.OSName = types.StringValue(acc.ExtraDetails.OSName)
	if acc.ExtraDetails.KeypairID != 0 {
		data.KeypairID = types.Int64Value(acc.ExtraDetails.KeypairID)
	}
	if v := acc.ExtraDetails.DiskSize; v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			data.DiskSize = types.Int64Value(n)
		}
	}
	if v := acc.ExtraDetails.Name; v != "" {
		data.VMName = types.StringValue(v)
	}
	data.LastInvoice = NeoliteLastInvoiceModel{
		ID:          types.Int64Value(acc.LastInvoice.ID),
		PaidID:      types.Int64Value(acc.LastInvoice.PaidID),
		Status:      types.StringValue(acc.LastInvoice.Status),
		Date:        types.StringValue(acc.LastInvoice.Date),
		Duedate:     types.StringValue(acc.LastInvoice.Duedate),
// wip 792
