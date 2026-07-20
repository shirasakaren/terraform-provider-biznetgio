package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/biznetgio/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BaremetalAdditionalIPAssignmentResource struct {
	client *biznetgio.Client
}

type BaremetalAdditionalIPAssignmentResourceModel struct {
	ID             types.String `tfsdk:"id"`
	AdditionalIPID types.Int64  `tfsdk:"additional_ip_id"`
	MetalAccountID types.Int64  `tfsdk:"metal_account_id"`
	Status         types.String `tfsdk:"status"`
	Raw            types.String `tfsdk:"raw"`
}

func NewBaremetalAdditionalIPAssignmentResource() resource.Resource {
	return &BaremetalAdditionalIPAssignmentResource{}
}

func (r *BaremetalAdditionalIPAssignmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_baremetal_additional_ip_assignment"
}

func (r *BaremetalAdditionalIPAssignmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Attach additional IP ke baremetal server via `PUT /baremetal-additional-ips/{account_id}/assigns`. Composite id: `<metal_account_id>:<additional_ip_id>`. Ganti target = destroy + recreate.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite id `<metal_account_id>:<additional_ip_id>`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"additional_ip_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Account id dari `biznetgio_baremetal_additional_ip` yang mau di-attach.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"metal_account_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Account id baremetal tujuan attach.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Status assignment kalau ada di response.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"raw": schema.StringAttribute{
				Sensitive:           true,
				Computed:            true,
				MarkdownDescription: "Full JSON response terakhir dari API.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *BaremetalAdditionalIPAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
