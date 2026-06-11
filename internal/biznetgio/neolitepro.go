package biznetgio

import (
	"context"
	"fmt"
)

type NeoliteProService struct {
	client *Client
}

func (s *NeoliteProService) VMCreate(ctx context.Context, req NeoliteCreateRequest) (*BillingResource, error) {
	var out BillingResource
	err := s.client.doJSON(ctx, "POST", "/neolite-pros", req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *NeoliteProService) AccountsList(ctx context.Context, status string) ([]AccountResource, error) {
	var out []AccountResource
	err := s.client.doJSON(ctx, "GET", "/neolite-pros/accounts"+statusQuery(status), nil, &out)
	return out, err
}

func (s *NeoliteProService) AccountGet(ctx context.Context, accountID int64) (AccountResource, error) {
	var out AccountResource
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/neolite-pros/accounts/%d", accountID), nil, &out)
	return out, err
}

func (s *NeoliteProService) VMDetails(ctx context.Context, accountID int64) (VirtualMachineResource, error) {
	var out VirtualMachineResource
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/neolite-pros/accounts/%d/vm-details", accountID), nil, &out)
	return out, err
}

func (s *NeoliteProService) VMDelete(ctx context.Context, accountID int64) error {
	return s.client.doJSON(ctx, "DELETE", fmt.Sprintf("/neolite-pros/%d", accountID), nil, nil)
