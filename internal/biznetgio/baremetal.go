package biznetgio

import (
	"context"
	"fmt"
)

type BaremetalService struct {
	client *Client
}

type BaremetalCreateRequest struct {
	ProductID        int64  `json:"product_id"`
	Cycle            string `json:"cycle"`
	SelectOS         string `json:"select_os,omitempty"`
	KeypairID        int64  `json:"keypair_id"`
	Label            string `json:"label"`
	PublicIP         int64  `json:"public_ip"`
	Promocode        string `json:"promocode,omitempty"`
	PayInvoiceWithCC string `json:"pay_invoice_with_cc,omitempty"`
}

type BaremetalUpdateLabelRequest struct {
	Label string `json:"label"`
}

type BaremetalRebuildRequest struct {
	OS string `json:"os"`
}

func (s *BaremetalService) Create(ctx context.Context, req BaremetalCreateRequest) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "POST", "/baremetals", req, &out)
	return out, err
}

func (s *BaremetalService) AccountsList(ctx context.Context, status string) ([]map[string]any, error) {
	var out []map[string]any
	err := s.client.doJSON(ctx, "GET", "/baremetals/accounts"+statusQuery(status), nil, &out)
	return out, err
}

func (s *BaremetalService) AccountGet(ctx context.Context, accountID int64) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/baremetals/accounts/%d", accountID), nil, &out)
	return out, err
}

func (s *BaremetalService) StateGet(ctx context.Context, accountID int64) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/baremetals/accounts/%d/state", accountID), nil, &out)
	return out, err
}

func (s *BaremetalService) StateSet(ctx context.Context, accountID int64, state string) error {
	return s.client.doJSON(ctx, "PUT", fmt.Sprintf("/baremetals/accounts/%d/state/%s", accountID, esc(state)), nil, nil)
}

func (s *BaremetalService) KeypairList(ctx context.Context) ([]map[string]any, error) {
	var out []map[string]any
	err := s.client.doJSON(ctx, "GET", "/baremetals/keypairs/", nil, &out)
	return out, err
}

func (s *BaremetalService) KeypairCreate(ctx context.Context, req KeypairCreateRequest) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "POST", "/baremetals/keypairs/", req, &out)
	return out, err
}

func (s *BaremetalService) KeypairImport(ctx context.Context, req KeypairImportRequest) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "POST", "/baremetals/keypairs/import", req, &out)
	return out, err
}

func (s *BaremetalService) KeypairDelete(ctx context.Context, keypairID int64) error {
	return s.client.doJSON(ctx, "DELETE", fmt.Sprintf("/baremetals/keypairs/%d", keypairID), nil, nil)
}

func (s *BaremetalService) OpenVPN(ctx context.Context) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "GET", "/baremetals/openvpn", nil, &out)
	return out, err
}

func (s *BaremetalService) ProductsList(ctx context.Context) ([]map[string]any, error) {
	var out []map[string]any
	err := s.client.doJSON(ctx, "GET", "/baremetals/products", nil, &out)
	return out, err
}

func (s *BaremetalService) ProductGet(ctx context.Context, productID int64) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/baremetals/products/%d", productID), nil, &out)
	return out, err
}

func (s *BaremetalService) ProductOSList(ctx context.Context, productID int64) ([]map[string]any, error) {
	var out []map[string]any
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/baremetals/products/%d/oss", productID), nil, &out)
	return out, err
}

func (s *BaremetalService) States(ctx context.Context) ([]string, error) {
	var out []string
	err := s.client.doJSON(ctx, "GET", "/baremetals/states", nil, &out)
	return out, err
}

func (s *BaremetalService) UpdateLabel(ctx context.Context, accountID int64, req BaremetalUpdateLabelRequest) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "PUT", fmt.Sprintf("/baremetals/%d", accountID), req, &out)
	return out, err
}

func (s *BaremetalService) Delete(ctx context.Context, accountID int64) error {
	return s.client.doJSON(ctx, "DELETE", fmt.Sprintf("/baremetals/%d", accountID), nil, nil)
}

func (s *BaremetalService) Rebuild(ctx context.Context, accountID int64, req BaremetalRebuildRequest) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "PUT", fmt.Sprintf("/baremetals/%d/rebuild", accountID), req, &out)
	return out, err
}

func (s *BaremetalService) RebuildOSList(ctx context.Context, accountID int64) ([]map[string]any, error) {
	var out []map[string]any
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/baremetals/%d/rebuild/oss", accountID), nil, &out)
	return out, err
}

func (s *BaremetalService) AccountID(v map[string]any) (int64, bool) {
	return mapInt64(v, "account_id")
}

func (s *BaremetalService) ProductID(v map[string]any) (int64, bool) {
	return mapInt64(v, "product_id")
}

func (s *BaremetalService) Status(v map[string]any) (string, bool) {
	return mapString(v, "status")
}
// wip 1180
