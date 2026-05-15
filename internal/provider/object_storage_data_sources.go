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
	}
}

func (d *ObjectStorageInstancesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*biznetgio.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *biznetgio.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}
	d.client = client
}

func (d *ObjectStorageInstancesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ObjectStorageInstancesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	items, err := d.client.ObjectStorage().AccountsList(ctx, data.Status.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list object storage instances: %s", err))
		return
	}

	var models []ObjectStorageInstanceModel
	for _, it := range items {
		m := ObjectStorageInstanceModel{
			Label:  objDSString(it, "label", "name"),
			Status: objDSString(it, "status", "state"),
			Raw:    objRawString(it),
		}
		if id, ok := objMapInt64(it, "account_id", "id"); ok {
			m.ID = types.StringValue(strconv.FormatInt(id, 10))
		} else if idStr, ok := objMapString(it, "account_id", "id"); ok {
			m.ID = types.StringValue(idStr)
		}
		if v, ok := objMapInt64(it, "product_id"); ok {
			m.ProductID = types.Int64Value(v)
		}
		if v, ok := objMapInt64(it, "quota"); ok {
			m.Quota = types.Int64Value(v)
		}
		models = append(models, m)
	}

	data.Instances, resp.Diagnostics = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: ObjectStorageInstanceModel{}.AttributeTypes()}, models)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ---- buckets ----

type ObjectStorageBucketsDataSource struct {
	client *biznetgio.Client
}

type ObjectStorageBucketsDataSourceModel struct {
	AccountID types.String `tfsdk:"account_id"`
	Buckets   types.List   `tfsdk:"buckets"`
}

type ObjectStorageBucketModel struct {
	Name types.String `tfsdk:"name"`
	ACL  types.String `tfsdk:"acl"`
	Raw  types.String `tfsdk:"raw"`
}

func (m ObjectStorageBucketModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
		"acl":  types.StringType,
		"raw":  types.StringType,
	}
}

func NewObjectStorageBucketsDataSource() datasource.DataSource {
	return &ObjectStorageBucketsDataSource{}
}

func (d *ObjectStorageBucketsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "biznetgio_object_storage_buckets"
}

func (d *ObjectStorageBucketsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List buckets inside an Object Storage instance.",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Object Storage instance account id.",
			},
			"buckets": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of buckets.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{Computed: true, MarkdownDescription: "Bucket name."},
						"acl":  schema.StringAttribute{Computed: true, MarkdownDescription: "Canned ACL."},
						"raw":  schema.StringAttribute{Computed: true, Sensitive: true, MarkdownDescription: "Full JSON of the item."},
					},
				},
			},
		},
	}
}

func (d *ObjectStorageBucketsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*biznetgio.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *biznetgio.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}
	d.client = client
}

func (d *ObjectStorageBucketsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ObjectStorageBucketsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountID, err := objParseAccountID(data.AccountID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid account_id", err.Error())
		return
	}

	items, err := d.client.ObjectStorage().BucketsList(ctx, accountID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list buckets: %s", err))
		return
	}

	var models []ObjectStorageBucketModel
	for _, it := range items {
		m := ObjectStorageBucketModel{
			Name: objDSString(it, "name", "bucket_name"),
			Raw:  objRawString(it),
		}
		if v, ok := objMapString(it, "acl"); ok {
			m.ACL = types.StringValue(v)
		}
		models = append(models, m)
	}

	data.Buckets, resp.Diagnostics = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: ObjectStorageBucketModel{}.AttributeTypes()}, models)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ---- credentials ----

type ObjectStorageCredentialsDataSource struct {
	client *biznetgio.Client
}

type ObjectStorageCredentialsDataSourceModel struct {
	AccountID   types.String `tfsdk:"account_id"`
	Credentials types.List   `tfsdk:"credentials"`
}

type ObjectStorageCredentialModel struct {
	AccessKey types.String `tfsdk:"access_key"`
	Active    types.Bool   `tfsdk:"active"`
}

func (m ObjectStorageCredentialModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"access_key": types.StringType,
		"active":     types.BoolType,
	}
}

func NewObjectStorageCredentialsDataSource() datasource.DataSource {
	return &ObjectStorageCredentialsDataSource{}
}

func (d *ObjectStorageCredentialsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "biznetgio_object_storage_credentials"
}

func (d *ObjectStorageCredentialsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List credentials of an Object Storage instance. Secret keys are never exposed.",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Object Storage instance account id.",
			},
			"credentials": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of credentials (access key + active state only; secret keys tidak pernah dikeluarkan).",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"access_key": schema.StringAttribute{Computed: true, Sensitive: true, MarkdownDescription: "Credential access key."},
						"active":     schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether the credential is enabled."},
					},
				},
			},
		},
	}
}

func (d *ObjectStorageCredentialsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*biznetgio.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *biznetgio.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}
	d.client = client
}

func (d *ObjectStorageCredentialsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ObjectStorageCredentialsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountID, err := objParseAccountID(data.AccountID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid account_id", err.Error())
		return
	}

	items, err := d.client.ObjectStorage().CredentialsList(ctx, accountID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list credentials: %s", err))
// wip 175
