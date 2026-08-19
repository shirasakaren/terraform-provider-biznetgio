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
	k.PublicKey = a.PublicKey
	return nil
}

type jsonInt64 int64

func (n *jsonInt64) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) > 0 && b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		*n = jsonInt64(v)
		return nil
	}
	var v int64
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	*n = jsonInt64(v)
	return nil
}

type PlanResource struct {
	ProductID    int64     `json:"product_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CategoryID   int64     `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Options      Options   `json:"options"`
	Billing      []Billing `json:"billing"`
}

type Options struct {
	Type           string `json:"type"`
	Cores          int64  `json:"cores"`
	Memory         int64  `json:"memory"`
	AllowDowngrade int64  `json:"allow_downgrade"`
}

type Billing struct {
	Label      string      `json:"label"`
	Cycle      string      `json:"cycle"`
	Price      int64       `json:"price"`
	Components []Component `json:"components"` // may be null → tolerate
}

type Component struct {
	Label  string  `json:"label"`
	Field  string  `json:"field"`
	Prices []Price `json:"prices"`
}

type Price struct {
	QtyMin int64 `json:"qty_min"`
	QtyMax int64 `json:"qty_max"`
	Price  int64 `json:"price"`
}

type OsResource struct {
	VMID   int64  `json:"vmid"`
	Node   string `json:"node"`
	Name   string `json:"name"`
	MaxMem int64  `json:"maxmem"`
	MaxCPU int64  `json:"maxcpu"`
}

type IPAvailability struct {
	Available bool `json:"available"`
}

type SnapshotAccountResource struct {
	AccountID    string               `json:"account_id"`
	Status       string               `json:"status"`
	ExtraDetails SnapshotExtraDetails `json:"extra_details"`
}

type SnapshotExtraDetails struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Region      string `json:"region"`
}

type SnapshotResource struct {
	ID   string `json:"id"` // = snapshot account_id
	Name string `json:"name"`
}

type KeypairCreateRequest struct {
	Name string `json:"name"`
}

type KeypairImportRequest struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key,omitempty"`
}

func mapInt64(v map[string]any, key string) (int64, bool) {
	x, ok := v[key]
	if !ok {
		return 0, false
	}
	switch n := x.(type) {
	case json.Number:
		i, err := n.Int64()
		if err != nil {
			return 0, false
		}
		return i, true
	case float64:
		return int64(n), true
	case string:
		i, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			return 0, false
		}
		return i, true
	}
	return 0, false
}

func mapString(v map[string]any, key string) (string, bool) {
	x, ok := v[key]
	if !ok {
		return "", false
	}
	s, ok := x.(string)
	return s, ok
}
