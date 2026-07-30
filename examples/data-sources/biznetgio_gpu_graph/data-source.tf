data "biznetgio_gpu_graph" "main" {
  account_id = 12345
  timeframe  = "day"
}

output "graph" {
  value = data.biznetgio_gpu_graph.main.graph
}
# wip 1155
