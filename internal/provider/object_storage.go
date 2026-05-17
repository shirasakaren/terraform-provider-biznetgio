package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/biznetgio/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ObjectStorageResource{}

type ObjectStorageResource struct {
	client *biznetgio.Client
}

type ObjectStorageResourceModel struct {
	ID                types.String   `tfsdk:"id"`
	ProductID         types.Int64    `tfsdk:"product_id"`
	Cycle             types.String   `tfsdk:"cycle"`
	Label             types.String   `tfsdk:"label"`
	Quota             types.Int64    `tfsdk:"quota"`
	Promocode         types.String   `tfsdk:"promocode"`
	PayWithCreditCard types.Bool     `tfsdk:"pay_with_credit_card"`
	OrderID           types.String   `tfsdk:"order_id"`
	Status            types.String   `tfsdk:"status"`
	Raw               types.String   `tfsdk:"raw"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
}

func NewObjectStorageResource() resource.Resource {
	return &ObjectStorageResource{}
}

func (r *ObjectStorageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "biznetgio_object_storage"
}

func (r *ObjectStorageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Object Storage subscribed instance (the S3-compatible tenant) on BiznetGIO. " +
			"Create orders the instance; quota can be grown via update. Label, cycle and product are immutable.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Account id of the object storage instance.",
			},
			"product_id": schema.Int64Attribute{
				Required:            true,
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
				MarkdownDescription: "Product/plan id from the object storage catalog.",
			},
			"cycle": schema.StringAttribute{
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators:          []validator.String{stringvalidator.OneOf("m", "a", "q", "s", "b", "t", "p4", "p5")},
				MarkdownDescription: "Billing cycle: m=monthly, q=quarterly, s=semi-annual, a=annual, b=biennial, t=triennial, p4/p5=4/5-year prepay.",
			},
			"label": schema.StringAttribute{
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators:          []validator.String{stringvalidator.LengthBetween(6, 16), stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9\-_]*$`), "label hanya boleh alphanumeric, dash dan underscore")},
				MarkdownDescription: "Instance label, 6-16 chars, `[a-zA-Z0-9-_]`. Immutable (no rename endpoint).",
			},
			"quota": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(10),
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
				Validators:          []validator.Int64{int64validator.AtLeast(10)},
				MarkdownDescription: "Total quota in GB, minimum 10. Only growable: update sends the delta as add_quota.",
			},
			"promocode": schema.StringAttribute{
				Optional:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				MarkdownDescription: "Promo/discount code applied at order time.",
			},
			"pay_with_credit_card": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Pay the invoice with the card on file. Set false to leave the invoice unpaid in the portal; the resource stays Pending until paid.",
			},
			"order_id": schema.StringAttribute{
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Order id from the create call.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Lifecycle status: Active, Pending, Suspended or Terminated.",
			},
			"raw": schema.StringAttribute{
				Sensitive:           true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Full JSON of the last read from the API, for anything not modeled yet.",
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
		},
	}
}

func (r *ObjectStorageResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ObjectStorageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ObjectStorageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	pay := "no"
	if data.PayWithCreditCard.ValueBool() {
		pay = "yes"
	}

	created, err := r.client.ObjectStorage().Create(ctx, biznetgio.ObjectStorageCreateRequest{
		ProductID:        data.ProductID.ValueInt64(),
		Cycle:            data.Cycle.ValueString(),
		Label:            data.Label.ValueString(),
		Quota:            data.Quota.ValueInt64(),
		Promocode:        data.Promocode.ValueString(),
		PayInvoiceWithCC: pay,
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create object storage: %s", err))
		return
	}

	accountID, ok := objMapInt64(created, "account_id", "id")
	if !ok {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Object storage create response missing account_id: %s", objRawString(created)))
		return
	}

	data.ID = types.StringValue(strconv.FormatInt(accountID, 10))
	if v, ok := objMapString(created, "order_id"); ok {
		data.OrderID = types.StringValue(v)
	}

	acc, err := biznetgio.WaitForStatus(ctx, 5*time.Second,
		func(ctx context.Context) (map[string]any, error) {
			return r.client.ObjectStorage().AccountGet(ctx, accountID)
		},
		func(m map[string]any) string {
			v, _ := objMapString(m, "status", "state")
			return v
		},
		[]string{"Active"}, []string{"Terminated", "Suspended", "error"},
	)
	if err != nil {
		resp.Diagnostics.AddError("Object Storage Provisioning Failed",
			fmt.Sprintf("Timed out waiting for object storage %d to become Active: %s", accountID, err))
		return
	}

	objSetAccountState(&data, acc)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ObjectStorageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ObjectStorageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountID, err := strconv.ParseInt(data.ID.ValueString(), 10, 64)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	acc, err := r.client.ObjectStorage().AccountGet(ctx, accountID)
	if err != nil {
		if biznetgio.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read object storage %d: %s", accountID, err))
		return
