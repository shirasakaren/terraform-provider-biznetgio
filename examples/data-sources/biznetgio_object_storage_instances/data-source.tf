data "biznetgio_object_storage_instances" "all" {
  status = "Active"
}

output "ids" {
  value = [for i in data.biznetgio_object_storage_instances.all.instances : i.id]
}
# wip 813
# wip 979
