package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/biznetgio/terraform-provider-biznetgio/internal/biznetgio"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccObjectStoragePreCheck(t *testing.T) {
	testAccPreCheck(t)
	if os.Getenv("BIZNETGIO_TEST_OS_PRODUCT_ID") == "" {
		t.Fatal("BIZNETGIO_TEST_OS_PRODUCT_ID harus set (angka) untuk acceptance test object storage")
	}
}

func testAccObjectStorageEnvVars() (int64, string, string) {
	productID, _ := strconv.ParseInt(os.Getenv("BIZNETGIO_TEST_OS_PRODUCT_ID"), 10, 64)
	cycle := os.Getenv("BIZNETGIO_TEST_OS_CYCLE")
	if cycle == "" {
		cycle = "m"
	}
	label := os.Getenv("BIZNETGIO_TEST_OS_LABEL")
	if label == "" {
		label = "tf-acc-os-test"
	}
	return productID, cycle, label
}

func TestAccObjectStorageResource(t *testing.T) {
	productID, cycle, label := testAccObjectStorageEnvVars()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccObjectStoragePreCheck(t) },
		CheckDestroy:             testAccCheckObjectStorageDestroy,
		Steps: []resource.TestStep{
			{
// wip 968
