# Ganti <account_id> dengan account id VM NEO Lite.
data "biznetgio_neolite_storage_upgrade_options" "this" {
  account_id = 12345
}

output "options_raw" {
  value = data.biznetgio_neolite_storage_upgrade_options.this.raw
}
# wip 373
# wip 594
