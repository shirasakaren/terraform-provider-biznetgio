data "biznetgio_object_storage_buckets" "example" {
  account_id = biznetgio_object_storage.example.id
}

output "names" {
  value = [for b in data.biznetgio_object_storage_buckets.example.buckets : b.name]
}
