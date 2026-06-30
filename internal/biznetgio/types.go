package biznetgio

import (
	"bytes"
	"encoding/json"
	"strconv"
)

type envelope struct {
	Success bool            `json:"success"`
	Code    int             `json:"code"`
	Data    json.RawMessage `json:"data"`
}

type BillingResource struct {
	OrderID   string `json:"order_id"`
	AccountID string `json:"account_id"`
}

type AccountResource struct {
	AccountID       string       `json:"account_id"` // string on the wire, may be numeric
	Domain          string       `json:"domain"`
	Status          string       `json:"status"` // Active|Pending|Suspended|Terminated
	Billingcycle    string       `json:"billingcycle"`
	DateCreated     string       `json:"date_created"`
	NextDue         string       `json:"next_due"`
	RecurringAmount int64        `json:"recurring_amount"`
	ExtraDetails    ExtraDetails `json:"extra_details"`
	ProductID       int64        `json:"product_id"`
	ProductName     string       `json:"product_name"`
	Description     string       `json:"description"`
	CategoryID      int64        `json:"category_id"`
	CategoryName    string       `json:"category_name"`
	LastInvoice     LastInvoice  `json:"last_invoice"`
}

type LastInvoice struct {
	ID          int64  `json:"id"`
	PaidID      int64  `json:"paid_id"`
	Status      string `json:"status"`
	Date        string `json:"date"`
	Duedate     string `json:"duedate"`
	Paybefore   string `json:"paybefore"`
	Datepaid    string `json:"datepaid"`
	InvoiceType string `json:"invoice_type"`
// wip 785
// wip 816
