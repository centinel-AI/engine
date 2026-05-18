# Grauss

Terraform / OpenTofu engine packaged as Docker images — one per cloud × engine combination. Each image contains a complete set of resource templates generated from the pinned provider schema. You bring credentials and JSON data files.

## Images

|         | Terraform              | OpenTofu              |
|---------|------------------------|-----------------------|
| Azure   | `engine-azure-terraform` | `engine-azure-opentofu` |
| AWS     | `engine-aws-terraform`   | `engine-aws-opentofu`   |
| GCP     | `engine-gcp-terraform`   | `engine-gcp-opentofu`   |
| OCI     | `engine-oci-terraform`   | `engine-oci-opentofu`   |
| OVH     | `engine-ovh-terraform`   | `engine-ovh-opentofu`   |

Azure covers AzureRM + AzureAD as a single stack. OVH bundles ovh/ovh + hashicorp/aws (configured for OVH S3-compatible Object Storage) for state backend support.

## How it works

```
provider schema ──▶ scripts/codegen.py ──▶ providers/<cloud>/*.tf ──▶ Docker image
                                                                           │
                    data/<cloud>/<project>/*.json ─────────────────────────▶ plan / apply
```

Resource templates are generated from the provider schema and baked into the image. Your JSON files drive `for_each` — no `.tf` changes needed to add or remove resource instances.

```
providers/<cloud>/          generated .tf templates  (gitignored)
data/<cloud>/<project>/     one JSON file per resource instance
  └─ <resource_type>/
       └─ <name>.json
workspace/                  local mirror of the container module  (gitignored)
scripts/                    codegen and tooling
Dockerfile                  parameterized by CLOUD_PROVIDER and ENGINE
compose.yml                 eight pre-configured services
```

## Versions

| Component | Version | Registry |
|-----------|---------|----------|
| Terraform | 1.15.3 | [developer.hashicorp.com](https://developer.hashicorp.com/terraform) |
| OpenTofu | 1.11.7 | [opentofu.org](https://opentofu.org) |
| hashicorp/azurerm | 4.72.0 | [registry.terraform.io](https://registry.terraform.io/providers/hashicorp/azurerm/latest) |
| hashicorp/azuread | 3.8.0 | [registry.terraform.io](https://registry.terraform.io/providers/hashicorp/azuread/latest) |
| hashicorp/google | 7.32.0 | [registry.terraform.io](https://registry.terraform.io/providers/hashicorp/google/latest) |
| hashicorp/aws | 6.44.0 | [registry.terraform.io](https://registry.terraform.io/providers/hashicorp/aws/latest) |
| oracle/oci | 8.13.0 | [registry.terraform.io](https://registry.terraform.io/providers/oracle/oci/latest) |
| ovh/ovh | 2.13.1 | [registry.terraform.io](https://registry.terraform.io/providers/ovh/ovh/latest) |

Pins live in `.env.example`. Update a pin → `task generate:<cloud>` → rebuild image.

## Prerequisites

| Tool | Required for |
|------|--------------|
| [Docker](https://docs.docker.com/get-docker/) + Compose v2 | Building and running images |
| [Task](https://taskfile.dev/) | All `task` commands |
| [Terraform](https://developer.hashicorp.com/terraform/install) ≥ 1.5 | Codegen (`task generate:*`) |
| [OpenTofu](https://opentofu.org/docs/intro/install/) ≥ 1.8 | Local OpenTofu workspace (optional) |

## Quick start

```bash
# 1. Initialize .env (then edit credentials)
task env:init

# 2. Generate provider templates
task generate:azure    # or: aws | gcp | oci | ovh | all

# 3. Build images
task build:azure:terraform
task build:all         # all eight

# 4. Run
task run:azure                           # plan with the project in .env
task run:azure PROJECT=networking        # override project
task run:aws -- terraform apply          # explicit command
task run:gcp:opentofu -- tofu destroy
```

Without Docker — populate `workspace/` and run the engine binary directly:

```bash
task workspace:azure:terraform PROJECT=project-01
terraform -chdir=workspace init
terraform -chdir=workspace plan
```

## Adding resources

**New instance** — drop a JSON file under `data/<cloud>/<project>/<resource_type>/`. Picked up automatically via `for_each`, no template changes needed.

**New resource type** — run `task generate:<cloud>` to regenerate templates, then rebuild the image.

**Provider version bump** — update the pin in `.env.example`, run `task generate:<cloud>`, rebuild.

## Contributing

[CONTRIBUTING.md](CONTRIBUTING.md) · [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) · [SECURITY.md](SECURITY.md) · [LICENSE](LICENSE)

## Documentation

| | |
|--|--|
| [docs/projects.md](docs/projects.md) | `data/` layout, `_bootstrap`, the six reference architectures and their resources per cloud |
| [docs/tasks.md](docs/tasks.md) | Full `task` reference |
