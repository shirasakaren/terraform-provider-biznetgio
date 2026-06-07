// SPDX-License-Identifier: MPL-2.0

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

func TestAccGpuInstanceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { gpuAccPreCheck(t) },
		CheckDestroy:             testAccCheckGpuInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGpuInstanceSubscriptionConfig("acc-gpu-sub"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("biznetgio_gpu_instance.test", "service_name", "acc-gpu-sub"),
					resource.TestCheckResourceAttrSet("biznetgio_gpu_instance.test", "id"),
					resource.TestCheckResourceAttrSet("biznetgio_gpu_instance.test", "status"),
					resource.TestCheckResourceAttrSet("biznetgio_gpu_instance.test", "raw"),
				),
			},
			{
				ResourceName:            "biznetgio_gpu_instance.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"console_password", "ssh_and_console_user", "promocode", "pay_with_credit_card", "subscription", "on_demand", "rebuild_trigger", "reserve_additional_hours_trigger", "select_os", "keypair_id"},
			},
		},
	})
}

func TestAccGpuInstanceOnDemandResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { gpuAccPreCheck(t) },
		CheckDestroy:             testAccCheckGpuInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGpuInstanceOnDemandConfig("acc-gpu-ond"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("biznetgio_gpu_instance.test", "service_name", "acc-gpu-ond"),
					resource.TestCheckResourceAttrSet("biznetgio_gpu_instance.test", "id"),
				),
			},
			{
				ResourceName:            "biznetgio_gpu_instance.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"console_password", "ssh_and_console_user", "promocode", "pay_with_credit_card", "subscription", "on_demand", "rebuild_trigger", "reserve_additional_hours_trigger", "select_os", "keypair_id"},
			},
		},
	})
}

func gpuAccPreCheck(t *testing.T) {
	testAccPreCheck(t)
	for _, env := range []string{"BIZNETGIO_TEST_GPU_PRODUCT_ID", "BIZNETGIO_TEST_GPU_KEYPAIR_ID", "BIZNETGIO_TEST_GPU_OS"} {
		if os.Getenv(env) == "" {
			t.Fatalf("%s must be set for gpu acceptance tests", env)
		}
	}
}

func testAccGpuInstanceSubscriptionConfig(name string) string {
	return fmt.Sprintf(`
provider "biznetgio" {}

resource "biznetgio_gpu_instance" "test" {
  product_id           = %s
  select_os            = %q
  keypair_id           = %s
  service_name         = %q
  ssh_and_console_user = "root"
  console_password     = "AccTestPass123!"

  subscription {
    cycle = "m"
  }
}
`, os.Getenv("BIZNETGIO_TEST_GPU_PRODUCT_ID"), os.Getenv("BIZNETGIO_TEST_GPU_OS"), os.Getenv("BIZNETGIO_TEST_GPU_KEYPAIR_ID"), name)
}

func testAccGpuInstanceOnDemandConfig(name string) string {
	return fmt.Sprintf(`
provider "biznetgio" {}

resource "biznetgio_gpu_instance" "test" {
  product_id           = %s
  select_os            = %q
  keypair_id           = %s
  service_name         = %q
  ssh_and_console_user = "root"
  console_password     = "AccTestPass123!"

  on_demand {
    additional_hours = 1
  }
}
`, os.Getenv("BIZNETGIO_TEST_GPU_PRODUCT_ID"), os.Getenv("BIZNETGIO_TEST_GPU_OS"), os.Getenv("BIZNETGIO_TEST_GPU_KEYPAIR_ID"), name)
}

func testAccCheckGpuInstanceDestroy(s *terraform.State) error {
	c := biznetgio.New(os.Getenv("BIZNETGIO_BASE_URL"), os.Getenv("BIZNETGIO_API_KEY"), 30*time.Second)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "biznetgio_gpu_instance" {
			continue
		}
		accountID, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return err
		}
		items, err := c.GPU().AccountsList(context.Background(), "")
		if err != nil {
			return err
		}
		for _, it := range items {
			if id, ok := gpuInt64(it, "account_id", "id"); ok && id == accountID {
				return fmt.Errorf("gpu instance %d masih ada", accountID)
			}
		}
	}
	return nil
}
