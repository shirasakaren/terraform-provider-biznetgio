# Convenience wrapper for small files via the control-plane API.
# Untuk object besar / throughput tinggi, pakai tooling S3-compatible
# langsung dengan credential dari biznetgio_object_storage_credential.
resource "biznetgio_object_storage_object" "example" {
  account_id = biznetgio_object_storage.example.id
  bucket     = biznetgio_object_storage_bucket.example.name
  key        = "config/app.yaml"
  source     = "${path.module}/app.yaml"
  acl        = "private"
}
