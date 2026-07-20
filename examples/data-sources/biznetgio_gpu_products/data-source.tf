provider "biznetgio" {}

data "biznetgio_gpu_products" "all" {}

output "gpu_products" {
  value = data.biznetgio_gpu_products.all.products
}
