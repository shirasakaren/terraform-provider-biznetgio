# Subscribed Object Storage instance. Polls until Active before proceeding.
resource "biznetgio_object_storage" "example" {
  product_id           = "8"
  cycle                = "m"
  label                = "data-archive"
  quota                = 10
  pay_with_credit_card = true
}
# wip 878
