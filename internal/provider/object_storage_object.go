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
