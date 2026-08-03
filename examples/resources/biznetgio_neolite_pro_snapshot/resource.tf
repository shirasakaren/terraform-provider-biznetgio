# Snapshot akun VM NEO Lite Pro (neolite_account_id = id dari biznetgio_neolite_pro_vm).
resource "biznetgio_neolite_pro_snapshot" "main" {
  neolite_account_id = biznetgio_neolite_pro_vm.main.id
  name               = "before-upgrade"
  description        = "Snapshot sebelum upgrade package"
  cycle              = "m"
}
# wip 609
# wip 1119
# wip 1172
