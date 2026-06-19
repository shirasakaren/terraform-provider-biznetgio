package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/biznetgio/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ObjectStorageBucketResource{}

type ObjectStorageBucketResource struct {
	client *biznetgio.Client
}

type ObjectStorageBucketResourceModel struct {
	ID        types.String `tfsdk:"id"`
	AccountID types.String `tfsdk:"account_id"`
	Name      types.String `tfsdk:"name"`
	ACL       types.String `tfsdk:"acl"`
	Raw       types.String `tfsdk:"raw"`
}

func NewObjectStorageBucketResource() resource.Resource {
	return &ObjectStorageBucketResource{}
}

func (r *ObjectStorageBucketResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "biznetgio_object_storage_bucket"
}

func (r *ObjectStorageBucketResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Bucket inside an Object Storage instance. Only the ACL is mutable; the name is immutable.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Composite id `<account_id>:<bucket_name>`.",
			},
			"account_id": schema.StringAttribute{
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				MarkdownDescription: "Object Storage instance account id.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators:          []validator.String{stringvalidator.LengthAtLeast(3), stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9\-_]*$`), "nama bucket hanya boleh alphanumeric, dash dan underscore")},
				MarkdownDescription: "Bucket name, minimum 3 chars, `[a-zA-Z0-9-_]`. Immutable.",
			},
			"acl": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators:          []validator.String{stringvalidator.OneOf("", "private", "public-read", "public-read-write", "authenticated-read", "log-delivery-write")},
				MarkdownDescription: "S3-style canned ACL: empty (default private), private, public-read, public-read-write, authenticated-read, log-delivery-write.",
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

func (r *ObjectStorageBucketResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ObjectStorageBucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ObjectStorageBucketResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountID, err := objParseAccountID(data.AccountID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid account_id", err.Error())
		return
	}

	if _, err := r.client.ObjectStorage().BucketCreate(ctx, accountID, biznetgio.BucketCreateRequest{
		Name: data.Name.ValueString(),
		ACL:  data.ACL.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create bucket %q: %s", data.Name.ValueString(), err))
		return
	}

	data.ID = types.StringValue(data.AccountID.ValueString() + ":" + data.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ObjectStorageBucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ObjectStorageBucketResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountID, err := objParseAccountID(data.AccountID.ValueString())
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	item, ok, err := objFindBucket(ctx, r.client, accountID, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list buckets: %s", err))
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

func (r *ObjectStorageBucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ObjectStorageBucketResourceModel
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

	if _, err := r.client.ObjectStorage().BucketSetACL(ctx, accountID, plan.Name.ValueString(), biznetgio.SetACLRequest{
		ACL: plan.ACL.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to set ACL on bucket %q: %s", plan.Name.ValueString(), err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
