provider "biznetgio" {}

resource "biznetgio_gpu_instance" "example" {
  product_id           = 123
  select_os            = "ubuntu-22"
  keypair_id           = 456
  service_name         = "neo-gpu-biznet"
  ssh_and_console_user = "root"
  console_password     = "ganti-password-ini"
  pay_with_credit_card = true

  subscription {
    cycle = "m"
  }

  # ubah value ini buat rebuild ulang (destructive) atau reserve jam tambahan
  rebuild_trigger                  = "init"
  reserve_additional_hours_trigger = ""
}
# wip 435
# wip 694
