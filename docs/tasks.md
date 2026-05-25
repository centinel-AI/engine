# Task reference

All tasks are defined in [`Taskfile.yml`](../Taskfile.yml). Run `task --list` for a live summary.

The engine Taskfile covers **codegen** (provider templates) and **Docker image builds** only. Version pins come from `.env.example` / `.env`; the first `generate:*` or `build:*` run creates `.env` if missing.

---

## Codegen

Generates `providers/<cloud>/*.tf` from the pinned provider schema. Results are cached under `tmp/codegen/<cloud>/`.

| Task | Description |
|------|-------------|
| `generate:azure` | Regenerate `providers/azure/` from AzureRM + AzureAD schemas. |
| `generate:aws` | Regenerate `providers/aws/` from the AWS provider schema. |
| `generate:gcp` | Regenerate `providers/gcp/` from the Google provider schema. |
| `generate:oci` | Regenerate `providers/oci/` from the OCI provider schema. |
| `generate:ovh` | Regenerate `providers/ovh/` from the OVH provider schema. |
| `generate:all` | All five clouds in parallel. |

---

## Build

Build engine Docker images via Compose. Run `generate:<cloud>` first if `providers/<cloud>/` is empty.

| Task | Description |
|------|-------------|
| `build:<cloud>:terraform` | Build the Terraform image for the given cloud. |
| `build:<cloud>:opentofu` | Build the OpenTofu image for the given cloud. |
| `build:all` | `generate:all`, then build all ten images in parallel. |

`<cloud>` is one of `azure` · `aws` · `gcp` · `oci` · `ovh`.

Image names: `grauss-<cloud>-terraform` · `grauss-<cloud>-opentofu`.

---

← [Documentation index](../README.md#documentation)
