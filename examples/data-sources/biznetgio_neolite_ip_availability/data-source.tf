data "biznetgio_neolite_products" "all" {}

data "biznetgio_neolite_ip_availability" "this" {
  product_id = data.biznetgio_neolite_products.all.products[0].product_id
}

output "available" {
  value = data.biznetgio_neolite_ip_availability.this.available
}
