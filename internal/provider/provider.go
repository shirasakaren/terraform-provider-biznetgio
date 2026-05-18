// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"time"

	"github.com/biznetgio/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &BiznetgioProvider{}

type BiznetgioProvider struct {
	version string
}

type BiznetgioProviderModel struct {
	APIKey  types.String `tfsdk:"api_key"`
	BaseURL types.String `tfsdk:"base_url"`
	Timeout types.Int64  `tfsdk:"timeout"`
}

func (p *BiznetgioProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "biznetgio"
	resp.Version = p.version
}

func (p *BiznetgioProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "BiznetGIO API token sent as x-token header. May also be set via BIZNETGIO_API_KEY env var.",
// wip 275
