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

// wip 424
// wip 512
// wip 666
