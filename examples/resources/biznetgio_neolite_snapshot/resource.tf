resource "biznetgio_neolite_keypair" "main" {
  name = "neo-lite-key"
}

resource "biznetgio_neolite_vm" "source" {
  ssh_and_console_user = "admin"
  console_password     = "s3cretP4ssw0rd"
  vm_name              = "neo-lite-source"
  product_id           = 123
  select_os            = "ubuntu-22"
  keypair_id           = biznetgio_neolite_keypair.main.keypair_id
  cycle                = "m"
}

# Snapshot akun VM; snapshot punya account id sendiri.
resource "biznetgio_neolite_snapshot" "main" {
  neolite_account_id = biznetgio_neolite_vm.source.id
  name               = "pre-update"
  description        = "backup sebelum upgrade package"
  cycle              = "m"
}
