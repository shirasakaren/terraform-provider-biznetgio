package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/biznetgio/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccNeoClient() *biznetgio.Client {
	return biznetgio.New(os.Getenv("BIZNETGIO_BASE_URL"), os.Getenv("BIZNETGIO_API_KEY"), 30*time.Second)
}

func testAccCheckNeoliteVMDestroy(s *terraform.State) error {
	client := testAccNeoClient()
	items, err := client.Neolite().AccountsList(context.Background(), "")
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "biznetgio_neolite_vm" && rs.Type != "biznetgio_neolite_vm_from_snapshot" {
			continue
		}
		for _, it := range items {
			if it.AccountID == rs.Primary.ID {
				return fmt.Errorf("neolite vm %s masih ada setelah destroy", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckNeoliteKeypairDestroy(s *terraform.State) error {
	client := testAccNeoClient()
	items, err := client.Neolite().KeypairList(context.Background())
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "biznetgio_neolite_keypair" {
			continue
		}
		for _, it := range items {
			if strconv.FormatInt(it.KeypairID, 10) == rs.Primary.ID {
				return fmt.Errorf("neolite keypair %s masih ada setelah destroy", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckNeoliteSnapshotDestroy(s *terraform.State) error {
	client := testAccNeoClient()
	items, err := client.Neolite().AccountSnapshotList(context.Background(), "")
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "biznetgio_neolite_snapshot" {
			continue
		}
		for _, it := range items {
			if it.AccountID == rs.Primary.ID {
				return fmt.Errorf("neolite snapshot %s masih ada setelah destroy", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckNeoliteDiskDestroy(s *terraform.State) error {
	client := testAccNeoClient()
	items, err := client.Neolite().DiskList(context.Background(), "")
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "biznetgio_neolite_disk" {
			continue
		}
		for _, it := range items {
			if strconv.FormatInt(aliasInt(it, "account_id", "id"), 10) == rs.Primary.ID {
				return fmt.Errorf("neolite disk %s masih ada setelah destroy", rs.Primary.ID)
			}
		}
	}
	return nil
}

// testAccNeoliteBaseConfig: products + os_list + keypair + vm, dipakai semua test neolite.
func testAccNeoliteBaseConfig() string {
	return `
data "biznetgio_neolite_products" "this" {}

data "biznetgio_neolite_os_list" "this" {
  product_id = data.biznetgio_neolite_products.this.products[0].product_id
}

resource "biznetgio_neolite_keypair" "test" {
  name = "tf-acc-neolite-key"
}

resource "biznetgio_neolite_vm" "test" {
  ssh_and_console_user = "tfaccuser"
  console_password     = "TfaccP4ssw0rd!"
  vm_name              = "tf-acc-neolite-vm"
  product_id           = data.biznetgio_neolite_products.this.products[0].product_id
  select_os            = data.biznetgio_neolite_os_list.this.oss[0].name
  keypair_id           = biznetgio_neolite_keypair.test.keypair_id
  cycle                = "m"
}
`
}

func TestAccNeoliteProductsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `data "biznetgio_neolite_products" "this" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.biznetgio_neolite_products.this", "id"),
					resource.TestCheckResourceAttrSet("data.biznetgio_neolite_products.this", "products.0.product_id"),
				),
			},
		},
	})
}

func TestAccNeoliteVM_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckNeoliteVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNeoliteBaseConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("biznetgio_neolite_vm.test", "id"),
					resource.TestCheckResourceAttrSet("biznetgio_neolite_vm.test", "order_id"),
					resource.TestCheckResourceAttr("biznetgio_neolite_vm.test", "vm_name", "tf-acc-neolite-vm"),
					resource.TestCheckResourceAttr("biznetgio_neolite_vm.test", "cycle", "m"),
					resource.TestCheckResourceAttr("biznetgio_neolite_vm.test", "pay_with_credit_card", "true"),
					resource.TestCheckResourceAttrSet("biznetgio_neolite_vm.test", "status"),
					resource.TestCheckResourceAttrSet("biznetgio_neolite_vm.test", "osname"),
					resource.TestCheckResourceAttrSet("biznetgio_neolite_vm.test", "cipassword"),
				),
			},
			{
				ResourceName:      "biznetgio_neolite_vm.test",
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
