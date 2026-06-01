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
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"regexp"
)

type NeoliteVMResource struct {
	client *biznetgio.Client
}

type NeoliteVMResourceModel struct {
	ID                types.String            `tfsdk:"id"`
	OrderID           types.String            `tfsdk:"order_id"`
	SSHAndConsoleUser types.String            `tfsdk:"ssh_and_console_user"`
	ConsolePassword   types.String            `tfsdk:"console_password"`
	VMName            types.String            `tfsdk:"vm_name"`
	Description       types.String            `tfsdk:"description"`
	ProductID         types.Int64             `tfsdk:"product_id"`
	SelectOS          types.String            `tfsdk:"select_os"`
	KeypairID         types.Int64             `tfsdk:"keypair_id"`
	Cycle             types.String            `tfsdk:"cycle"`
	PayWithCreditCard types.Bool              `tfsdk:"pay_with_credit_card"`
	Promocode         types.String            `tfsdk:"promocode"`
	PowerState        types.String            `tfsdk:"power_state"`
	RebuildOS         types.String            `tfsdk:"rebuild_os"`
	MigrateToPro      types.String            `tfsdk:"migrate_to_pro"`
	DiskSize          types.Int64             `tfsdk:"disk_size"`
	Status            types.String            `tfsdk:"status"`
	Uptime            types.Int64             `tfsdk:"uptime"`
	MaxDisk           types.Int64             `tfsdk:"maxdisk"`
	MaxMem            types.Int64             `tfsdk:"maxmem"`
	Mem               types.Int64             `tfsdk:"mem"`
	CPUs              types.Int64             `tfsdk:"cpus"`
	CIUser            types.String            `tfsdk:"ciuser"`
	CIPassword        types.String            `tfsdk:"cipassword"`
	OSName            types.String            `tfsdk:"osname"`
	Region            types.String            `tfsdk:"region"`
	RegionLabel       types.String            `tfsdk:"region_label"`
	NextDue           types.String            `tfsdk:"next_due"`
	RecurringAmount   types.Int64             `tfsdk:"recurring_amount"`
	BillingCycle      types.String            `tfsdk:"billingcycle"`
	ProductName       types.String            `tfsdk:"product_name"`
	LastInvoice       NeoliteLastInvoiceModel `tfsdk:"last_invoice"`
// wip 492
