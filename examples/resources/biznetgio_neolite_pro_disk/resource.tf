# Disk tambahan NEO Lite Pro (neolite_account_id = id dari biznetgio_neolite_pro_vm).
resource "biznetgio_neolite_pro_disk" "main" {
  product_id         = 30
  cycle              = "m"
  neolite_account_id = biznetgio_neolite_pro_vm.main.id
  service_name       = "extra-data"
  size               = 30
}

# Upgrade disk: naikin size (increment otomatis dihitung dari selisih).
