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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.biznetgio_neolite_pro_products.this", "id"),
					resource.TestCheckResourceAttrSet("data.biznetgio_neolite_pro_products.this", "products.0.product_id"),
				),
			},
		},
	})
}

func TestAccNeoliteProVM_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckNeoliteProVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNeoliteProBaseConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("biznetgio_neolite_pro_vm.test", "id"),
					resource.TestCheckResourceAttrSet("biznetgio_neolite_pro_vm.test", "order_id"),
					resource.TestCheckResourceAttr("biznetgio_neolite_pro_vm.test", "vm_name", "tf-acc-neolite-pro-vm"),
					resource.TestCheckResourceAttr("biznetgio_neolite_pro_vm.test", "cycle", "m"),
					resource.TestCheckResourceAttr("biznetgio_neolite_pro_vm.test", "pay_with_credit_card", "true"),
					resource.TestCheckResourceAttrSet("biznetgio_neolite_pro_vm.test", "status"),
					resource.TestCheckResourceAttrSet("biznetgio_neolite_pro_vm.test", "osname"),
					resource.TestCheckResourceAttrSet("biznetgio_neolite_pro_vm.test", "cipassword"),
				),
			},
			{
				ResourceName:      "biznetgio_neolite_pro_vm.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// write-only / create-only: ga bisa di-refetch
					"console_password",
					"ssh_and_console_user",
					"select_os",
					"cycle",
					"pay_with_credit_card",
					"promocode",
					"power_state",
					"rebuild_os",
					"timeouts",
				},
			},
		},
	})
}

func TestAccNeoliteProKeypair_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckNeoliteProKeypairDestroy,
		Steps: []resource.TestStep{
			{
				Config: `resource "biznetgio_neolite_pro_keypair" "test" {
					name = "tf-acc-neolite-pro-key"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("biznetgio_neolite_pro_keypair.test", "id"),
					resource.TestCheckResourceAttrSet("biznetgio_neolite_pro_keypair.test", "keypair_id"),
					resource.TestCheckResourceAttrSet("biznetgio_neolite_pro_keypair.test", "public_key"),
					resource.TestCheckResourceAttrSet("biznetgio_neolite_pro_keypair.test", "private_key"),
				),
			},
			{
				ResourceName:      "biznetgio_neolite_pro_keypair.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// private key cuma muncul sekali di response create
					"private_key",
				},
			},
		},
	})
}

func TestAccNeoliteProSnapshot_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckNeoliteProSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNeoliteProBaseConfig() + `
resource "biznetgio_neolite_pro_snapshot" "test" {
  neolite_account_id = biznetgio_neolite_pro_vm.test.id
  name               = "tf-acc-snapshot"
  cycle              = "m"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("biznetgio_neolite_pro_snapshot.test", "id"),
					resource.TestCheckResourceAttr("biznetgio_neolite_pro_snapshot.test", "name", "tf-acc-snapshot"),
					resource.TestCheckResourceAttrSet("biznetgio_neolite_pro_snapshot.test", "status"),
				),
			},
			{
				ResourceName:      "biznetgio_neolite_pro_snapshot.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// create-only / ga ada di detail response
					"neolite_account_id",
					"cycle",
					"pay_with_credit_card",
					"promocode",
					"timeouts",
				},
			},
		},
	})
}

func TestAccNeoliteProDisk_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckNeoliteProDiskDestroy,
// wip 90
// wip 133
// wip 311
