# Keypair di-generate server; private_key cuma keluar sekali di state.
resource "biznetgio_neolite_pro_keypair" "main" {
  name = "neo-lite-pro-key"
}

# Simpen private key-nya, misal pake output sensitive:
output "private_key" {
  value     = biznetgio_neolite_pro_keypair.main.private_key
  sensitive = true
}
