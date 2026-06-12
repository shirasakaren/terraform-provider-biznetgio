data "biznetgio_baremetal_rebuild_os_list" "main" {
  account_id = 12345
}

output "rebuild_oss" {
  value = data.biznetgio_baremetal_rebuild_os_list.main.oss
}
# wip 590
# wip 665
