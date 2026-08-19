data "biznetgio_neolite_pro_products" "all" {}

data "biznetgio_neolite_pro_os_list" "this" {
  product_id = data.biznetgio_neolite_pro_products.all.products[0].product_id
}

output "first_os" {
  value = data.biznetgio_neolite_pro_os_list.this.oss[0].name
}
