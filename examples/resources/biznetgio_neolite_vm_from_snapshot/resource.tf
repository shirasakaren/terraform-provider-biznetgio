data "biznetgio_neolite_products" "all" {}

resource "biznetgio_neolite_keypair" "main" {
  name = "neo-lite-key"
}

resource "biznetgio_neolite_vm" "source" {
  ssh_and_console_user = "admin"
  console_password     = "s3cretP4ssw0rd"
  vm_name              = "neo-lite-source"
  product_id           = data.biznetgio_neolite_products.all.products[0].product_id
  select_os            = "ubuntu-22"
  keypair_id           = biznetgio_neolite_keypair.main.keypair_id
  cycle                = "m"
}

resource "biznetgio_neolite_snapshot" "main" {
  neolite_account_id = biznetgio_neolite_vm.source.id
  name               = "pre-update"
  cycle              = "m"
}

# VM baru yang di-restore dari snapshot. Delete resource = delete VM.
resource "biznetgio_neolite_vm_from_snapshot" "restored" {
  snapshot_id          = biznetgio_neolite_snapshot.main.id
  product_id           = data.biznetgio_neolite_products.all.products[0].product_id
  cycle                = "m"
  keypair_id           = biznetgio_neolite_keypair.main.keypair_id
  name                 = "neo-lite-restored"
  ssh_and_console_user = "admin"
  console_password     = "s3cretP4ssw0rd"
}
