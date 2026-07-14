// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/biznetgio/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.Resource                   = &GpuInstanceResource{}
	_ resource.ResourceWithImportState    = &GpuInstanceResource{}
	_ resource.ResourceWithValidateConfig = &GpuInstanceResource{}
)

type GpuInstanceResource struct {
	client *biznetgio.Client
}

type GpuInstanceResourceModel struct {
	ID                            types.String   `tfsdk:"id"`
	ProductID                     types.Int64    `tfsdk:"product_id"`
	SelectOS                      types.String   `tfsdk:"select_os"`
	KeypairID                     types.Int64    `tfsdk:"keypair_id"`
	ServiceName                   types.String   `tfsdk:"service_name"`
	SSHAndConsoleUser             types.String   `tfsdk:"ssh_and_console_user"`
	ConsolePassword               types.String   `tfsdk:"console_password"`
	Promocode                     types.String   `tfsdk:"promocode"`
	PayWithCreditCard             types.Bool     `tfsdk:"pay_with_credit_card"`
	Subscription                  types.Object   `tfsdk:"subscription"`
	OnDemand                      types.Object   `tfsdk:"on_demand"`
	RebuildTrigger                types.String   `tfsdk:"rebuild_trigger"`
	ReserveAdditionalHoursTrigger types.String   `tfsdk:"reserve_additional_hours_trigger"`
	Status                        types.String   `tfsdk:"status"`
	OrderID                       types.String   `tfsdk:"order_id"`
// wip 834
// wip 942
