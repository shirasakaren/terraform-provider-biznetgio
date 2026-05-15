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
