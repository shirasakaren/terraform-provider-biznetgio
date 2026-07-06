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
				Config: testAccObjectStorageInstanceConfig(productID, cycle, label, 10),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("biznetgio_object_storage.test", "quota", "10"),
					resource.TestCheckResourceAttr("biznetgio_object_storage.test", "status", "Active"),
				),
			},
			{
				ResourceName:      "biznetgio_object_storage.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccObjectStorageInstanceConfig(productID, cycle, label, 20),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("biznetgio_object_storage.test", "quota", "20"),
				),
			},
		},
	})
}

func TestAccObjectStorageBucketResource(t *testing.T) {
	productID, cycle, label := testAccObjectStorageEnvVars()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccObjectStoragePreCheck(t) },
		CheckDestroy:             testAccCheckObjectStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectStorageBucketConfig(productID, cycle, label, "private"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("biznetgio_object_storage_bucket.test", "acl", "private"),
				),
			},
			{
				ResourceName:      "biznetgio_object_storage_bucket.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccObjectStorageBucketConfig(productID, cycle, label, "public-read"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("biznetgio_object_storage_bucket.test", "acl", "public-read"),
				),
			},
		},
	})
}

func TestAccObjectStorageCredentialResource(t *testing.T) {
	productID, cycle, label := testAccObjectStorageEnvVars()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccObjectStoragePreCheck(t) },
		CheckDestroy:             testAccCheckObjectStorageCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectStorageCredentialConfig(productID, cycle, label),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("biznetgio_object_storage_credential.test", "access_key"),
					resource.TestCheckResourceAttrSet("biznetgio_object_storage_credential.test", "secret_key"),
				),
			},
			{
				ResourceName:            "biznetgio_object_storage_credential.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_key"},
			},
			{
				Config: testAccObjectStorageCredentialDisabledConfig(productID, cycle, label),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("biznetgio_object_storage_credential.test", "active", "false"),
				),
			},
		},
	})
}

func TestAccObjectStorageObjectResource(t *testing.T) {
	productID, cycle, label := testAccObjectStorageEnvVars()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccObjectStoragePreCheck(t) },
		CheckDestroy:             testAccCheckObjectStorageObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectStorageObjectConfig(productID, cycle, label),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("biznetgio_object_storage_object.test", "key", "test/hello.txt"),
				),
			},
			{
				ResourceName:            "biznetgio_object_storage_object.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"content"},
			},
		},
	})
}

func TestAccObjectStorageInstancesDataSource(t *testing.T) {
	productID, cycle, label := testAccObjectStorageEnvVars()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccObjectStoragePreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccObjectStorageInstancesDSConfig(productID, cycle, label),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.biznetgio_object_storage_instances.test", "instances.#"),
				),
			},
		},
	})
}

func TestAccObjectStorageBucketsDataSource(t *testing.T) {
	productID, cycle, label := testAccObjectStorageEnvVars()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccObjectStoragePreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccObjectStorageBucketsDSConfig(productID, cycle, label),
