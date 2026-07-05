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
