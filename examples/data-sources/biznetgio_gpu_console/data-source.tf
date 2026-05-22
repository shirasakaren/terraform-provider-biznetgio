data "biznetgio_gpu_console" "main" {
  account_id = 12345
}

output "console_url" {
  value     = data.biznetgio_gpu_console.main.url
  sensitive = true
}
