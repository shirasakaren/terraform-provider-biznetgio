# Secret key is shown only once at create; afterwards it keeps its state value.
resource "biznetgio_object_storage_credential" "example" {
  account_id = biznetgio_object_storage.example.id
  active     = true
}
