package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/biznetgio/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BaremetalKeypairResource struct {
	client *biznetgio.Client
}

type BaremetalKeypairResourceModel struct {
	ID         types.String `tfsdk:"id"`
	KeypairID  types.Int64  `tfsdk:"keypair_id"`
	Name       types.String `tfsdk:"name"`
	PublicKey  types.String `tfsdk:"public_key"`
	PrivateKey types.String `tfsdk:"private_key"`
	Raw        types.String `tfsdk:"raw"`
}

func NewBaremetalKeypairResource() resource.Resource {
	return &BaremetalKeypairResource{}
}

func (r *BaremetalKeypairResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_baremetal_keypair"
}

func (r *BaremetalKeypairResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "SSH keypair untuk NEO Metal (pool keypair baremetal, terpisah dari neolite/gpu). Keypair di-generate server; private key cuma keluar sekali di response create.",
// wip 664
