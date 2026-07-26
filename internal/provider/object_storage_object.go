package provider

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/biznetgio/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ObjectStorageObjectResource{}

type ObjectStorageObjectResource struct {
	client *biznetgio.Client
}

type ObjectStorageObjectResourceModel struct {
	ID        types.String `tfsdk:"id"`
	AccountID types.String `tfsdk:"account_id"`
	Bucket    types.String `tfsdk:"bucket"`
	Key       types.String `tfsdk:"key"`
	Source    types.String `tfsdk:"source"`
	Content   types.String `tfsdk:"content"`
	ACL       types.String `tfsdk:"acl"`
	Raw       types.String `tfsdk:"raw"`
}

func NewObjectStorageObjectResource() resource.Resource {
	return &ObjectStorageObjectResource{}
}

func (r *ObjectStorageObjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "biznetgio_object_storage_object"
}

func (r *ObjectStorageObjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Thin convenience wrapper to upload/delete a single object via the BiznetGIO control-plane API. " +
			"Tidak cocok untuk object besar — untuk scale, pakai tooling S3-compatible langsung dengan credential dari resource credential.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Composite id `<account_id>:<bucket>:<key>`.",
			},
			"account_id": schema.StringAttribute{
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				MarkdownDescription: "Object Storage instance account id.",
			},
			"bucket": schema.StringAttribute{
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				MarkdownDescription: "Bucket name.",
			},
			"key": schema.StringAttribute{
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				MarkdownDescription: "Object key inside the bucket.",
			},
			"source": schema.StringAttribute{
				Optional:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				MarkdownDescription: "Path ke file lokal yang mau di-upload. Exactly one of `source`/`content` required.",
			},
			"content": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				MarkdownDescription: "Inline content yang mau di-upload. Exactly one of `source`/`content` required.",
			},
			"acl": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators:          []validator.String{stringvalidator.OneOf("", "private", "public-read", "public-read-write", "authenticated-read", "log-delivery-write")},
				MarkdownDescription: "S3-style canned ACL applied to the object.",
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

func (r *ObjectStorageObjectResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(tfpath.MatchRoot("source"), tfpath.MatchRoot("content")),
	}
}

func (r *ObjectStorageObjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ObjectStorageObjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ObjectStorageObjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountID, err := objParseAccountID(data.AccountID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid account_id", err.Error())
		return
	}

	content, err := objObjectBytes(data.Source.ValueString(), data.Content.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Object Source", err.Error())
		return
	}

	directory, filename := objSplitKey(data.Key.ValueString())
	if _, err := r.client.ObjectStorage().ObjectUpload(ctx, accountID, data.Bucket.ValueString(), directory, filename, content); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to upload object %q: %s", data.Key.ValueString(), err))
		return
	}

	data.ID = types.StringValue(data.AccountID.ValueString() + ":" + data.Bucket.ValueString() + ":" + data.Key.ValueString())
	if data.ACL.ValueString() != "" {
		if _, err := r.client.ObjectStorage().ObjectSetACL(ctx, accountID, data.Bucket.ValueString(), data.Key.ValueString(), biznetgio.SetACLRequest{
			ACL: data.ACL.ValueString(),
		}); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to set ACL on object %q: %s", data.Key.ValueString(), err))
			return
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ObjectStorageObjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ObjectStorageObjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountID, err := objParseAccountID(data.AccountID.ValueString())
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	item, ok, err := objFindObject(ctx, r.client, accountID, data.Bucket.ValueString(), data.Key.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list objects: %s", err))
		return
	}
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}

	if v, ok := objMapString(item, "acl"); ok {
		data.ACL = types.StringValue(v)
	}
	data.Raw = objRawString(item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ObjectStorageObjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ObjectStorageObjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ACL.Equal(state.ACL) {
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	accountID, err := objParseAccountID(plan.AccountID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid account_id", err.Error())
		return
	}

	if _, err := r.client.ObjectStorage().ObjectSetACL(ctx, accountID, plan.Bucket.ValueString(), plan.Key.ValueString(), biznetgio.SetACLRequest{
		ACL: plan.ACL.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to set ACL on object %q: %s", plan.Key.ValueString(), err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ObjectStorageObjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ObjectStorageObjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountID, err := objParseAccountID(data.AccountID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid account_id", err.Error())
		return
	}

	if err := r.client.ObjectStorage().ObjectDelete(ctx, accountID, data.Bucket.ValueString(), data.Key.ValueString()); err != nil && !biznetgio.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete object %q: %s", data.Key.ValueString(), err))
	}
}

func (r *ObjectStorageObjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, ":", 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Expected `<account_id>:<bucket>:<key>`")
		return
	}
