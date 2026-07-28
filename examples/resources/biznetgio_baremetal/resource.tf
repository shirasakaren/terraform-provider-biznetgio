# Auth: set BIZNETGIO_API_KEY (dan optional BIZNETGIO_BASE_URL) di env var.

data "biznetgio_baremetal_products" "all" {}

resource "biznetgio_baremetal_keypair" "main" {
  name = "neo-metal-key"
}

resource "biznetgio_baremetal" "main" {
  product_id = data.biznetgio_baremetal_products.all.products[0].product_id
  cycle      = "m"
  select_os  = "ubuntu-22"
  keypair_id = biznetgio_baremetal_keypair.main.keypair_id
  label      = "neo-metal-main"
  public_ip  = 1

  # power_state = "on"   # ganti ke "off" buat stop
  # reset_trigger = ""   # ganti nilainya buat reboot sekali
  # rebuild_os = "centos7-base"  # ganti buat wipe & reinstall OS
}
# wip 1070
