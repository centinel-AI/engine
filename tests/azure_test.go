package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// azurePlanOpts builds terraform.Options for an Azure plan-only test.
func azurePlanOpts(t *testing.T, project string) *terraform.Options {
	t.Helper()
	workspaceDir := prepareWorkspace(t, "azure", project)
	return terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: workspaceDir,
		EnvVars:      mergeEnvVars(t, azureEnvVars()),
		NoColor:      true,
	})
}

// TestAzureNetworkingPlan verifies that the networking project produces a valid
// plan containing the expected VNet and NSG resources.
func TestAzureNetworkingPlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t,
		"ENGINE_ARM_CLIENT_ID", "ENGINE_ARM_CLIENT_SECRET",
		"ENGINE_ARM_TENANT_ID", "ENGINE_ARM_SUBSCRIPTION_ID",
	)

	opts := azurePlanOpts(t, "networking")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.NotEmpty(t, plan.ResourceChangesMap, "plan should contain resource changes")
	assert.Contains(t, plan.ResourceChangesMap, `azurerm_resource_group.grauss["rg-grauss-networking"]`)
	assert.Contains(t, plan.ResourceChangesMap, `azurerm_virtual_network.grauss["vnet-grauss-networking"]`)
	assert.Contains(t, plan.ResourceChangesMap, `azurerm_network_security_group.grauss["nsg-grauss-networking"]`)
}

// TestAzureVmSimplePlan verifies that the vm-simple project plans a VM with a
// public IP and network interface.
func TestAzureVmSimplePlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t,
		"ENGINE_ARM_CLIENT_ID", "ENGINE_ARM_CLIENT_SECRET",
		"ENGINE_ARM_TENANT_ID", "ENGINE_ARM_SUBSCRIPTION_ID",
	)

	opts := azurePlanOpts(t, "vm-simple")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.Contains(t, plan.ResourceChangesMap, `azurerm_linux_virtual_machine.grauss["vm-grauss-simple"]`)
	assert.Contains(t, plan.ResourceChangesMap, `azurerm_public_ip.grauss["pip-grauss-vm-simple"]`)
	assert.Contains(t, plan.ResourceChangesMap, `azurerm_network_interface.grauss["nic-grauss-vm-simple"]`)
}

// TestAzureWebAppPlan verifies that the web-app project plans a VM and a
// PostgreSQL Flexible Server.
func TestAzureWebAppPlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t,
		"ENGINE_ARM_CLIENT_ID", "ENGINE_ARM_CLIENT_SECRET",
		"ENGINE_ARM_TENANT_ID", "ENGINE_ARM_SUBSCRIPTION_ID",
	)

	opts := azurePlanOpts(t, "web-app")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.Contains(t, plan.ResourceChangesMap, `azurerm_linux_virtual_machine.grauss["vm-grauss-web-app"]`)
	assert.Contains(t, plan.ResourceChangesMap, `azurerm_postgresql_flexible_server.grauss["db-grauss-web-app"]`)
}

// TestAzureKubernetesPlan verifies that the kubernetes project plans an AKS cluster.
func TestAzureKubernetesPlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t,
		"ENGINE_ARM_CLIENT_ID", "ENGINE_ARM_CLIENT_SECRET",
		"ENGINE_ARM_TENANT_ID", "ENGINE_ARM_SUBSCRIPTION_ID",
	)

	opts := azurePlanOpts(t, "kubernetes")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.Contains(t, plan.ResourceChangesMap, `azurerm_kubernetes_cluster.grauss["aks-grauss-kubernetes"]`)
}

// TestAzureDataLakePlan verifies that the data-lake project plans a storage
// account and a PostgreSQL Flexible Server.
func TestAzureDataLakePlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t,
		"ENGINE_ARM_CLIENT_ID", "ENGINE_ARM_CLIENT_SECRET",
		"ENGINE_ARM_TENANT_ID", "ENGINE_ARM_SUBSCRIPTION_ID",
	)

	opts := azurePlanOpts(t, "data-lake")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.Contains(t, plan.ResourceChangesMap, `azurerm_storage_account.grauss["stgraussdatalake01"]`)
	assert.Contains(t, plan.ResourceChangesMap, `azurerm_postgresql_flexible_server.grauss["db-grauss-data-lake"]`)
}

// TestAzureMessagingPlan verifies that the messaging project plans a Service Bus
// namespace and queue.
func TestAzureMessagingPlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t,
		"ENGINE_ARM_CLIENT_ID", "ENGINE_ARM_CLIENT_SECRET",
		"ENGINE_ARM_TENANT_ID", "ENGINE_ARM_SUBSCRIPTION_ID",
	)

	opts := azurePlanOpts(t, "messaging")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.Contains(t, plan.ResourceChangesMap, `azurerm_servicebus_namespace.grauss["sb-grauss-messaging"]`)
	assert.Contains(t, plan.ResourceChangesMap, `azurerm_servicebus_queue.grauss["queue-grauss-main"]`)
}
