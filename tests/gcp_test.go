package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func gcpPlanOpts(t *testing.T, project string) *terraform.Options {
	t.Helper()
	workspaceDir := prepareWorkspace(t, "gcp", project)
	return terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: workspaceDir,
		EnvVars:      mergeEnvVars(t, gcpEnvVars()),
		NoColor:      true,
	})
}

func TestGcpNetworkingPlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t, "ENGINE_GOOGLE_PROJECT", "ENGINE_GOOGLE_CREDENTIALS")

	opts := gcpPlanOpts(t, "networking")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.Contains(t, plan.ResourceChangesMap, `google_compute_network.grauss["vpc-grauss-networking"]`)
	assert.Contains(t, plan.ResourceChangesMap, `google_compute_subnetwork.grauss["subnet-grauss-networking"]`)
	assert.Contains(t, plan.ResourceChangesMap, `google_compute_firewall.grauss["fw-allow-ssh-grauss-networking"]`)
}

func TestGcpVmSimplePlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t, "ENGINE_GOOGLE_PROJECT", "ENGINE_GOOGLE_CREDENTIALS")

	opts := gcpPlanOpts(t, "vm-simple")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.Contains(t, plan.ResourceChangesMap, `google_compute_instance.grauss["vm-grauss-simple"]`)
	assert.Contains(t, plan.ResourceChangesMap, `google_compute_address.grauss["addr-grauss-vm-simple"]`)
}

func TestGcpWebAppPlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t, "ENGINE_GOOGLE_PROJECT", "ENGINE_GOOGLE_CREDENTIALS")

	opts := gcpPlanOpts(t, "web-app")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.Contains(t, plan.ResourceChangesMap, `google_compute_instance.grauss["vm-grauss-web-app"]`)
	assert.Contains(t, plan.ResourceChangesMap, `google_sql_database_instance.grauss["db-grauss-web-app"]`)
}

func TestGcpKubernetesPlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t, "ENGINE_GOOGLE_PROJECT", "ENGINE_GOOGLE_CREDENTIALS")

	opts := gcpPlanOpts(t, "kubernetes")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.Contains(t, plan.ResourceChangesMap, `google_container_cluster.grauss["gke-grauss-kubernetes"]`)
	assert.Contains(t, plan.ResourceChangesMap, `google_container_node_pool.grauss["np-grauss-kubernetes"]`)
}

func TestGcpDataLakePlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t, "ENGINE_GOOGLE_PROJECT", "ENGINE_GOOGLE_CREDENTIALS")

	opts := gcpPlanOpts(t, "data-lake")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.Contains(t, plan.ResourceChangesMap, `google_storage_bucket.grauss["bkt-grauss-data-lake-raw"]`)
	assert.Contains(t, plan.ResourceChangesMap, `google_bigquery_dataset.grauss["ds-grauss-data-lake"]`)
	assert.Contains(t, plan.ResourceChangesMap, `google_sql_database_instance.grauss["db-grauss-data-lake"]`)
}

func TestGcpMessagingPlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t, "ENGINE_GOOGLE_PROJECT", "ENGINE_GOOGLE_CREDENTIALS")

	opts := gcpPlanOpts(t, "messaging")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.Contains(t, plan.ResourceChangesMap, `google_pubsub_topic.grauss["topic-grauss-messaging"]`)
	assert.Contains(t, plan.ResourceChangesMap, `google_pubsub_subscription.grauss["sub-grauss-messaging"]`)
}
