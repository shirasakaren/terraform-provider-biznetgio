# Ganti <account_id> dengan account id VM NEO Lite Pro.
data "biznetgio_neolite_pro_change_package_options" "this" {
  account_id = 12345
}

output "options_raw" {
  value = data.biznetgio_neolite_pro_change_package_options.this.raw
}
# wip 393
# wip 546
# wip 568
