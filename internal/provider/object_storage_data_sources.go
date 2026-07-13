package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/biznetgio/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource = &ObjectStorageInstancesDataSource{}
	_ datasource.DataSource = &ObjectStorageBucketsDataSource{}
	_ datasource.DataSource = &ObjectStorageCredentialsDataSource{}
)

// ---- instances ----

type ObjectStorageInstancesDataSource struct {
	client *biznetgio.Client
}

type ObjectStorageInstancesDataSourceModel struct {
	Status    types.String `tfsdk:"status"`
	Instances types.List   `tfsdk:"instances"`
}

type ObjectStorageInstanceModel struct {
	ID        types.String `tfsdk:"id"`
	Label     types.String `tfsdk:"label"`
	Status    types.String `tfsdk:"status"`
	ProductID types.Int64  `tfsdk:"product_id"`
	Quota     types.Int64  `tfsdk:"quota"`
	Raw       types.String `tfsdk:"raw"`
}

func (m ObjectStorageInstanceModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":         types.StringType,
		"label":      types.StringType,
		"status":     types.StringType,
		"product_id": types.Int64Type,
		"quota":      types.Int64Type,
		"raw":        types.StringType,
	}
}

func NewObjectStorageInstancesDataSource() datasource.DataSource {
	return &ObjectStorageInstancesDataSource{}
}

func (d *ObjectStorageInstancesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "biznetgio_object_storage_instances"
}

func (d *ObjectStorageInstancesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List Object Storage instances on the account.",
		Attributes: map[string]schema.Attribute{
			"status": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Filter by status (Active, Pending, Suspended, Terminated).",
			},
			"instances": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of object storage instances.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":         schema.StringAttribute{Computed: true, MarkdownDescription: "Account id."},
						"label":      schema.StringAttribute{Computed: true, MarkdownDescription: "Instance label."},
						"status":     schema.StringAttribute{Computed: true, MarkdownDescription: "Lifecycle status."},
						"product_id": schema.Int64Attribute{Computed: true, MarkdownDescription: "Product/plan id."},
						"quota":      schema.Int64Attribute{Computed: true, MarkdownDescription: "Quota in GB."},
						"raw":        schema.StringAttribute{Computed: true, Sensitive: true, MarkdownDescription: "Full JSON of the item."},
					},
				},
			},
		},
