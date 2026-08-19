package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/shirasakaren/terraform-provider-biznetgio/internal/biznetgio"
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

func (r *BaremetalAdditionalIPAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data BaremetalAdditionalIPAssignmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, err := r.client.BaremetalAdditionalIP().Assign(ctx, data.AdditionalIPID.ValueInt64(),
		biznetgio.AssignToMachineRequest{MetalAccountID: data.MetalAccountID.ValueInt64()})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to assign additional ip: %s", err))
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%d:%d", data.MetalAccountID.ValueInt64(), data.AdditionalIPID.ValueInt64()))
	data.Status = types.StringValue(aliasStr(out, "status", "state"))
	data.Raw = types.StringValue(rawJSON(out))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update gak ada endpointnya - semua input RequiresReplace, method ini cuma formalitas interface.
func (r *BaremetalAdditionalIPAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data BaremetalAdditionalIPAssignmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BaremetalAdditionalIPAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data BaremetalAdditionalIPAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ipID, metalID, ok := parseAssignmentID(data.ID.ValueString())
	if !ok {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Invalid assignment id: %q", data.ID.ValueString()))
		return
	}

	out, err := r.client.BaremetalAdditionalIP().AssignGet(ctx, ipID, metalID)
	if err != nil {
		if biznetgio.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read assignment: %s", err))
		return
	}

	data.AdditionalIPID = types.Int64Value(ipID)
	data.MetalAccountID = types.Int64Value(metalID)
	data.Status = types.StringValue(aliasStr(out, "status", "state"))
	data.Raw = types.StringValue(rawJSON(out))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BaremetalAdditionalIPAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data BaremetalAdditionalIPAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ipID, metalID, ok := parseAssignmentID(data.ID.ValueString())
	if !ok {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Invalid assignment id: %q", data.ID.ValueString()))
		return
	}
	if err := r.client.BaremetalAdditionalIP().Unassign(ctx, ipID, metalID); err != nil && !biznetgio.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to unassign additional ip: %s", err))
	}
}

func (r *BaremetalAdditionalIPAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ipID, metalID, ok := parseAssignmentID(req.ID)
	if !ok {
		resp.Diagnostics.AddError("Invalid Import ID", "Import id harus format `<metal_account_id>:<additional_ip_id>`")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("metal_account_id"), metalID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("additional_ip_id"), ipID)...)
}

func parseAssignmentID(id string) (ipID, metalID int64, ok bool) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return 0, 0, false
	}
	metal, err1 := strconv.ParseInt(parts[0], 10, 64)
	ip, err2 := strconv.ParseInt(parts[1], 10, 64)
	if err1 != nil || err2 != nil {
		return 0, 0, false
	}
	return ip, metal, true
}
