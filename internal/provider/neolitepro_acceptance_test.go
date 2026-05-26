package provider

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccCheckNeoliteProVMDestroy(s *terraform.State) error {
	client := testAccNeoClient()
	items, err := client.NeolitePro().AccountsList(context.Background(), "")
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "biznetgio_neolite_pro_vm" {
			continue
		}
		for _, it := range items {
			if it.AccountID == rs.Primary.ID {
				return fmt.Errorf("neolite pro vm %s masih ada setelah destroy", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckNeoliteProKeypairDestroy(s *terraform.State) error {
	client := testAccNeoClient()
	items, err := client.NeolitePro().KeypairList(context.Background())
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "biznetgio_neolite_pro_keypair" {
			continue
		}
		for _, it := range items {
			if strconv.FormatInt(it.KeypairID, 10) == rs.Primary.ID {
				return fmt.Errorf("neolite pro keypair %s masih ada setelah destroy", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckNeoliteProSnapshotDestroy(s *terraform.State) error {
	client := testAccNeoClient()
	items, err := client.NeolitePro().AccountSnapshotList(context.Background(), "")
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "biznetgio_neolite_pro_snapshot" {
			continue
		}
		for _, it := range items {
			if it.AccountID == rs.Primary.ID {
				return fmt.Errorf("neolite pro snapshot %s masih ada setelah destroy", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckNeoliteProDiskDestroy(s *terraform.State) error {
	client := testAccNeoClient()
	items, err := client.NeolitePro().DiskList(context.Background(), "")
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "biznetgio_neolite_pro_disk" {
			continue
		}
		for _, it := range items {
			if strconv.FormatInt(aliasInt(it, "account_id", "id"), 10) == rs.Primary.ID {
				return fmt.Errorf("neolite pro disk %s masih ada setelah destroy", rs.Primary.ID)
			}
		}
	}
	return nil
}

// testAccNeoliteProBaseConfig: products + os_list + keypair + vm, dipakai semua test neolite pro.
func testAccNeoliteProBaseConfig() string {
	return `
data "biznetgio_neolite_pro_products" "this" {}

data "biznetgio_neolite_pro_os_list" "this" {
  product_id = data.biznetgio_neolite_pro_products.this.products[0].product_id
}

resource "biznetgio_neolite_pro_keypair" "test" {
  name = "tf-acc-neolite-pro-key"
}

resource "biznetgio_neolite_pro_vm" "test" {
  ssh_and_console_user = "tfaccuser"
  console_password     = "TfaccP4ssw0rd!"
  vm_name              = "tf-acc-neolite-pro-vm"
  product_id           = data.biznetgio_neolite_pro_products.this.products[0].product_id
  select_os            = data.biznetgio_neolite_pro_os_list.this.oss[0].name
  keypair_id           = biznetgio_neolite_pro_keypair.test.keypair_id
  cycle                = "m"
}
`
}

func TestAccNeoliteProProductsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `data "biznetgio_neolite_pro_products" "this" {}`,
