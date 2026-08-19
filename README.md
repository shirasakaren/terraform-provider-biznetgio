# Terraform Provider for BiznetGIO

> **Unofficial** community provider by [Shirasaka Ren](https://shirasaka.ren),
> not affiliated with or endorsed by PT Biznet Gio Nusantara.

A [Terraform](https://www.terraform.io) provider for managing
[BiznetGIO](https://www.biznetgio.com) cloud infrastructure via the
[BiznetGIO Portal API](https://api.portal.biznetgio.com/v1/docs).

BiznetGIO is an Indonesian cloud provider offering NEO Metal (bare metal
servers), NEO Lite / NEO Lite Pro (virtual machines), NEO GPU (GPU instances),
and NEO Object Storage (S3-compatible storage).

Documentation: https://biznetgio.creations.ren

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- A BiznetGIO account with an API token generated from the
  [customer portal](https://portal.biznetgio.com)

## Authentication

The provider authenticates with the `x-token` header. Set the token via the
`api_key` provider attribute or the `BIZNETGIO_API_KEY` environment variable.

```hcl
terraform {
  required_providers {
    biznetgio = {
      source  = "registry.terraform.io/shirasakaren/biznetgio"
    }
  }
}

provider "biznetgio" {
  # api_key is also read from BIZNETGIO_API_KEY
  # base_url is also read from BIZNETGIO_BASE_URL (default: https://api.portal.biznetgio.com/v1)
}

# catalog lookup
data "biznetgio_neolite_products" "plans" {}

# a NEO Lite VM
resource "biznetgio_neolite_vm" "web" {
  vm_name         = "web-1"
  product_id      = data.biznetgio_neolite_products.plans.products[0].product_id
  select_os       = "Ubuntu 22.04"
  keypair_id      = biznetgio_neolite_keypair.deploy.id
  cycle           = "m"
  # defaults to true: the invoice is paid automatically with the stored card.
  # set false to keep the order pending until paid manually in the portal.
  pay_with_credit_card = true
}

resource "biznetgio_neolite_keypair" "deploy" {
  name = "deploy-key"
}
```

> **Billing note**: every create/upgrade call places a real order and may
> charge the credit card on file. Resources created with
> `pay_with_credit_card = false` stay `Pending` until the invoice is paid in
> the portal.

## Resources

| Resource | Description |
|---|---|
| `biznetgio_baremetal` | NEO Metal bare metal server (power state, OS rebuild, additional-IP/elastic-storage attach via sub-resources) |
| `biznetgio_baremetal_keypair` | SSH keypair for NEO Metal |
| `biznetgio_baremetal_additional_ip` | Additional IP address for NEO Metal |
| `biznetgio_baremetal_additional_ip_assignment` | Assign an additional IP to a bare metal server |
| `biznetgio_baremetal_elastic_storage` | NEO Elastic Storage volume attached to a bare metal server |
| `biznetgio_gpu_instance` | NEO GPU instance (subscription or on-demand billing) |
| `biznetgio_gpu_keypair` | SSH keypair for NEO GPU |
| `biznetgio_neolite_vm` | NEO Lite virtual machine |
| `biznetgio_neolite_keypair` | SSH keypair for NEO Lite |
| `biznetgio_neolite_snapshot` | NEO Lite VM snapshot |
| `biznetgio_neolite_vm_from_snapshot` | Restore a NEO Lite VM from a snapshot |
| `biznetgio_neolite_disk` | Additional disk for NEO Lite |
| `biznetgio_neolite_pro_vm` | NEO Lite Pro virtual machine |
| `biznetgio_neolite_pro_keypair` | SSH keypair for NEO Lite Pro |
| `biznetgio_neolite_pro_snapshot` | NEO Lite Pro VM snapshot |
| `biznetgio_neolite_pro_disk` | Additional disk for NEO Lite Pro |
| `biznetgio_object_storage` | NEO Object Storage subscription |
| `biznetgio_object_storage_bucket` | Bucket in a NEO Object Storage account |
| `biznetgio_object_storage_credential` | S3 credential (access/secret key) |
| `biznetgio_object_storage_object` | Thin object upload wrapper - use S3-compatible tooling for bulk data |

## Data Sources

- `biznetgio_baremetal_products` / `_rebuild_os_list` / `_openvpn`
- `biznetgio_gpu_products` / `_console` / `_graph`
- `biznetgio_neolite_products` / `_os_list` / `_change_package_options` / `_storage_upgrade_options` / `_ip_availability`
- `biznetgio_neolite_pro_products` / `_os_list` / `_change_package_options` / `_storage_upgrade_options` / `_ip_availability`
- `biznetgio_object_storage_instances` / `_buckets` / `_credentials`

Every resource also exposes a computed `raw` attribute with the full last-read
API payload (secrets redacted, marked sensitive) as an escape hatch for fields
not yet modeled.

## Notes on the BiznetGIO API

- The Portal API does not publish response schemas. Response handling is
  defensive (case-insensitive alias lookups) and was cross-checked against
  BiznetGIO's own SDKs and CLI. Report any field mismatch as an issue.
- Power actions (start/stop/suspend) are declarative via `power_state`; the
  API is only called when the value changes.
- One-shot actions (reset, rebuild, reserve GPU hours, migrate-to-pro) are
  trigger attributes: change the string value to re-fire.
- Some products (NEO Virtual Compute, NEO Kubernetes, NEO DNS, domains, web
  hosting, gio-private, gio-enterprise-cloud, gio-backup) are provisioned
  manually in the portal and have no public API - they are out of scope.

## Development

```sh
make generate   # regenerate docs/ from schema + examples/ (needs tools/ go generate)
make testacc    # acceptance tests; requires TF_ACC=1 + BIZNETGIO_API_KEY
```

Acceptance tests create real, billable resources. Prefer a sandbox account.

## Publishing

Publishing to the Terraform Registry requires the GitHub repo to be named
`terraform-provider-biznetgio`, an RSA GPG signing key registered with the
registry, and a `v*` tag release (see `.goreleaser.yml`). Replace the
`biznetgio` namespace in `main.go` if you publish under a different GitHub
organization.

## License

MPL-2.0
