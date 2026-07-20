data "biznetgio_baremetal_products" "all" {}

output "first_product_id" {
  value = data.biznetgio_baremetal_products.all.products[0].product_id
}
# wip 976
