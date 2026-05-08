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

type NeoliteKeypairResource struct {
	client *biznetgio.Client
}

type NeoliteKeypairResourceModel struct {
	ID         types.String `tfsdk:"id"`
	KeypairID  types.Int64  `tfsdk:"keypair_id"`
	Name       types.String `tfsdk:"name"`
	PublicKey  types.String `tfsdk:"public_key"`
	PrivateKey types.String `tfsdk:"private_key"`
}

func NewNeoliteKeypairResource() resource.Resource {
	return &NeoliteKeypairResource{}
}

func (r *NeoliteKeypairResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_neolite_keypair"
}

func (r *NeoliteKeypairResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "SSH keypair untuk NEO Lite (pool keypair neolite, terpisah dari baremetal/gpu). Keypair di-generate server; private key cuma keluar sekali di response create.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource id = keypair_id.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"keypair_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Id keypair di BiznetGIO.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Nama keypair.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"public_key": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Public key yang di-generate.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"private_key": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "Private key (sensitive). Write-only: cuma ada di response create, ga bisa di-refetch.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *NeoliteKeypairResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NeoliteKeypairResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NeoliteKeypairResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// pake raw response: field private key undocumented, aliasnya bisa beda-beda.
	raw, err := r.client.Neolite().KeypairCreateRaw(ctx, biznetgio.KeypairCreateRequest{Name: data.Name.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create neolite keypair: %s", err))
		return
	}

	keypairID := aliasInt(raw, "keypair_id", "neosshkey_id", "id")
	if keypairID == 0 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Create neolite keypair response tidak ada keypair_id: %s", rawJSON(raw)))
		return
	}
	data.ID = types.StringValue(strconv.FormatInt(keypairID, 10))
	data.KeypairID = types.Int64Value(keypairID)
	data.PublicKey = types.StringValue(aliasStr(raw, "public_key", "pubkey"))
	data.PrivateKey = types.StringValue(aliasStr(raw, "private_key", "private", "secret_key", "pem"))
	if data.PrivateKey.ValueString() == "" {
		resp.Diagnostics.AddWarning("Private key tidak ditemukan",
			fmt.Sprintf("Response create ga memuat private key (field undocumented), coba cek: %s", rawJSON(raw)))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NeoliteKeypairResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NeoliteKeypairResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	list, err := r.client.Neolite().KeypairList(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list neolite keypairs: %s", err))
		return
	}
