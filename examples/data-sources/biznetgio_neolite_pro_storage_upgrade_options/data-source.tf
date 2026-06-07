# Ganti <account_id> dengan account id VM NEO Lite Pro.
data "biznetgio_neolite_pro_storage_upgrade_options" "this" {
  account_id = 12345
}

output "options_raw" {
  value = data.biznetgio_neolite_pro_storage_upgrade_options.this.raw
}
