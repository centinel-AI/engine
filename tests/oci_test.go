package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func ociPlanOpts(t *testing.T, project string) *terraform.Options {
	t.Helper()
	workspaceDir := prepareWorkspace(t, "oci", project)
	return terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: workspaceDir,
		EnvVars:      mergeEnvVars(t, ociEnvVars()),
		NoColor:      true,
	})
}

func TestOciNetworkingPlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t, "ENGINE_OCI_TENANCY_OCID", "ENGINE_OCI_USER_OCID", "ENGINE_OCI_FINGERPRINT", "ENGINE_OCI_PRIVATE_KEY_FILE")

	opts := ociPlanOpts(t, "networking")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.Contains(t, plan.ResourceChangesMap, `oci_core_vcn.grauss["vcn-grauss-networking"]`)
	assert.Contains(t, plan.ResourceChangesMap, `oci_core_internet_gateway.grauss["igw-grauss-networking"]`)
	assert.Contains(t, plan.ResourceChangesMap, `oci_core_security_list.grauss["sl-grauss-networking"]`)
}

func TestOciVmSimplePlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t, "ENGINE_OCI_TENANCY_OCID", "ENGINE_OCI_USER_OCID", "ENGINE_OCI_FINGERPRINT", "ENGINE_OCI_PRIVATE_KEY_FILE")

	opts := ociPlanOpts(t, "vm-simple")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.Contains(t, plan.ResourceChangesMap, `oci_core_instance.grauss["vm-grauss-simple"]`)
	assert.Contains(t, plan.ResourceChangesMap, `oci_core_vcn.grauss["vcn-grauss-vm-simple"]`)
}

func TestOciWebAppPlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t, "ENGINE_OCI_TENANCY_OCID", "ENGINE_OCI_USER_OCID", "ENGINE_OCI_FINGERPRINT", "ENGINE_OCI_PRIVATE_KEY_FILE")

	opts := ociPlanOpts(t, "web-app")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.Contains(t, plan.ResourceChangesMap, `oci_core_instance.grauss["vm-grauss-web-app"]`)
	assert.Contains(t, plan.ResourceChangesMap, `oci_database_autonomous_database.grauss["db-grauss-web-app"]`)
}

func TestOciKubernetesPlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t, "ENGINE_OCI_TENANCY_OCID", "ENGINE_OCI_USER_OCID", "ENGINE_OCI_FINGERPRINT", "ENGINE_OCI_PRIVATE_KEY_FILE")

	opts := ociPlanOpts(t, "kubernetes")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.Contains(t, plan.ResourceChangesMap, `oci_containerengine_cluster.grauss["oke-grauss-kubernetes"]`)
	assert.Contains(t, plan.ResourceChangesMap, `oci_containerengine_node_pool.grauss["np-grauss-kubernetes"]`)
}

func TestOciDataLakePlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t, "ENGINE_OCI_TENANCY_OCID", "ENGINE_OCI_USER_OCID", "ENGINE_OCI_FINGERPRINT", "ENGINE_OCI_PRIVATE_KEY_FILE")

	opts := ociPlanOpts(t, "data-lake")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.Contains(t, plan.ResourceChangesMap, `oci_objectstorage_bucket.grauss["bkt-grauss-data-lake-raw"]`)
	assert.Contains(t, plan.ResourceChangesMap, `oci_database_autonomous_database.grauss["db-grauss-data-lake"]`)
}

func TestOciMessagingPlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t, "ENGINE_OCI_TENANCY_OCID", "ENGINE_OCI_USER_OCID", "ENGINE_OCI_FINGERPRINT", "ENGINE_OCI_PRIVATE_KEY_FILE")

	opts := ociPlanOpts(t, "messaging")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.Contains(t, plan.ResourceChangesMap, `oci_streaming_stream_pool.grauss["pool-grauss-messaging"]`)
	assert.Contains(t, plan.ResourceChangesMap, `oci_streaming_stream.grauss["stream-grauss-messaging"]`)
}
