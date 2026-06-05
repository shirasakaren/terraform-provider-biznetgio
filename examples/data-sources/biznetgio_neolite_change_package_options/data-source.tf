# Ganti <account_id> dengan account id VM NEO Lite.
data "biznetgio_neolite_change_package_options" "this" {
  account_id = 12345
}

output "options_raw" {
  value = data.biznetgio_neolite_change_package_options.this.raw
}
