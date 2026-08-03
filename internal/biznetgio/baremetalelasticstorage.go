package biznetgio

import (
	"context"
	"fmt"
)

type BaremetalElasticStorageService struct {
	client *Client
}

type ElasticStorageCreateRequest struct {
	ProductID        int64  `json:"product_id"`
	Cycle            string `json:"cycle"`
	StorageName      string `json:"storage_name"`
	MetalAccountID   int64  `json:"metal_account_id"`
	Size             int64  `json:"size,omitempty"`
	Promocode        string `json:"promocode,omitempty"`
	PayInvoiceWithCC string `json:"pay_invoice_with_cc,omitempty"`
}

type UpgradeElasticStorageRequest struct {
	Size             int64  `json:"size"`
	PayInvoiceWithCC string `json:"pay_invoice_with_cc,omitempty"`
}

type ChangeElasticStoragePackageRequest struct {
	NewProductID     int64  `json:"new_product_id"`
	PayInvoiceWithCC string `json:"pay_invoice_with_cc,omitempty"`
}

func (s *BaremetalElasticStorageService) Create(ctx context.Context, req ElasticStorageCreateRequest) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "POST", "/baremetal-neo-elastic-storages", req, &out)
	return out, err
}

func (s *BaremetalElasticStorageService) List(ctx context.Context, status string) ([]map[string]any, error) {
	var out []map[string]any
	err := s.client.doJSON(ctx, "GET", "/baremetal-neo-elastic-storages"+statusQuery(status), nil, &out)
	return out, err
}
