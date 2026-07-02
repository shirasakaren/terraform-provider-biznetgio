package biznetgio

import (
	"context"
	"fmt"
)

type BaremetalAdditionalIPService struct {
	client *Client
}

type AdditionalIPCreateRequest struct {
	ProductID        int64  `json:"product_id"`
	Cycle            string `json:"cycle"`
	Region           string `json:"region,omitempty"`
	Promocode        string `json:"promocode,omitempty"`
	PayInvoiceWithCC string `json:"pay_invoice_with_cc,omitempty"`
}

type AssignToMachineRequest struct {
	MetalAccountID int64 `json:"metal_account_id"`
}

func (s *BaremetalAdditionalIPService) Create(ctx context.Context, req AdditionalIPCreateRequest) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "POST", "/baremetal-additional-ips", req, &out)
	return out, err
}

func (s *BaremetalAdditionalIPService) List(ctx context.Context, status string) ([]map[string]any, error) {
	var out []map[string]any
	err := s.client.doJSON(ctx, "GET", "/baremetal-additional-ips"+statusQuery(status), nil, &out)
	return out, err
}

func (s *BaremetalAdditionalIPService) Get(ctx context.Context, accountID int64) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/baremetal-additional-ips/%d", accountID), nil, &out)
	return out, err
}

func (s *BaremetalAdditionalIPService) Delete(ctx context.Context, accountID int64) error {
	return s.client.doJSON(ctx, "DELETE", fmt.Sprintf("/baremetal-additional-ips/%d", accountID), nil, nil)
}

func (s *BaremetalAdditionalIPService) AssignsList(ctx context.Context, accountID int64) ([]map[string]any, error) {
	var out []map[string]any
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/baremetal-additional-ips/%d/assigns", accountID), nil, &out)
	return out, err
}

func (s *BaremetalAdditionalIPService) Assign(ctx context.Context, accountID int64, req AssignToMachineRequest) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "PUT", fmt.Sprintf("/baremetal-additional-ips/%d/assigns", accountID), req, &out)
	return out, err
}

func (s *BaremetalAdditionalIPService) AssignGet(ctx context.Context, accountID, metalAccountID int64) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/baremetal-additional-ips/%d/assigns/%d", accountID, metalAccountID), nil, &out)
	return out, err
}

func (s *BaremetalAdditionalIPService) Unassign(ctx context.Context, accountID, metalAccountID int64) error {
	return s.client.doJSON(ctx, "DELETE", fmt.Sprintf("/baremetal-additional-ips/%d/assigns/%d", accountID, metalAccountID), nil, nil)
}
