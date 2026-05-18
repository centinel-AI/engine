# Projects

Projects are the unit of deployment in Grauss. Each project is a directory under `data/<cloud>/` that contains one JSON file per resource instance. The engine picks up all files at plan/apply time via `for_each` — no `.tf` changes are needed to add or remove resources.

## Directory layout

```
data/
├── aws/
│   ├── _bootstrap/          # Remote state backend setup (run once)
│   ├── networking/          # Reference architecture
│   ├── vm-simple/
│   ├── web-app/
│   ├── kubernetes/
│   ├── data-lake/
│   └── messaging/
├── azure/   (same structure)
├── gcp/     (same structure)
└── oci/     (same structure)
```

Inside each project, subdirectories map to Terraform resource types (after stripping the provider prefix):

```
data/aws/networking/
├── backend.tf.json          # Optional: remote state backend config
├── vpc/
│   └── vpc-grauss-networking.json
├── subnet/
│   ├── subnet-public-a-grauss-networking.json
│   └── subnet-private-a-grauss-networking.json
├── internet_gateway/
│   └── igw-grauss-networking.json
├── route_table/
│   └── rt-grauss-networking.json
└── security_group/
    └── sg-default-grauss-networking.json
```

The subdirectory name (`vpc`, `subnet`, …) is the key into `local.json_data` inside the generated template. Each JSON file becomes one resource instance keyed by its `name` field.

---

## `_bootstrap` — remote state setup

The `_bootstrap` project provisions the backend infrastructure for Terraform remote state. It must be applied **before any other project** and uses local state (no backend block).

| Cloud | Resources provisioned |
|-------|-----------------------|
| AWS   | `s3_bucket` (`grauss-tfstate-bootstrap`), `dynamodb_table` (`grauss-tfstate-locks`) |
| Azure | `resource_group`, `storage_account`, `storage_container` |
| GCP   | `storage_bucket` (`grauss-tfstate-bootstrap`, region EU) |
| OCI   | `identity_compartment` |

After applying `_bootstrap`, update `backend.tf.json` in each project to point to the created bucket/container.

```bash
# Apply bootstrap (local state, no remote backend needed)
docker compose run --rm -e TF_VAR_project=_bootstrap engine-aws-terraform terraform apply
```

---

## Reference architectures

Six architectures are provided, each implemented across all four clouds. They are independent and can be deployed in any order (except that some reference IDs produced by `networking` or `_bootstrap`).

### networking

Foundational network layer. Establishes the VPC/VNet, subnets, internet gateway, route tables, and security groups that other architectures build on top of.

| Cloud | Resources |
|-------|-----------|
| AWS | `vpc`, `subnet` ×2 (public-a, private-a), `internet_gateway`, `route_table`, `security_group` |
| Azure | `resource_group`, `virtual_network`, `subnet` ×2 (public, private), `network_security_group` |
| GCP | `compute_network`, `compute_subnetwork`, `compute_firewall` ×2 (allow-internal, allow-ssh) |
| OCI | `identity_compartment`, `core_vcn`, `core_internet_gateway`, `core_route_table`, `core_security_list`, `core_subnet` ×2 (public, private) |

CIDR block: `10.0.0.0/16`

---

### vm-simple

Minimal single-VM deployment with a public IP. Useful as a jump host or to verify end-to-end network connectivity. Uses the smallest available instance type on each cloud.

| Cloud | Resources |
|-------|-----------|
| AWS | `vpc`, `subnet`, `security_group` (SSH only), `instance` (t3.micro), `eip` |
| Azure | `resource_group`, `virtual_network`, `subnet`, `public_ip`, `network_interface`, `linux_virtual_machine` (Standard_B1s) |
| GCP | `compute_network`, `compute_subnetwork`, `compute_firewall` (SSH), `compute_address` (EXTERNAL), `compute_instance` (e2-micro) |
| OCI | `identity_compartment`, `core_vcn`, `core_internet_gateway`, `core_route_table`, `core_security_list`, `core_subnet`, `core_instance` (VM.Standard.E2.1.Micro) |

CIDR block: `10.1.0.0/16`

---

### web-app

Three-tier web application: load balancer → application VM → managed PostgreSQL. Demonstrates cross-tier security group rules and database subnet delegation.

| Cloud | Resources |
|-------|-----------|
| AWS | `vpc`, `subnet` ×3 (public-a, public-b, private-a), `internet_gateway`, `security_group` ×3 (alb, app, db), `alb`, `alb_target_group`, `alb_listener`, `instance` (t3.micro), `db_subnet_group`, `db_instance` (db.t3.micro, PostgreSQL 16) |
| Azure | `resource_group`, `virtual_network`, `subnet` ×2 (app, db with PostgreSQL delegation), `network_security_group`, `public_ip`, `network_interface`, `linux_virtual_machine` (Standard_B1s), `postgresql_flexible_server` (Standard_B1ms, PostgreSQL 16) |
| GCP | `compute_network`, `compute_subnetwork`, `compute_firewall` ×2 (http, ssh), `compute_instance` (e2-micro), `sql_database_instance` (db-f1-micro, POSTGRES_16) |
| OCI | `identity_compartment`, `core_vcn`, `core_internet_gateway`, `core_route_table`, `core_security_list` (ports 22/80/443), `core_subnet`, `core_instance` (VM.Standard.E2.1.Micro), `database_autonomous_database` (OLTP, free tier) |

CIDR block: `10.2.0.0/16`

---

### kubernetes

Managed Kubernetes cluster with a single node pool. One node, minimum viable configuration. Control plane is fully managed by the cloud provider.

| Cloud | Resources |
|-------|-----------|
| AWS | `vpc`, `subnet` ×2 (private-a, private-b — tagged for EKS), `internet_gateway`, `security_group`, `eks_cluster` (v1.31), `eks_node_group` (t3.small, 1 node) |
| Azure | `resource_group`, `kubernetes_cluster` (AKS v1.31, Standard_B2s node, 1 node, Free SKU tier, SystemAssigned identity) |
| GCP | `compute_network`, `compute_subnetwork`, `container_cluster` (GKE, `remove_default_node_pool = true`), `container_node_pool` (e2-small, 1 node) |
| OCI | `identity_compartment`, `core_vcn`, `core_internet_gateway`, `core_route_table`, `core_security_list` (port 6443), `core_subnet`, `containerengine_cluster` (OKE v1.31.1), `containerengine_node_pool` (VM.Standard.E2.1.Micro, 1 node) |

CIDR block: `10.3.0.0/16`

> **Cost warning:** AWS EKS and GCP GKE charge ~$0.10/hr for the control plane regardless of workloads. Destroy when not in use. Azure AKS Free tier and OCI OKE have no control-plane charge.

---

### data-lake

Object storage for raw and processed data layers plus a managed database or data warehouse. Demonstrates multi-bucket patterns and cloud-native analytics services.

| Cloud | Resources |
|-------|-----------|
| AWS | `vpc`, `subnet` (private), `db_subnet_group`, `db_instance` (db.t3.micro, PostgreSQL 16, encrypted), `s3_bucket` ×2 (raw, processed) |
| Azure | `resource_group`, `storage_account` (StorageV2, HNS enabled = Data Lake Gen2), `postgresql_flexible_server` (Standard_B1ms, PostgreSQL 16) |
| GCP | `storage_bucket` ×3 (raw, processed — EU multi-region), `bigquery_dataset` (EU), `sql_database_instance` (db-f1-micro, POSTGRES_16) |
| OCI | `identity_compartment`, `objectstorage_bucket` ×2 (raw, processed), `database_autonomous_database` (DW workload, free tier) |

CIDR block: `10.4.0.0/16` (AWS only — other clouds have no VPC in this architecture)

---

### messaging

Asynchronous message queue / pub-sub service. No VPC or networking resources — messaging services are fully managed and accessed via endpoints.

| Cloud | Resources |
|-------|-----------|
| AWS | `sqs_queue` ×2 — main queue (1-day retention) + DLQ (14-day retention) |
| Azure | `resource_group`, `servicebus_namespace` (Basic SKU), `servicebus_queue` |
| GCP | `pubsub_topic`, `pubsub_subscription` (30s ack deadline, 8h message retention) |
| OCI | `identity_compartment`, `streaming_stream_pool`, `streaming_stream` (1 partition, 24h retention) |

---

## How to deploy a project

```bash
# Via Docker (recommended)
docker compose run --rm \
  -e TF_VAR_project=networking \
  engine-aws-terraform terraform plan

docker compose run --rm \
  -e TF_VAR_project=networking \
  engine-aws-terraform terraform apply

# Via host workspace (no Docker)
task workspace:aws PROJECT=networking
terraform -chdir=workspace plan
terraform -chdir=workspace apply
```

Replace `aws` / `networking` with any cloud and project name. See [`docs/docker-and-workspace.md`](docker-and-workspace.md) for full options and [`docs/tasks.md`](tasks.md) for the complete task reference.

---

## Adding a custom project

1. Create a directory under `data/<cloud>/<your-project>/`.
2. Add a `backend.tf.json` if you want remote state (copy from an existing project and update the key).
3. Add one JSON file per resource instance under `data/<cloud>/<your-project>/<resource_type>/`.  
   The resource type name is the Terraform type with the provider prefix stripped — e.g. `aws_s3_bucket` → `s3_bucket`.
4. Run a plan to verify:
   ```bash
   docker compose run --rm -e TF_VAR_project=your-project engine-aws-terraform terraform plan
   ```

No `.tf` changes are needed. The generated templates in `providers/<cloud>/` already cover every resource type the provider exposes.

---

← [Documentation index](README.md)
