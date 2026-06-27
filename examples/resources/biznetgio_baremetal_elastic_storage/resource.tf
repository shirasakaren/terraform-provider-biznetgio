resource "biznetgio_baremetal_elastic_storage" "main" {
  product_id       = 1
  cycle            = "m"
  storage_name     = "data-volume"
  metal_account_id = 12345
  size             = 100
}
