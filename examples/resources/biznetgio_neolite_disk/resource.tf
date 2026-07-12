resource "biznetgio_neolite_keypair" "main" {
  name = "neo-lite-key"
}

resource "biznetgio_neolite_vm" "main" {
  ssh_and_console_user = "admin"
  console_password     = "s3cretP4ssw0rd"
  vm_name              = "neo-lite-1"
  product_id           = 123
  select_os            = "ubuntu-22"
  keypair_id           = biznetgio_neolite_keypair.main.keypair_id
  cycle                = "m"
}

# Disk tambahan; product_id dari endpoint /neolites/disks/products.
resource "biznetgio_neolite_disk" "main" {
  product_id         = 456
  cycle              = "m"
  neolite_account_id = biznetgio_neolite_vm.main.id
  service_name       = "data-disk"
  size               = 100
}

# Upgrade disk: naikin `size` (increment dihitung otomatis).
