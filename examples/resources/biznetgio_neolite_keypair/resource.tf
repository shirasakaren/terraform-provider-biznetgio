# Keypair di-generate server; private key cuma keluar sekali di response create.
resource "biznetgio_neolite_keypair" "main" {
  name = "neo-lite-key"
}

# Import keypair yang sudah ada:
# resource "biznetgio_neolite_keypair" "imported" {
#   name = "existing-key"
# }
# wip 749
