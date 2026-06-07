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

// wip 570
// wip 605
