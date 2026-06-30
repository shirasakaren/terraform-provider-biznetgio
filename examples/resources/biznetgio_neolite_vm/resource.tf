data "biznetgio_neolite_products" "all" {}

data "biznetgio_neolite_os_list" "ubuntu" {
  product_id = data.biznetgio_neolite_products.all.products[0].product_id
}

resource "biznetgio_neolite_keypair" "main" {
  name = "neo-lite-key"
}

# Order VM, tunggu sampai status Active.
resource "biznetgio_neolite_vm" "main" {
  ssh_and_console_user = "admin"
  console_password     = "s3cretP4ssw0rd"
  vm_name              = "neo-lite-1"
  product_id           = data.biznetgio_neolite_products.all.products[0].product_id
  select_os            = data.biznetgio_neolite_os_list.ubuntu.oss[0].name
  keypair_id           = biznetgio_neolite_keypair.main.keypair_id
  cycle                = "m"
}

# Upgrade disk (absolute target GB — cuma bisa naik):
# disk_size = 100

# Ganti power: power_state = "stop" (start/stop/suspend/resume/shutdown)

# Rebuild OS: rebuild_os = data.biznetgio_neolite_os_list.ubuntu.oss[1].name

# Migrate ke NEO Lite Pro: migrate_to_pro = "12345"  # neolitepro_product_id
# wip 841
