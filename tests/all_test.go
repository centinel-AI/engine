//go:build exhaustive

package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAzureExhaustivePlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t,
		"ENGINE_ARM_CLIENT_ID", "ENGINE_ARM_CLIENT_SECRET",
		"ENGINE_ARM_TENANT_ID", "ENGINE_ARM_SUBSCRIPTION_ID",
	)
	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, azurePlanOpts(t, "_test-all"))
	assertExhaustivePlan(t, plan, "azure")
}

func TestAwsExhaustivePlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t, "ENGINE_AWS_ACCESS_KEY_ID", "ENGINE_AWS_SECRET_ACCESS_KEY", "ENGINE_AWS_DEFAULT_REGION")
	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, awsPlanOpts(t, "_test-all"))
	assertExhaustivePlan(t, plan, "aws")
}

func TestGcpExhaustivePlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t, "ENGINE_GOOGLE_PROJECT", "ENGINE_GOOGLE_CREDENTIALS")
	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, gcpPlanOpts(t, "_test-all"))
	assertExhaustivePlan(t, plan, "gcp")
}

func TestOciExhaustivePlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t,
		"ENGINE_OCI_TENANCY_OCID", "ENGINE_OCI_USER_OCID",
		"ENGINE_OCI_FINGERPRINT", "ENGINE_OCI_PRIVATE_KEY_FILE",
	)
	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, ociPlanOpts(t, "_test-all"))
	assertExhaustivePlan(t, plan, "oci")
}

// assertExhaustivePlan reads every fixture directory under data/<cloud>/_test-all/,
// extracts the "_resource_type" field from fixture.json, and asserts that the
// corresponding resource appears in the Terraform plan.
func assertExhaustivePlan(t *testing.T, plan *terraform.PlanStruct, cloud string) {
	t.Helper()

	fixtureDir := filepath.Join(repoRoot(t), "data", cloud, "_test-all")
	entries, err := os.ReadDir(fixtureDir)
	require.NoError(t, err,
		"test fixtures not found — run: task generate:fixtures:%s", cloud)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(fixtureDir, entry.Name(), "fixture.json"))
		if err != nil {
			continue
		}
		var fixture map[string]any
		if err := json.Unmarshal(raw, &fixture); err != nil {
			continue
		}
		resourceType, _ := fixture["_resource_type"].(string)
		if resourceType == "" {
			continue
		}
		assert.Contains(t, plan.ResourceChangesMap,
			resourceType+`.grauss["test-fixture"]`,
			"resource type %s not found in plan", resourceType,
		)
	}
}
