package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/biznetgio/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type BaremetalResource struct {
	client *biznetgio.Client
}

type BaremetalResourceModel struct {
	ID                types.String   `tfsdk:"id"`
	AccountID         types.Int64    `tfsdk:"account_id"`
	ProductID         types.Int64    `tfsdk:"product_id"`
	Cycle             types.String   `tfsdk:"cycle"`
	SelectOS          types.String   `tfsdk:"select_os"`
	KeypairID         types.Int64    `tfsdk:"keypair_id"`
	Label             types.String   `tfsdk:"label"`
	PublicIP          types.Int64    `tfsdk:"public_ip"`
	Promocode         types.String   `tfsdk:"promocode"`
	PayWithCreditCard types.Bool     `tfsdk:"pay_with_credit_card"`
	PowerState        types.String   `tfsdk:"power_state"`
	ResetTrigger      types.String   `tfsdk:"reset_trigger"`
	RebuildOS         types.String   `tfsdk:"rebuild_os"`
	Status            types.String   `tfsdk:"status"`
// wip 314
// wip 431
