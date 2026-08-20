<p align="center">
  <img src="./assets/logo.png" width="88" alt="BiznetGIO logo" />
</p>

<h1 align="center">Terraform Provider for BiznetGIO</h1>

<p align="center">
  Manage <a href="https://www.biznetgio.com">BiznetGIO</a> cloud infrastructure - bare metal, VMs, GPUs, object storage - with Terraform.
</p>

<p align="center">
  <a href="https://github.com/shirasakaren/terraform-provider-biznetgio/actions/workflows/test.yml"><img src="https://github.com/shirasakaren/terraform-provider-biznetgio/actions/workflows/test.yml/badge.svg" alt="CI"></a>
  <a href="https://registry.terraform.io/providers/shirasakaren/biznetgio/latest"><img src="https://img.shields.io/badge/terraform%20registry-shirasakaren%2Fbiznetgio-844fba?logo=terraform" alt="Terraform Registry"></a>
  <a href="https://github.com/shirasakaren/terraform-provider-biznetgio/releases/latest"><img src="https://img.shields.io/github/v/release/shirasakaren/terraform-provider-biznetgio" alt="Latest release"></a>
  <a href="./LICENSE"><img src="https://img.shields.io/github/license/shirasakaren/terraform-provider-biznetgio" alt="License"></a>
</p>

<p align="center">
  <a href="https://biznetgio.creations.ren"><img src="https://img.shields.io/badge/docs-biznetgio.creations.ren-008541?style=for-the-badge" alt="Documentation site"></a>
</p>

> **Unofficial** community provider maintained by [Shirasaka Ren](https://shirasaka.ren) - not affiliated with or
> endorsed by PT Biznet Gio Nusantara.

A [Terraform](https://developer.hashicorp.com/terraform) provider for BiznetGIO, an Indonesian cloud platform,
built on the [BiznetGIO Portal API](https://api.portal.biznetgio.com/v1/docs). Covers NEO Metal (bare metal),
NEO Lite / NEO Lite Pro (VMs), NEO GPU, and NEO Object Storage (S3-compatible).

📖 **Full docs, every resource, and step-by-step guides: [biznetgio.creations.ren](https://biznetgio.creations.ren).**
New to IaC? Start with [What is Infrastructure as Code?](https://biznetgio.creations.ren/what-is-iac).

## Install

The provider is published on the [Terraform Registry](https://registry.terraform.io/providers/shirasakaren/biznetgio/latest):

```hcl
terraform {
  required_providers {
    biznetgio = {
      source = "registry.terraform.io/shirasakaren/biznetgio"
    }
  }
}
```

Requires Terraform >= 1.0.

## Quickstart

Get an API token from the [portal](https://portal.biznetgio.com) and export it:

```bash
export BIZNETGIO_API_KEY="<your-token>"
```

Then order a NEO Lite VM:

```hcl
data "biznetgio_neolite_products" "plans" {}

resource "biznetgio_neolite_keypair" "deploy" {
  name = "deploy-key"
}

resource "biznetgio_neolite_vm" "web" {
  vm_name                = "web-1"
  product_id             = data.biznetgio_neolite_products.plans.products[0].product_id
  select_os              = "Ubuntu 22.04"
  keypair_id             = biznetgio_neolite_keypair.deploy.id
  ssh_and_console_user   = "root"
  console_password       = "change-this-now!"
  cycle                  = "m"
  pay_with_credit_card   = true # bills the card on file, see below
}
```

> **Billing note**: `pay_with_credit_card` defaults to `true`, so the first `terraform apply` places a real
> order and may charge the card on file. Set it to `false` to leave the order pending until you pay in the portal.

## What's covered

<details>
<summary><b>Resources</b> (20 total, click to expand)</summary>

**NEO Metal**: `biznetgio_baremetal`, `biznetgio_baremetal_keypair`, `biznetgio_baremetal_additional_ip`,
`biznetgio_baremetal_additional_ip_assignment`, `biznetgio_baremetal_elastic_storage`.

**NEO Lite**: `biznetgio_neolite_vm`, `biznetgio_neolite_keypair`, `biznetgio_neolite_snapshot`,
`biznetgio_neolite_vm_from_snapshot`, `biznetgio_neolite_disk`.

**NEO Lite Pro**: `biznetgio_neolitepro_vm`, `biznetgio_neolitepro_keypair`, `biznetgio_neolitepro_snapshot`,
`biznetgio_neolitepro_disk`.

**NEO GPU**: `biznetgio_gpu_instance`, `biznetgio_gpu_keypair`.

**NEO Object Storage**: `biznetgio_object_storage`, `biznetgio_object_storage_bucket`,
`biznetgio_object_storage_credential`, `biznetgio_object_storage_object`.

</details>

<details>
<summary><b>Data sources</b> (19 total, click to expand)</summary>

Catalog lookups for every service: `biznetgio_neolite_products`, `biznetgio_neolite_os_list`,
`biznetgio_neolite_change_package_options`, `biznetgio_neolite_storage_upgrade_options`,
`biznetgio_neolite_ip_availability`, the same five for `neolitepro_*`, `biznetgio_baremetal_products`,
`biznetgio_baremetal_rebuild_os_list`, `biznetgio_baremetal_openvpn`, `biznetgio_gpu_products`,
`biznetgio_gpu_console`, `biznetgio_gpu_graph`, `biznetgio_object_storage_instances`,
`biznetgio_object_storage_buckets`, `biznetgio_object_storage_credentials`.

</details>

## Why this provider

- **Everything the API can do**: 20 resources and 19 data sources across all five public API groups.
- **Billing-aware**: every create and upgrade is a real paid order; you control auto-payment per resource.
- **Defensive by design**: the API publishes no response schemas, so every resource also exposes a redacted
  `raw` output with the full last-read payload, for anything not modeled yet.
- **Pulumi twin**: prefer general purpose languages? There is a [matching Pulumi provider](https://github.com/shirasakaren/pulumi-biznetgio)
  with the same resource boundaries.

## Development

```sh
make testacc    # acceptance tests (needs BIZNETGIO_API_KEY and a real account)
golangci-lint run
```

See the [development guide](https://biznetgio.creations.ren/guides/development) for the full workflow.

## License

[MPL-2.0](LICENSE)


