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
