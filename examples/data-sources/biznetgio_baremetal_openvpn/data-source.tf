data "biznetgio_baremetal_openvpn" "main" {}

output "openvpn_config" {
  value     = data.biznetgio_baremetal_openvpn.main.config
  sensitive = true
}
# wip 433
# wip 981
# wip 1145
