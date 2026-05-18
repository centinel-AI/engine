# Grauss Reference Architectures

This document describes the six reference architecture projects included in the Grauss IaC engine. Each architecture is implemented across all five supported clouds: AWS, Azure, GCP, OCI, and OVH.

---

## Overview

| Architecture  | CIDR (2nd octet) | Purpose                                              |
|---------------|------------------|------------------------------------------------------|
| networking    | 10.0.x.x         | Foundational network layer: VPC, subnets, firewall   |
| vm-simple     | 10.1.x.x         | Minimal single-VM deployment with public IP          |
| web-app       | 10.2.x.x         | VM + load balancer + managed PostgreSQL database     |
| kubernetes    | 10.3.x.x         | Managed Kubernetes cluster + node pool               |
| data-lake     | 10.4.x.x         | Object storage + managed database/data warehouse     |
| messaging     | n/a              | Message queue / pub-sub service (no networking)      |

All architectures target **eu-west-1** (AWS), **westeurope** (Azure), **europe-west1** (GCP), **eu-frankfurt-1** (OCI), and **GRA9** (OVH — Gravelines, France) as default regions.

---

## Architecture 1: networking

**Purpose:** Establishes a baseline network layer that other architectures can reference. Demonstrates VPC creation, public and private subnets, internet gateway, route tables, and security groups across all four clouds.

### Resources per cloud

| Cloud | Resources                                                                                                      |
|-------|----------------------------------------------------------------------------------------------------------------|
| AWS   | vpc, subnet (public-a, private-a), internet_gateway, route_table, security_group                              |
| Azure | resource_group, virtual_network, subnet (public, private), network_security_group                             |
| GCP   | compute_network, compute_subnetwork, compute_firewall (allow-internal, allow-ssh)                              |
| OCI   | identity_compartment, core_vcn, core_internet_gateway, core_route_table, core_security_list, core_subnet (public, private) |
| OVH   | cloud_project_network_private (VLAN 100, GRA9), cloud_project_network_private_subnet (public 10.0.1.0/24, private 10.0.2.0/24) |

### Cost notes

All resources in this architecture are **free or negligible cost** when idle. VPCs, subnets, route tables, and security groups incur no charges on their own. The internet gateway on AWS is free; data transfer costs apply only when traffic flows through it. OVH private networks (vRack) and VLAN configuration in Public Cloud carry no fixed charge.

---

## Architecture 2: vm-simple

**Purpose:** The smallest possible VM deployment with a public IP address. Useful as a jump host, a scratchpad instance, or to verify that the networking layer works end-to-end.

### Resources per cloud

| Cloud | Resources                                                                                         |
|-------|---------------------------------------------------------------------------------------------------|
| AWS   | vpc, subnet, security_group (SSH only), instance (t3.micro), eip                                 |
| Azure | resource_group, virtual_network, subnet, public_ip, network_interface, linux_virtual_machine (Standard_B1s) |
| GCP   | compute_network, compute_subnetwork, compute_firewall (SSH), compute_address (EXTERNAL), compute_instance (e2-micro) |
| OCI   | identity_compartment, core_vcn, core_internet_gateway, core_route_table, core_security_list, core_subnet, core_instance (VM.Standard.E2.1.Micro) |
| OVH   | cloud_project_network_private (VLAN 101), cloud_project_network_private_subnet (10.1.0.0/16), cloud_project_instance (b2-7, Ubuntu 22.04) |

### Cost notes

These are the **smallest available instance types** on each cloud:
- AWS `t3.micro`: ~$0.0104/hr (Free Tier eligible)
- Azure `Standard_B1s`: ~$0.0104/hr
- GCP `e2-micro`: ~$0.007/hr (Free Tier: 1 f1-micro free in us-* regions)
- OCI `VM.Standard.E2.1.Micro`: Always Free eligible (up to 2 instances)
- OVH `b2-7` (2 vCPU, 7 GB RAM, 50 GB SSD): ~€0.055/hr — no free tier; check the [OVH Public Cloud pricing](https://www.ovhcloud.com/en/public-cloud/prices/) for the current rate

---

## Architecture 3: web-app

**Purpose:** A typical three-tier web application: load balancer → application VM → managed database. Demonstrates how to wire security groups between tiers and configure a managed PostgreSQL instance.

### Resources per cloud

| Cloud | Resources                                                                                                                       |
|-------|---------------------------------------------------------------------------------------------------------------------------------|
| AWS   | vpc, subnet (public-a, public-b, private-a), internet_gateway, security_group (alb, app, db), alb, alb_target_group, alb_listener, instance (t3.micro), db_subnet_group, db_instance (db.t3.micro, PostgreSQL 16) |
| Azure | resource_group, virtual_network, subnet (app, db with delegation), network_security_group, public_ip, network_interface, linux_virtual_machine (Standard_B1s), postgresql_flexible_server (Standard_B1ms) |
| GCP   | compute_network, compute_subnetwork, compute_firewall (http, ssh), compute_instance (e2-micro), sql_database_instance (db-f1-micro, POSTGRES_16) |
| OCI   | identity_compartment, core_vcn, core_internet_gateway, core_route_table, core_security_list (ports 22/80/443), core_subnet, core_instance (VM.Standard.E2.1.Micro), database_autonomous_database (OLTP, free tier) |
| OVH   | cloud_project_network_private (VLAN 102), cloud_project_network_private_subnet (10.2.0.0/16), cloud_project_instance (b2-7), cloud_project_database (PostgreSQL 16, essential/db1-4) |

### Cost notes

- AWS ALB: ~$0.008/hr + $0.008 per LCU. `db.t3.micro` RDS: ~$0.017/hr
- Azure `postgresql_flexible_server` Standard_B1ms: ~$0.041/hr
- GCP Cloud SQL `db-f1-micro`: ~$0.025/hr (not free tier)
- OCI Autonomous Database with `is_free_tier: true`: Always Free (up to 2 ADB instances)
- OVH `b2-7` instance: ~€0.055/hr. `cloud_project_database` essential/db1-4 (1 vCPU, 4 GB): ~€0.069/hr — no free tier

> The AWS ALB and RDS instance are the main cost drivers. Stop or destroy these when not in use. On OVH, the managed database is the main cost driver; destroy it between tests.

---

## Architecture 4: kubernetes

**Purpose:** A minimal managed Kubernetes cluster with a single node pool. Suitable for testing Kubernetes workloads without self-managing the control plane.

### Resources per cloud

| Cloud | Resources                                                                                                  |
|-------|------------------------------------------------------------------------------------------------------------|
| AWS   | vpc, subnet (private-a, private-b with EKS tags), internet_gateway, security_group, eks_cluster (v1.31), eks_node_group (t3.small, 1 node) |
| Azure | resource_group, kubernetes_cluster (AKS v1.31, Standard_B2s, 1 node, Free tier, SystemAssigned identity)  |
| GCP   | compute_network, compute_subnetwork, container_cluster (GKE, remove_default_node_pool), container_node_pool (e2-small, 1 node) |
| OCI   | identity_compartment, core_vcn, core_internet_gateway, core_route_table, core_security_list (port 6443), core_subnet, containerengine_cluster (OKE v1.31.1), containerengine_node_pool (VM.Standard.E2.1.Micro, 1 node) |
| OVH   | cloud_project_network_private (VLAN 103), cloud_project_network_private_subnet (10.3.0.0/16), cloud_project_kube (MKS v1.31, GRA9), cloud_project_kube_nodepool (b2-7, 1 node) |

### Cost notes

Kubernetes clusters are **more expensive** than simple VMs due to the control plane charges:
- AWS EKS control plane: **$0.10/hr** (~$72/month) + node costs. `t3.small` node: ~$0.023/hr
- Azure AKS with `sku_tier: Free`: **control plane is free**. `Standard_B2s` node: ~$0.042/hr
- GCP GKE Autopilot is free for one Autopilot cluster; Standard GKE charges $0.10/hr for cluster management + node costs. `e2-small`: ~$0.017/hr
- OCI OKE control plane: **free**. `VM.Standard.E2.1.Micro` node: Always Free eligible
- OVH MKS (Managed Kubernetes Service): **control plane is free**. `b2-7` worker node: ~€0.055/hr

> Destroy Kubernetes clusters when not in use. EKS and GKE cluster management fees accumulate even with zero workloads.

---

## Architecture 5: data-lake

**Purpose:** Object storage for raw and processed data layers, plus a managed database or data warehouse. Demonstrates multi-bucket patterns and cloud-native analytics services.

### Resources per cloud

| Cloud | Resources                                                                                                                 |
|-------|---------------------------------------------------------------------------------------------------------------------------|
| AWS   | vpc, subnet (private), db_subnet_group, db_instance (db.t3.micro, PostgreSQL 16, encrypted), s3_bucket (raw, processed)  |
| Azure | resource_group, storage_account (StorageV2, HNS enabled = Data Lake Gen2), postgresql_flexible_server (Standard_B1ms)    |
| GCP   | storage_bucket (raw, processed, EU multi-region), bigquery_dataset (EU), sql_database_instance (db-f1-micro, PostgreSQL 16) |
| OCI   | identity_compartment, objectstorage_bucket (raw, processed), database_autonomous_database (DW workload, free tier)        |
| OVH   | cloud_project_database (PostgreSQL 16, essential/db1-4); additional s3_bucket resources can be added via the same `s3_bucket` data directory using OVH Object Storage (requires `ENGINE_OVH_S3_*` credentials set from the `_bootstrap` output) |

### Cost notes

- AWS S3: $0.023/GB/month. `db.t3.micro` RDS: ~$0.017/hr
- Azure Storage (LRS): $0.018/GB/month. `Standard_B1ms` PostgreSQL: ~$0.041/hr
- GCP Cloud Storage (EU multi-region): $0.026/GB/month. BigQuery: first 10 GB/month storage free, $5/TB queries
- OCI Object Storage: 20 GB free (Always Free). ADB with `db_workload: DW` and `is_free_tier: true`: Always Free
- OVH `cloud_project_database` essential/db1-4: ~€0.069/hr. OVH Object Storage (S3-compatible): ~€0.0119/GB/month — no free tier

---

## Architecture 6: messaging

**Purpose:** A message queue or pub-sub service for asynchronous workloads. This architecture has no networking layer — messaging services are managed and accessed via endpoints.

### Resources per cloud

| Cloud | Resources                                                                               |
|-------|-----------------------------------------------------------------------------------------|
| AWS   | sqs_queue (main queue, 1-day retention), sqs_queue (DLQ, 14-day retention)             |
| Azure | resource_group, servicebus_namespace (Basic SKU), servicebus_queue                     |
| GCP   | pubsub_topic, pubsub_subscription                                                       |
| OCI   | identity_compartment, streaming_stream_pool, streaming_stream (1 partition, 24hr retention) |
| OVH   | cloud_project_database (Kafka 3.7, business/db2-7, 3 nodes — OVH Managed Kafka; minimum 3-node HA cluster required by the business plan) |

### Cost notes

All messaging resources chosen here fall in the **free or very low cost** tier, except OVH:
- AWS SQS: first 1M requests/month free, then $0.40/million
- Azure Service Bus Basic: $0.05/million operations. No namespace idle charge
- GCP Pub/Sub: first 10 GB/month free, then $0.04/GB
- OCI Streaming: 1 MB/s ingress and 1 partition included in Always Free
- OVH Managed Kafka (business/db2-7, 3 nodes): significant ongoing cost — the business plan requires 3 nodes for HA. **Destroy immediately after testing.** Check the [OVH Managed Databases pricing](https://www.ovhcloud.com/en/public-cloud/prices/) for the current rate.

---

## How to deploy

Grauss deploys via Docker images or a local workspace. See [`docs/docker-and-workspace.md`](../docs/docker-and-workspace.md) for full details.

**Via Docker (recommended):**

```bash
# 1. Build the image for the target cloud + engine
task build:aws:terraform

# 2. Plan
docker compose run --rm -e TF_VAR_project=networking engine-aws-terraform terraform plan

# 3. Apply
docker compose run --rm -e TF_VAR_project=networking engine-aws-terraform terraform apply

# 4. Destroy
docker compose run --rm -e TF_VAR_project=networking engine-aws-terraform terraform destroy
```

**Via host workspace (no Docker):**

```bash
# Populate workspace/ for the target cloud + project
task workspace:aws PROJECT=networking

# Then run Terraform directly
terraform -chdir=workspace init
terraform -chdir=workspace plan
terraform -chdir=workspace apply
```

Replace `aws` / `engine-aws-terraform` with `azure`, `gcp`, `oci`, or `ovh` as needed, and `networking` with any architecture name (`vm-simple`, `web-app`, `kubernetes`, `data-lake`, `messaging`).

For OVH, set `TF_VAR_ovh_s3_*` env vars (from `_bootstrap` output) before running projects that use `s3_bucket`. The `_bootstrap` project itself uses a two-phase apply — see [`docs/projects.md`](../docs/projects.md) for the OVH bootstrap sequence.

See [`docs/tasks.md`](../docs/tasks.md) for the full list of available task commands.

---

## Placeholder values that require replacement

Before deploying any architecture, search the relevant data directory for `REPLACE_WITH_` tokens and substitute real values:

| Placeholder                                        | What to substitute                                                   |
|----------------------------------------------------|----------------------------------------------------------------------|
| `REPLACE_WITH_VPC_ID`                              | The `id` output of the `aws_vpc` resource after first apply          |
| `REPLACE_WITH_SUBNET_ID` / `_A_ID` / `_B_ID`      | Subnet IDs from the `aws_subnet` resource outputs                    |
| `REPLACE_WITH_SG_ID` / `_ALB_SG_ID` / `_APP_SG_ID`| Security group IDs from `aws_security_group` outputs                 |
| `REPLACE_WITH_IGW_ID`                              | Internet gateway ID from `aws_internet_gateway`                      |
| `REPLACE_WITH_ALB_ARN` / `REPLACE_WITH_TG_ARN`    | ARNs from `aws_lb` and `aws_lb_target_group` outputs                 |
| `REPLACE_WITH_EKS_ROLE_ARN`                        | IAM role ARN with `AmazonEKSClusterPolicy` attached                  |
| `REPLACE_WITH_NODE_ROLE_ARN`                       | IAM role ARN with `AmazonEKSWorkerNodePolicy` etc. attached          |
| `ami-REPLACE_WITH_UBUNTU_2204_AMI`                 | Region-specific Ubuntu 22.04 AMI ID (use SSM or EC2 console)        |
| `REPLACE_WITH_INSTANCE_ID`                         | EC2 instance ID for EIP association                                   |
| `REPLACE_WITH_NIC_ID` / `REPLACE_WITH_PIP_ID`     | Azure NIC / Public IP resource IDs from portal or outputs            |
| `REPLACE_WITH_SSH_PUBLIC_KEY`                      | Contents of your `~/.ssh/id_rsa.pub` or equivalent                   |
| `REPLACE_WITH_SERVICEBUS_NAMESPACE_ID`             | Azure Service Bus namespace resource ID                               |
| `REPLACE_WITH_COMPARTMENT_OCID`                    | OCI compartment OCID from Console → Identity → Compartments          |
| `ocid1.tenancy.oc1..aaaaaaaaREPLACE_...`           | Your OCI tenancy OCID (root compartment)                              |
| `REPLACE_WITH_VCN_OCID`                            | OCI VCN OCID from Console → Networking → Virtual Cloud Networks      |
| `REPLACE_WITH_IGW_OCID`                            | OCI Internet Gateway OCID                                             |
| `REPLACE_WITH_SUBNET_OCID`                         | OCI Subnet OCID                                                       |
| `REPLACE_WITH_IMAGE_OCID`                          | OCI compute image OCID (use `oci compute image list`)                 |
| `REPLACE_WITH_AVAILABILITY_DOMAIN`                 | OCI availability domain name (e.g. `efXJ:EU-FRANKFURT-1-AD-1`)       |
| `REPLACE_WITH_OCI_NAMESPACE`                       | OCI Object Storage namespace (use `oci os ns get`)                    |
| `REPLACE_WITH_OKE_CLUSTER_OCID`                    | OCI OKE cluster OCID after creating the cluster                       |
| `REPLACE_WITH_STREAM_POOL_OCID`                    | OCI Streaming stream pool OCID                                        |
| `REPLACE_WITH_STRONG_PASSWORD`                     | A password meeting cloud provider complexity requirements             |
| `REPLACE_WITH_SERVICE_NAME`                        | OVH Cloud project ID (UUID visible in the OVH Control Panel under Public Cloud) |
| `REPLACE_WITH_NETWORK_ID`                          | OVH private network ID — output of `cloud_project_network_private` after first apply |
| `REPLACE_WITH_KUBE_ID`                             | OVH MKS cluster ID — output of `cloud_project_kube` after apply      |
| `REPLACE_WITH_USER_ID`                             | OVH cloud project user ID — output of `cloud_project_user` (used by `cloud_project_user_s3_credential`) |
| `REPLACE_WITH_SSH_KEY_NAME`                        | Name of the SSH key registered in the OVH Public Cloud project (Control Panel → SSH Keys) |

### Password requirements
- AWS RDS: 8–41 chars, no `/`, `"`, `@`, or space
- Azure PostgreSQL Flexible: 8–128 chars, must include uppercase, lowercase, digit, and symbol
- OCI ADB: 12–30 chars, at least 2 uppercase, 2 lowercase, 2 digits, 2 special chars (not `"` or `@`)
- OVH Managed Database: 8–64 chars, must include uppercase, lowercase, digit, and symbol

> **Never commit real passwords or OCIDs to version control.** Use environment variables, a secrets manager, or Terraform variable files (`.tfvars`) that are excluded via `.gitignore`.
