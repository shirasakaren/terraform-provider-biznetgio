data "biznetgio_neolite_pro_products" "all" {}

output "first_product_id" {
  value = data.biznetgio_neolite_pro_products.all.products[0].product_id
}
