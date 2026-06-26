# Cek dulu IP available sebelum order product.
data "biznetgio_neolite_pro_ip_availability" "this" {
  product_id = 30
}

output "ip_available" {
  value = data.biznetgio_neolite_pro_ip_availability.this.available
}
