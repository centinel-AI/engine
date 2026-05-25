# Task reference

All tasks are defined in [`Taskfile.yml`](../Taskfile.yml). Run `task --list` for a live summary.

Project JSON is maintained in a **separate repository** (see [centinel-AI/data](https://github.com/centinel-AI/data)). At run time, point the engine at that checkout with `ENGINE_DATA_REPO_PATH` in `.env` (default `../data`).

---

## Setup

| Task | Description |
|------|-------------|
| `env:init` | Create `.env` from `.env.example`. No-op if `.env` already exists. |

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

Build engine Docker images. Run `generate:<cloud>` first if templates are missing.

| Task | Description |
|------|-------------|
| `build:<cloud>:terraform` | Build the Terraform image for the given cloud. |
| `build:<cloud>:opentofu` | Build the OpenTofu image for the given cloud. |
| `build:all` | Generate all providers, then build all ten images in parallel. |

`<cloud>` is one of `azure` · `aws` · `gcp` · `oci` · `ovh`.

---

## Workspace

Populate `./workspace/` with the selected cloud module — the same layout baked into Docker images. Lets you run Terraform or OpenTofu directly on the host without Docker.

```bash
task workspace:azure:terraform PROJECT=project-01
terraform -chdir=workspace init
terraform -chdir=workspace plan
```

| Task | Description |
|------|-------------|
| `workspace:<cloud>` | Shorthand — populates with Terraform engine. |
| `workspace:<cloud>:terraform` | Populate `workspace/` with Terraform as the engine binary. |
| `workspace:<cloud>:opentofu` | Populate `workspace/` with OpenTofu as the engine binary. |

`PROJECT=<name>` overrides `ENGINE_TF_VAR_project` from `.env`. Project JSON is copied from the external checkout configured in `ENGINE_DATA_REPO_PATH`.

---

## Run

Execute a built engine image via Docker Compose. Default command is `terraform plan` / `tofu plan`. Append `-- <command>` to run anything else.

Compose bind-mounts the external JSON checkout at `/app/workspace/data` using `ENGINE_DATA_REPO_PATH` from `.env`.

```bash
task run:aws                              # plan with the project from .env
task run:aws PROJECT=networking           # override project
task run:aws -- terraform apply
task run:gcp:opentofu -- tofu destroy -auto-approve
```

| Task | Description |
|------|-------------|
| `run:<cloud>` | Shorthand — runs the Terraform container. |
| `run:<cloud>:terraform` | Run `engine-<cloud>-terraform`. |
| `run:<cloud>:opentofu` | Run `engine-<cloud>-opentofu`. |

`PROJECT=<name>` overrides `ENGINE_TF_VAR_project` for the container run.

---

## Cleanup

| Task | Description |
|------|-------------|
| `clean` | Reset `providers/` and `workspace/` to `.gitkeep`; remove `tmp/`. Does not delete `.env`. |
