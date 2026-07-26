resource "biznetgio_baremetal_additional_ip" "main" {
  product_id = 1
  cycle      = "m"
}

resource "biznetgio_baremetal" "main" {
  product_id = 1
  cycle      = "m"
  label      = "neo-metal-main"
  keypair_id = 1
  public_ip  = 1
}

resource "biznetgio_baremetal_additional_ip_assignment" "main" {
  additional_ip_id = biznetgio_baremetal_additional_ip.main.account_id
  metal_account_id = biznetgio_baremetal.main.account_id
}
# wip 567
# wip 1027
