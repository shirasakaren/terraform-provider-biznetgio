package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/shirasakaren/terraform-provider-biznetgio/internal/biznetgio"
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
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Kalau diisi, keypair di-import via `POST /baremetals/keypairs/import`. Kalau kosong, keypair di-generate server. Output-nya selalu public key dari API.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"private_key": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "Private key dari response create (kalau API ngasih). Cuma muncul sekali; di Read berikutnya di-keep dari state.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"raw": schema.StringAttribute{
				Sensitive:           true,
				Computed:            true,
				MarkdownDescription: "Full JSON item keypair dari list API.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *BaremetalKeypairResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BaremetalKeypairResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data BaremetalKeypairResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var out map[string]any
	var err error
	if data.PublicKey.IsNull() || data.PublicKey.ValueString() == "" {
		out, err = r.client.Baremetal().KeypairCreate(ctx, biznetgio.KeypairCreateRequest{Name: data.Name.ValueString()})
	} else {
		out, err = r.client.Baremetal().KeypairImport(ctx, biznetgio.KeypairImportRequest{Name: data.Name.ValueString(), PublicKey: data.PublicKey.ValueString()})
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create baremetal keypair: %s", err))
		return
	}

	keypairID := aliasInt(out, "keypair_id", "id")
	if keypairID == 0 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Create keypair response tidak ada id: %s", rawJSON(out)))
		return
	}
	data.ID = types.StringValue(strconv.FormatInt(keypairID, 10))
	data.KeypairID = types.Int64Value(keypairID)
	if v := aliasStr(out, "public_key", "publickey"); v != "" {
		data.PublicKey = types.StringValue(v)
	}
	// private key cuma di response create — alias defensif, jangan nebak nama
	if v := aliasStr(out, "private_key", "private", "secret_key", "pem"); v != "" {
		data.PrivateKey = types.StringValue(v)
	}
	data.Raw = types.StringValue(rawJSON(out))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update gak ada endpointnya — semua input RequiresReplace, method ini cuma formalitas interface.
func (r *BaremetalKeypairResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data BaremetalKeypairResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BaremetalKeypairResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data BaremetalKeypairResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	keypairID, err := strconv.ParseInt(data.ID.ValueString(), 10, 64)
	if err != nil || keypairID == 0 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Invalid keypair id: %q", data.ID.ValueString()))
		return
	}

	items, err := r.client.Baremetal().KeypairList(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list baremetal keypairs: %s", err))
		return
	}
	var found map[string]any
	for _, it := range items {
		if aliasInt(it, "keypair_id", "id") == keypairID {
			found = it
			break
		}
	}
	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	data.KeypairID = types.Int64Value(keypairID)
	if v := aliasStr(found, "name"); v != "" {
		data.Name = types.StringValue(v)
	}
	if v := aliasStr(found, "public_key", "publickey"); v != "" {
		data.PublicKey = types.StringValue(v)
	}
	// list gak bawa private key — keep value lama
	data.Raw = types.StringValue(rawJSON(found))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BaremetalKeypairResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data BaremetalKeypairResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	keypairID, err := strconv.ParseInt(data.ID.ValueString(), 10, 64)
	if err != nil || keypairID == 0 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Invalid keypair id: %q", data.ID.ValueString()))
		return
	}
	if err := r.client.Baremetal().KeypairDelete(ctx, keypairID); err != nil && !biznetgio.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete baremetal keypair %d: %s", keypairID, err))
	}
}

func (r *BaremetalKeypairResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
