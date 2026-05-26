package biznetgio

import (
	"context"
	"fmt"
	"net/url"
)

type GPUService struct {
	client *Client
}

type GPUCreateRequest struct {
	ProductID         int64  `json:"product_id"`
	SelectOS          string `json:"select_os"`
	KeypairID         int64  `json:"keypair_id"`
	ServiceName       string `json:"service_name,omitempty"`
	SSHAndConsoleUser string `json:"ssh_and_console_user"`
	ConsolePassword   string `json:"console_password"`
	Promocode         string `json:"promocode,omitempty"`
	PayInvoiceWithCC  string `json:"pay_invoice_with_cc,omitempty"`
	Cycle             string `json:"cycle"`
}

type GPUOneTimeCreateRequest struct {
	ProductID         int64  `json:"product_id"`
	SelectOS          string `json:"select_os"`
	KeypairID         int64  `json:"keypair_id"`
	ServiceName       string `json:"service_name,omitempty"`
	SSHAndConsoleUser string `json:"ssh_and_console_user"`
	ConsolePassword   string `json:"console_password"`
	Promocode         string `json:"promocode,omitempty"`
	PayInvoiceWithCC  string `json:"pay_invoice_with_cc,omitempty"`
	AdditionalHours   int64  `json:"additional_hours,omitempty"`
}

type SetVMStatusRequest struct {
	Status string `json:"status"`
}

type GPURebuildRequest struct {
	SelectOS string `json:"select_os"`
}

type ReserveAdditionalHoursRequest struct {
	Hours int64 `json:"hours,omitempty"`
}

func (s *GPUService) Create(ctx context.Context, req GPUCreateRequest) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "POST", "/neo-gpus", req, &out)
	return out, err
}

func (s *GPUService) CreateOneTime(ctx context.Context, req GPUOneTimeCreateRequest) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "POST", "/neo-gpus/one-time", req, &out)
	return out, err
}

func (s *GPUService) AccountsList(ctx context.Context, status string) ([]map[string]any, error) {
	var out []map[string]any
	err := s.client.doJSON(ctx, "GET", "/neo-gpus/accounts"+statusQuery(status), nil, &out)
	return out, err
}

func (s *GPUService) AccountGet(ctx context.Context, accountID int64) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/neo-gpus/accounts/%d", accountID), nil, &out)
	return out, err
}

func (s *GPUService) Delete(ctx context.Context, accountID int64) error {
	return s.client.doJSON(ctx, "DELETE", fmt.Sprintf("/neo-gpus/%d", accountID), nil, nil)
}

func (s *GPUService) VMStatusGet(ctx context.Context, accountID int64) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/neo-gpus/accounts/%d/vm-status", accountID), nil, &out)
	return out, err
}

func (s *GPUService) VMStatusSet(ctx context.Context, accountID int64, req SetVMStatusRequest) (map[string]any, error) {
	var out map[string]any
// wip 200
// wip 322
// wip 329
// wip 410
