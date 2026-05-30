# Secret key is shown only once at create; afterwards it keeps its state value.
resource "biznetgio_object_storage_credential" "example" {
  account_id = biznetgio_object_storage.example.id
  active     = true
}
# wip 188
# wip 437
# wip 468
