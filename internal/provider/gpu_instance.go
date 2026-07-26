// SPDX-License-Identifier: MPL-2.0

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
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.Resource                   = &GpuInstanceResource{}
	_ resource.ResourceWithImportState    = &GpuInstanceResource{}
	_ resource.ResourceWithValidateConfig = &GpuInstanceResource{}
)

type GpuInstanceResource struct {
	client *biznetgio.Client
}

type GpuInstanceResourceModel struct {
	ID                            types.String   `tfsdk:"id"`
	ProductID                     types.Int64    `tfsdk:"product_id"`
	SelectOS                      types.String   `tfsdk:"select_os"`
	KeypairID                     types.Int64    `tfsdk:"keypair_id"`
	ServiceName                   types.String   `tfsdk:"service_name"`
	SSHAndConsoleUser             types.String   `tfsdk:"ssh_and_console_user"`
	ConsolePassword               types.String   `tfsdk:"console_password"`
	Promocode                     types.String   `tfsdk:"promocode"`
	PayWithCreditCard             types.Bool     `tfsdk:"pay_with_credit_card"`
	Subscription                  types.Object   `tfsdk:"subscription"`
	OnDemand                      types.Object   `tfsdk:"on_demand"`
	RebuildTrigger                types.String   `tfsdk:"rebuild_trigger"`
	ReserveAdditionalHoursTrigger types.String   `tfsdk:"reserve_additional_hours_trigger"`
	Status                        types.String   `tfsdk:"status"`
	OrderID                       types.String   `tfsdk:"order_id"`
	Raw                           types.String   `tfsdk:"raw"`
	Timeouts                      timeouts.Value `tfsdk:"timeouts"`
}

type GpuSubscriptionModel struct {
	Cycle types.String `tfsdk:"cycle"`
}

type GpuOnDemandModel struct {
	AdditionalHours types.Int64 `tfsdk:"additional_hours"`
}

func gpuSubscriptionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cycle": types.StringType,
	}
}

func gpuOnDemandAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"additional_hours": types.Int64Type,
	}
}

func (r *GpuInstanceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "biznetgio_gpu_instance"
}

func (r *GpuInstanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NEO GPU instance. Billing mode is either `subscription` (cycle-based) or `on_demand` (hourly), exactly one must be set.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "GPU instance account id.",
			},
			"product_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "GPU product id from `biznetgio_gpu_products`.",
			},
			"select_os": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "OS to install, from the product's select-os catalog.",
			},
			"keypair_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "GPU keypair id, from `biznetgio_gpu_keypair`.",
			},
			"service_name": schema.StringAttribute{
				Optional:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				MarkdownDescription: "Display name of the instance. Create-only, changing it replaces the instance.",
			},
			"ssh_and_console_user": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "SSH and console user for the instance.",
			},
			"console_password": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				MarkdownDescription: "Console password. Create-only, changing it replaces the instance.",
			},
			"promocode": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Promo code to apply at creation.",
			},
			"pay_with_credit_card": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Pay the invoice with the registered credit card. Defaults to true. Set false to leave the invoice unpaid in the portal; the resource stays Pending until paid.",
			},
			"subscription": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Subscription billing mode (cycle-based). Exactly one of `subscription` or `on_demand` must be set.",
				Attributes: map[string]schema.Attribute{
					"cycle": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Billing cycle, e.g. `m` for monthly or `a` for annual.",
					},
				},
			},
			"on_demand": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "On-demand billing mode (hourly). Exactly one of `subscription` or `on_demand` must be set.",
				Attributes: map[string]schema.Attribute{
					"additional_hours": schema.Int64Attribute{
						Optional:            true,
						Computed:            true,
						Default:             int64default.StaticInt64(0),
						MarkdownDescription: "Additional hours to reserve. Defaults to 0.",
					},
				},
			},
			"rebuild_trigger": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Change this value to rebuild the instance (destructive, reinstalls the OS).",
			},
			"reserve_additional_hours_trigger": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Change this value to reserve 1 additional hour on an on-demand instance.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Current status of the instance.",
			},
			"order_id": schema.StringAttribute{
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Order id from the creation response.",
			},
			"raw": schema.StringAttribute{
				Sensitive:           true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Raw JSON of the last read response, for anything not modeled yet.",
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
		},
	}
}

func (r *GpuInstanceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *GpuInstanceResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data GpuInstanceResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	hasSub := !data.Subscription.IsNull() && !data.Subscription.IsUnknown()
	hasOnDemand := !data.OnDemand.IsNull() && !data.OnDemand.IsUnknown()
	if hasSub == hasOnDemand {
		resp.Diagnostics.AddError(
			"Invalid Billing Mode",
			"exactly one of `subscription` or `on_demand` must be set on biznetgio_gpu_instance",
		)
	}
}

func (r *GpuInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data GpuInstanceResourceModel
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

	pay := "yes"
	if !data.PayWithCreditCard.IsNull() && !data.PayWithCreditCard.ValueBool() {
		pay = "no"
	}

	var raw map[string]any
	var err error
	if !data.Subscription.IsNull() {
		var sub GpuSubscriptionModel
		resp.Diagnostics.Append(data.Subscription.As(ctx, &sub, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}
		raw, err = r.client.GPU().Create(ctx, biznetgio.GPUCreateRequest{
			ProductID:         data.ProductID.ValueInt64(),
			SelectOS:          data.SelectOS.ValueString(),
			KeypairID:         data.KeypairID.ValueInt64(),
			ServiceName:       data.ServiceName.ValueString(),
			SSHAndConsoleUser: data.SSHAndConsoleUser.ValueString(),
			ConsolePassword:   data.ConsolePassword.ValueString(),
			Promocode:         data.Promocode.ValueString(),
			PayInvoiceWithCC:  pay,
			Cycle:             sub.Cycle.ValueString(),
		})
	} else if !data.OnDemand.IsNull() {
		var od GpuOnDemandModel
		resp.Diagnostics.Append(data.OnDemand.As(ctx, &od, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}
		raw, err = r.client.GPU().CreateOneTime(ctx, biznetgio.GPUOneTimeCreateRequest{
			ProductID:         data.ProductID.ValueInt64(),
			SelectOS:          data.SelectOS.ValueString(),
			KeypairID:         data.KeypairID.ValueInt64(),
			ServiceName:       data.ServiceName.ValueString(),
			SSHAndConsoleUser: data.SSHAndConsoleUser.ValueString(),
			ConsolePassword:   data.ConsolePassword.ValueString(),
			Promocode:         data.Promocode.ValueString(),
			PayInvoiceWithCC:  pay,
			AdditionalHours:   od.AdditionalHours.ValueInt64(),
		})
	} else {
		resp.Diagnostics.AddError("Invalid Billing Mode", "exactly one of `subscription` or `on_demand` must be set on biznetgio_gpu_instance")
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to create gpu instance: %s", err))
		return
	}

	accountID, ok := gpuInt64(raw, "account_id", "id")
	if !ok {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("gpu create response missing account_id: %v", raw))
		return
	}
	data.ID = types.StringValue(strconv.FormatInt(accountID, 10))
	if v, ok := gpuString(raw, "order_id"); ok {
		data.OrderID = types.StringValue(v)
	}

	final, err := biznetgio.WaitForStatus(ctx, 5*time.Second,
		func(ctx context.Context) (map[string]any, error) { return r.client.GPU().AccountGet(ctx, accountID) },
		gpuStatusOf, []string{"active"}, []string{"terminated", "failed", "error", "deleted", "cancelled"},
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to wait for gpu instance %d to become active: %s", accountID, err))
		return
	}

	resp.Diagnostics.Append(gpuInstanceSetFromMap(ctx, &data, final)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Raw = types.StringValue(string(redactJSON(final)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GpuInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data GpuInstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountID, err := strconv.ParseInt(data.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("invalid gpu instance id %q: %s", data.ID.ValueString(), err))
		return
	}

	items, err := r.client.GPU().AccountsList(ctx, "")
	if err != nil {
		if biznetgio.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to list gpu instances: %s", err))
		return
	}

	var found map[string]any
	for _, it := range items {
		if id, ok := gpuInt64(it, "account_id", "id"); ok && id == accountID {
			found = it
			break
		}
	}
	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(gpuInstanceSetFromMap(ctx, &data, found)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Raw = types.StringValue(string(redactJSON(found)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GpuInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state GpuInstanceResourceModel
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

	accountID, err := strconv.ParseInt(plan.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("invalid gpu instance id %q: %s", plan.ID.ValueString(), err))
		return
	}

	rebuild := plan.RebuildTrigger.ValueString()
	if rebuild != "" && rebuild != state.RebuildTrigger.ValueString() {
		if _, err := r.client.GPU().Rebuild(ctx, accountID, biznetgio.GPURebuildRequest{SelectOS: plan.SelectOS.ValueString()}); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to rebuild gpu instance %d: %s", accountID, err))
			return
		}
	}
	reserve := plan.ReserveAdditionalHoursTrigger.ValueString()
	if reserve != "" && reserve != state.ReserveAdditionalHoursTrigger.ValueString() {
		if _, err := r.client.GPU().ReserveAdditionalHours(ctx, accountID, biznetgio.ReserveAdditionalHoursRequest{Hours: 1}); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to reserve additional hours for gpu instance %d: %s", accountID, err))
			return
		}
	}

	final, err := biznetgio.WaitForStatus(ctx, 5*time.Second,
		func(ctx context.Context) (map[string]any, error) { return r.client.GPU().AccountGet(ctx, accountID) },
		gpuStatusOf, []string{"active"}, []string{"terminated", "failed", "error", "deleted", "cancelled"},
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to wait for gpu instance %d: %s", accountID, err))
		return
	}

	data := plan
	resp.Diagnostics.Append(gpuInstanceSetFromMap(ctx, &data, final)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Raw = types.StringValue(string(redactJSON(final)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GpuInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data GpuInstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountID, err := strconv.ParseInt(data.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("invalid gpu instance id %q: %s", data.ID.ValueString(), err))
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, 10*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

