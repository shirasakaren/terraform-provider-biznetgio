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

func testAccClient() *biznetgio.Client {
	return biznetgio.New(os.Getenv("BIZNETGIO_BASE_URL"), os.Getenv("BIZNETGIO_API_KEY"), 30*time.Second)
}

func testAccCheckBaremetalDestroy(s *terraform.State) error {
	client := testAccClient()
	items, err := client.Baremetal().AccountsList(context.Background(), "")
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "biznetgio_baremetal" {
			continue
		}
		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return err
		}
		for _, it := range items {
			if aliasInt(it, "account_id", "id") == id {
				return fmt.Errorf("baremetal %d masih ada setelah destroy", id)
			}
		}
	}
	return nil
}

func testAccCheckBaremetalKeypairDestroy(s *terraform.State) error {
	client := testAccClient()
	items, err := client.Baremetal().KeypairList(context.Background())
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "biznetgio_baremetal_keypair" {
			continue
		}
		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return err
		}
		for _, it := range items {
			if aliasInt(it, "keypair_id", "id") == id {
				return fmt.Errorf("baremetal keypair %d masih ada setelah destroy", id)
			}
		}
	}
	return nil
}

func TestAccBaremetal_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckBaremetalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBaremetalConfig("tf-acc-baremetal"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("biznetgio_baremetal.test", "id"),
					resource.TestCheckResourceAttrSet("biznetgio_baremetal.test", "account_id"),
					resource.TestCheckResourceAttr("biznetgio_baremetal.test", "label", "tf-acc-baremetal"),
					resource.TestCheckResourceAttr("biznetgio_baremetal.test", "cycle", "m"),
					resource.TestCheckResourceAttr("biznetgio_baremetal.test", "pay_with_credit_card", "true"),
					resource.TestCheckResourceAttrSet("biznetgio_baremetal.test", "status"),
				),
			},
			{
				ResourceName:      "biznetgio_baremetal.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

func TestAccBaremetalKeypair_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckBaremetalKeypairDestroy,
		Steps: []resource.TestStep{
			{
				Config: `resource "biznetgio_baremetal_keypair" "test" {
					name = "tf-acc-test-key"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("biznetgio_baremetal_keypair.test", "id"),
					resource.TestCheckResourceAttrSet("biznetgio_baremetal_keypair.test", "keypair_id"),
					resource.TestCheckResourceAttrSet("biznetgio_baremetal_keypair.test", "public_key"),
				),
			},
			{
				ResourceName:      "biznetgio_baremetal_keypair.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// private key cuma muncul sekali di response create, gak bisa di-refresh
					"private_key",
				},
			},
		},
	})
}

func testAccBaremetalConfig(label string) string {
	return fmt.Sprintf(`
data "biznetgio_baremetal_products" "all" {}

resource "biznetgio_baremetal_keypair" "test" {
	name = "tf-acc-test-key"
}

resource "biznetgio_baremetal" "test" {
	product_id = data.biznetgio_baremetal_products.all.products[0].product_id
	cycle      = "m"
	label      = %q
	keypair_id = biznetgio_baremetal_keypair.test.keypair_id
	public_ip  = 1
}
`, label)
}
