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
// wip 152
