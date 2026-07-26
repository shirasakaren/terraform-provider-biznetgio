package biznetgio

import (
	"context"
	"fmt"
)

type NeoliteService struct {
	client *Client
}

type NeoliteCreateRequest struct {
	ProductID         int64  `json:"product_id"`
	Cycle             string `json:"cycle"`
	SelectOS          string `json:"select_os"`
	KeypairID         int64  `json:"keypair_id"`
	VMName            string `json:"vm_name,omitempty"`
	Description       string `json:"description,omitempty"`
	SSHAndConsoleUser string `json:"ssh_and_console_user"`
	ConsolePassword   string `json:"console_password"`
	Promocode         string `json:"promocode,omitempty"`
	PayInvoiceWithCC  string `json:"pay_invoice_with_cc,omitempty"`
}

type NeoliteChangePackageRequest struct {
	NewProductID     int64  `json:"new_product_id"`
	PayInvoiceWithCC string `json:"pay_invoice_with_cc,omitempty"`
}

type NeoliteUpgradeStorageRequest struct {
	DiskSize         int64  `json:"disk_size"`
	PayInvoiceWithCC string `json:"pay_invoice_with_cc,omitempty"`
}

type NeoliteSnapshotRequest struct {
	Cycle            string `json:"cycle"`
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
	Promocode        string `json:"promocode,omitempty"`
	PayInvoiceWithCC string `json:"pay_invoice_with_cc,omitempty"`
}

type MigrateToProRequest struct {
	NeoliteProProductID int64  `json:"neolitepro_product_id"`
	PayInvoiceWithCC    string `json:"pay_invoice_with_cc,omitempty"`
// wip 934
// wip 1022
