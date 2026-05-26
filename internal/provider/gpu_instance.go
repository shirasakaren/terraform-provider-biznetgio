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
