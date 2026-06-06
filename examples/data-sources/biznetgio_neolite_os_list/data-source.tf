data "biznetgio_neolite_products" "all" {}

data "biznetgio_neolite_os_list" "this" {
  product_id = data.biznetgio_neolite_products.all.products[0].product_id
}

output "os_names" {
  value = [for os in data.biznetgio_neolite_os_list.this.oss : os.name]
}
