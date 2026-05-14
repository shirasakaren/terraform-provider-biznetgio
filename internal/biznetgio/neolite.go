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
}

type NeoliteDiskCreateRequest struct {
	ProductID        int64  `json:"product_id"`
	Cycle            string `json:"cycle"`
	NeoliteAccountID int64  `json:"neolite_account_id"`
	ServiceName      string `json:"service_name,omitempty"`
	Promocode        string `json:"promocode,omitempty"`
	PayInvoiceWithCC string `json:"pay_invoice_with_cc,omitempty"`
	Size             int64  `json:"size,omitempty"`
}

type NeoliteDiskUpgradeRequest struct {
	AdditionalSize   int64  `json:"additional_size,omitempty"`
	PayInvoiceWithCC string `json:"pay_invoice_with_cc,omitempty"`
}

type NeoliteFromSnapshotRequest struct {
	ProductID         int64  `json:"product_id"`
	Cycle             string `json:"cycle"`
	KeypairID         int64  `json:"keypair_id"`
	Name              string `json:"name"`
	Description       string `json:"description,omitempty"`
	SSHAndConsoleUser string `json:"ssh_and_console_user"`
	ConsolePassword   string `json:"console_password"`
	Promocode         string `json:"promocode,omitempty"`
	PayInvoiceWithCC  string `json:"pay_invoice_with_cc,omitempty"`
}

func (s *NeoliteService) VMCreate(ctx context.Context, req NeoliteCreateRequest) (*BillingResource, error) {
	var out BillingResource
	err := s.client.doJSON(ctx, "POST", "/neolites", req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *NeoliteService) AccountsList(ctx context.Context, status string) ([]AccountResource, error) {
	var out []AccountResource
	err := s.client.doJSON(ctx, "GET", "/neolites/accounts"+statusQuery(status), nil, &out)
	return out, err
}

func (s *NeoliteService) AccountGet(ctx context.Context, accountID int64) (AccountResource, error) {
	var out AccountResource
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/neolites/accounts/%d", accountID), nil, &out)
	return out, err
}

func (s *NeoliteService) VMDetails(ctx context.Context, accountID int64) (VirtualMachineResource, error) {
	var out VirtualMachineResource
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/neolites/accounts/%d/vm-details", accountID), nil, &out)
	return out, err
}

func (s *NeoliteService) VMDelete(ctx context.Context, accountID int64) error {
	return s.client.doJSON(ctx, "DELETE", fmt.Sprintf("/neolites/%d", accountID), nil, nil)
}

func (s *NeoliteService) VMState(ctx context.Context, accountID int64, state string) error {
	return s.client.doJSON(ctx, "PUT", fmt.Sprintf("/neolites/accounts/%d/vm-state/%s", accountID, esc(state)), nil, nil)
}

func (s *NeoliteService) VMRebuild(ctx context.Context, accountID int64, osName string) error {
	return s.client.doJSON(ctx, "PUT", fmt.Sprintf("/neolites/accounts/%d/rebuild", accountID), struct {
		SelectOS string `json:"select_os"`
	}{SelectOS: osName}, nil)
}

func (s *NeoliteService) VMChangeName(ctx context.Context, accountID int64, name string) error {
	return s.client.doJSON(ctx, "PUT", fmt.Sprintf("/neolites/accounts/%d/change-vm-name", accountID), struct {
		Name string `json:"name"`
	}{Name: name}, nil)
}

func (s *NeoliteService) VMChangeKeypair(ctx context.Context, accountID, keypairID int64) error {
	return s.client.doJSON(ctx, "PUT", fmt.Sprintf("/neolites/accounts/%d/keypair", accountID), struct {
		KeypairID int64 `json:"keypair_id"`
	}{KeypairID: keypairID}, nil)
}

func (s *NeoliteService) VMChangePackage(ctx context.Context, accountID int64, req NeoliteChangePackageRequest) (*BillingResource, error) {
	var out BillingResource
	err := s.client.doJSON(ctx, "POST", fmt.Sprintf("/neolites/accounts/%d/change-package", accountID), req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
// wip 129
