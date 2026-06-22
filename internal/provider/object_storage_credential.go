package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/biznetgio/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ObjectStorageCredentialResource{}

type ObjectStorageCredentialResource struct {
	client *biznetgio.Client
}

type ObjectStorageCredentialResourceModel struct {
	ID        types.String `tfsdk:"id"`
	AccountID types.String `tfsdk:"account_id"`
	AccessKey types.String `tfsdk:"access_key"`
	SecretKey types.String `tfsdk:"secret_key"`
	Active    types.Bool   `tfsdk:"active"`
	Raw       types.String `tfsdk:"raw"`
}

func NewObjectStorageCredentialResource() resource.Resource {
	return &ObjectStorageCredentialResource{}
}

func (r *ObjectStorageCredentialResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "biznetgio_object_storage_credential"
}

func (r *ObjectStorageCredentialResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "S3 credential (access/secret key pair) for an Object Storage instance. " +
			"The secret key is shown only once at create time and cannot be re-fetched. " +
			"Import with `terraform import biznetgio_object_storage_credential.example <account_id>:<access_key>`; " +
			"the id is then normalized to the hashed form `<account_id>:<sha256(access_key)[:16]>`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Composite id `<account_id>:<sha256 hex dari access key, 16 char pertama>` — access key plaintext ga pernah masuk id.",
			},
			"account_id": schema.StringAttribute{
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				MarkdownDescription: "Object Storage instance account id.",
			},
			"access_key": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Credential access key, returned once at create.",
			},
			"secret_key": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Credential secret key, shown only once at create; keeps its last state value on refresh.",
			},
			"active": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Whether the credential is enabled. Set false to disable without deleting.",
			},
			"raw": schema.StringAttribute{
				Sensitive:           true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Full JSON of the last read from the API, for anything not modeled yet.",
			},
		},
	}
}

func (r *ObjectStorageCredentialResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*biznetgio.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *biznetgio.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}
	r.client = client
}

func (r *ObjectStorageCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ObjectStorageCredentialResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountID, err := objParseAccountID(data.AccountID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid account_id", err.Error())
		return
	}

	created, err := r.client.ObjectStorage().CredentialCreate(ctx, accountID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create credential: %s", err))
		return
	}

	accessKey, ok := objMapString(created, "accessKey", "access_key", "accesskey")
	if !ok {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Credential create response missing access key: %s", objRawString(created)))
		return
	}
	if secret, ok := objMapString(created, "secretKey", "secret_key", "secretkey"); ok {
		data.SecretKey = types.StringValue(secret)
	}
	if v, ok := objMapBool(created, "active"); ok {
		data.Active = types.BoolValue(v)
	}

	data.AccessKey = types.StringValue(accessKey)
	data.ID = types.StringValue(data.AccountID.ValueString() + ":" + objHashKey(accessKey))
	data.Raw = objRawString(created)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ObjectStorageCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ObjectStorageCredentialResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountID, keyForm, err := objParseCredentialID(data.ID.ValueString(), data.AccountID.ValueString(), data.AccessKey.ValueString())
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	item, ok, err := objFindCredential(ctx, r.client, accountID, keyForm)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list credentials: %s", err))
		return
	}
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}

	accessKey, ok := objMapString(item, "accessKey", "access_key", "accesskey")
	if !ok {
// wip 741
