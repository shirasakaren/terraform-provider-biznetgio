# Secret keys are never exposed by this data source.
data "biznetgio_object_storage_credentials" "example" {
  account_id = biznetgio_object_storage.example.id
}

output "access_keys" {
  value = [for c in data.biznetgio_object_storage_credentials.example.credentials : c.access_key]
}
