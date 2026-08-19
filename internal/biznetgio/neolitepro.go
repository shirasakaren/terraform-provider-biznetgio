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
}

func (s *NeoliteProService) VMState(ctx context.Context, accountID int64, state string) error {
	return s.client.doJSON(ctx, "PUT", fmt.Sprintf("/neolite-pros/accounts/%d/vm-state/%s", accountID, esc(state)), nil, nil)
}

func (s *NeoliteProService) VMRebuild(ctx context.Context, accountID int64, osName string) error {
	return s.client.doJSON(ctx, "PUT", fmt.Sprintf("/neolite-pros/accounts/%d/rebuild", accountID), struct {
		SelectOS string `json:"select_os"`
	}{SelectOS: osName}, nil)
}

func (s *NeoliteProService) VMChangeName(ctx context.Context, accountID int64, name string) error {
	return s.client.doJSON(ctx, "PUT", fmt.Sprintf("/neolite-pros/accounts/%d/change-vm-name", accountID), struct {
		Name string `json:"name"`
	}{Name: name}, nil)
}

func (s *NeoliteProService) VMChangeKeypair(ctx context.Context, accountID, keypairID int64) error {
	return s.client.doJSON(ctx, "PUT", fmt.Sprintf("/neolite-pros/accounts/%d/keypair", accountID), struct {
		KeypairID int64 `json:"keypair_id"`
	}{KeypairID: keypairID}, nil)
}

func (s *NeoliteProService) VMChangePackage(ctx context.Context, accountID int64, req NeoliteChangePackageRequest) (*BillingResource, error) {
	var out BillingResource
	err := s.client.doJSON(ctx, "POST", fmt.Sprintf("/neolite-pros/accounts/%d/change-package", accountID), req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *NeoliteProService) VMChangeStorage(ctx context.Context, accountID int64, req NeoliteUpgradeStorageRequest) (*BillingResource, error) {
	var out BillingResource
	err := s.client.doJSON(ctx, "PUT", fmt.Sprintf("/neolite-pros/accounts/%d/storage", accountID), req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *NeoliteProService) ChangePackagePrepare(ctx context.Context, accountID int64) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/neolite-pros/accounts/%d/change-package", accountID), nil, &out)
	return out, err
}

func (s *NeoliteProService) StoragePrepare(ctx context.Context, accountID int64) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/neolite-pros/accounts/%d/storage", accountID), nil, &out)
	return out, err
}

func (s *NeoliteProService) SnapshotCreate(ctx context.Context, accountID int64, req NeoliteSnapshotRequest) (*BillingResource, error) {
	var out BillingResource
	err := s.client.doJSON(ctx, "POST", fmt.Sprintf("/neolite-pros/accounts/%d/snapshot", accountID), req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *NeoliteProService) DiskCreate(ctx context.Context, req NeoliteDiskCreateRequest) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "POST", "/neolite-pros/disks", req, &out)
	return out, err
}

func (s *NeoliteProService) DiskList(ctx context.Context, status string) ([]map[string]any, error) {
	var out []map[string]any
	err := s.client.doJSON(ctx, "GET", "/neolite-pros/disks/accounts"+statusQuery(status), nil, &out)
	return out, err
}

func (s *NeoliteProService) DiskGet(ctx context.Context, accountID int64) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/neolite-pros/disks/accounts/%d", accountID), nil, &out)
	return out, err
}

func (s *NeoliteProService) DiskUpgrade(ctx context.Context, accountID int64, req NeoliteDiskUpgradeRequest) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "PUT", fmt.Sprintf("/neolite-pros/disks/accounts/%d", accountID), req, &out)
	return out, err
}

func (s *NeoliteProService) DiskDelete(ctx context.Context, accountID int64) error {
	return s.client.doJSON(ctx, "DELETE", fmt.Sprintf("/neolite-pros/disks/%d", accountID), nil, nil)
}

func (s *NeoliteProService) DiskProducts(ctx context.Context) ([]map[string]any, error) {
	var out []map[string]any
	err := s.client.doJSON(ctx, "GET", "/neolite-pros/disks/products", nil, &out)
	return out, err
}

func (s *NeoliteProService) DiskProductGet(ctx context.Context, productID int64) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/neolite-pros/disks/products/%d", productID), nil, &out)
	return out, err
}

func (s *NeoliteProService) AccountSnapshotList(ctx context.Context, status string) ([]SnapshotAccountResource, error) {
	var out []SnapshotAccountResource
	err := s.client.doJSON(ctx, "GET", "/neolite-pros/snapshots/accounts"+statusQuery(status), nil, &out)
	return out, err
}

func (s *NeoliteProService) AccountSnapshotGet(ctx context.Context, accountID int64) (SnapshotAccountResource, error) {
	var out SnapshotAccountResource
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/neolite-pros/snapshots/accounts/%d", accountID), nil, &out)
	return out, err
}

func (s *NeoliteProService) SnapshotRestoreWith(ctx context.Context, accountID int64, req NeoliteFromSnapshotRequest) (*BillingResource, error) {
	var out BillingResource
	err := s.client.doJSON(ctx, "POST", fmt.Sprintf("/neolite-pros/snapshots/accounts/%d/create", accountID), req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *NeoliteProService) SnapshotRestore(ctx context.Context, accountID int64) error {
	return s.client.doJSON(ctx, "PUT", fmt.Sprintf("/neolite-pros/snapshots/accounts/%d/restore", accountID), nil, nil)
}

func (s *NeoliteProService) SnapshotDelete(ctx context.Context, accountID int64) error {
	return s.client.doJSON(ctx, "DELETE", fmt.Sprintf("/neolite-pros/snapshots/%d", accountID), nil, nil)
}

func (s *NeoliteProService) SnapshotProducts(ctx context.Context) ([]PlanResource, error) {
	var out []PlanResource
	err := s.client.doJSON(ctx, "GET", "/neolite-pros/snapshots/products", nil, &out)
	return out, err
}

func (s *NeoliteProService) SnapshotProductGet(ctx context.Context, productID int64) (PlanResource, error) {
	var out PlanResource
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/neolite-pros/snapshots/products/%d", productID), nil, &out)
	return out, err
}

func (s *NeoliteProService) KeypairList(ctx context.Context) ([]KeypairResource, error) {
	var out []KeypairResource
	err := s.client.doJSON(ctx, "GET", "/neolite-pros/keypairs/", nil, &out)
	return out, err
}

func (s *NeoliteProService) KeypairCreate(ctx context.Context, req KeypairCreateRequest) (*KeypairResource, error) {
	var out KeypairResource
	err := s.client.doJSON(ctx, "POST", "/neolite-pros/keypairs/", req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// KeypairCreateRaw create keypair, return raw response map - private_key cuma
// keluar di response ini dan field-nya undocumented, jadi jangan di-decode typed.
func (s *NeoliteProService) KeypairCreateRaw(ctx context.Context, req KeypairCreateRequest) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "POST", "/neolite-pros/keypairs/", req, &out)
	return out, err
}

func (s *NeoliteProService) KeypairImport(ctx context.Context, req KeypairImportRequest) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "POST", "/neolite-pros/keypairs/import", req, &out)
	return out, err
}

func (s *NeoliteProService) KeypairDelete(ctx context.Context, keypairID int64) error {
	return s.client.doJSON(ctx, "DELETE", fmt.Sprintf("/neolite-pros/keypairs/%d", keypairID), nil, nil)
}

func (s *NeoliteProService) ProductList(ctx context.Context) ([]PlanResource, error) {
	var out []PlanResource
	err := s.client.doJSON(ctx, "GET", "/neolite-pros/products", nil, &out)
	return out, err
}

func (s *NeoliteProService) ProductGet(ctx context.Context, productID int64) (PlanResource, error) {
	var out PlanResource
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/neolite-pros/products/%d", productID), nil, &out)
	return out, err
}

func (s *NeoliteProService) ProductOSList(ctx context.Context, productID int64) ([]OsResource, error) {
	var out []OsResource
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/neolite-pros/products/%d/oss", productID), nil, &out)
	return out, err
}

func (s *NeoliteProService) ProductIPAvailability(ctx context.Context, productID int64) (IPAvailability, error) {
	var out IPAvailability
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/neolite-pros/products/%d/ip-availability", productID), nil, &out)
	return out, err
}
