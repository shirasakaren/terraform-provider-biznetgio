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
