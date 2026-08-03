# Keypair di-generate server; private key cuma keluar sekali di response create.
resource "biznetgio_baremetal_keypair" "main" {
  name = "neo-metal-key"
}

# Import keypair yang sudah ada:
# resource "biznetgio_baremetal_keypair" "imported" {
#   name       = "existing-key"
#   public_key = "ssh-rsa AAAA..."
# }
