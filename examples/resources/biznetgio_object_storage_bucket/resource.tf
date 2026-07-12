resource "biznetgio_object_storage_bucket" "example" {
  account_id = biznetgio_object_storage.example.id
  name       = "my-app-assets"
  acl        = "private"
}
# wip 875
# wip 904
