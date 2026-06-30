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
