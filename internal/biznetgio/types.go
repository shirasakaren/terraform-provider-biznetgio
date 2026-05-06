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
}

type ExtraDetails struct {
	Region      string  `json:"region"`
	RegionLabel string  `json:"region_label"`
	Description string  `json:"description"`
	Name        string  `json:"name"`
	TenantID    *string `json:"tenant_id"`
	CIUser      string  `json:"ciuser"`
	CIPassword  string  `json:"cipassword"`
	KeypairID   int64   `json:"neosshkey_id"`
	SSHKeys     string  `json:"sshkeys"`
	OSName      string  `json:"osname"`
	DiskSize    string  `json:"disk_size"` // string on the wire
}

type VirtualMachineResource struct {
	VMID    int64  `json:"vmid"` // Proxmox-internal; resource ID is always AccountID
	Name    string `json:"name"`
	Status  string `json:"status"`
	Uptime  int64  `json:"uptime"`
	MaxDisk int64  `json:"maxdisk"`
	MaxMem  int64  `json:"maxmem"`
	Mem     int64  `json:"mem"`
	CPUs    int64  `json:"cpus"`
}

type KeypairResource struct {
	KeypairID int64  `json:"keypair_id"` // note: NOT `id` on the wire
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
	// no private key in the documented response — see resources-plan.md note
}

func (k *KeypairResource) UnmarshalJSON(b []byte) error {
	var a struct {
		KeypairID jsonInt64 `json:"keypair_id"`
		Name      string    `json:"name"`
		PublicKey string    `json:"public_key"`
	}
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	k.KeypairID = int64(a.KeypairID)
	k.Name = a.Name
// wip 44
