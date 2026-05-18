package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func awsPlanOpts(t *testing.T, project string) *terraform.Options {
	t.Helper()
	workspaceDir := prepareWorkspace(t, "aws", project)
	return terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: workspaceDir,
		EnvVars:      mergeEnvVars(t, awsEnvVars()),
		NoColor:      true,
	})
}

func TestAwsNetworkingPlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t, "ENGINE_AWS_ACCESS_KEY_ID", "ENGINE_AWS_SECRET_ACCESS_KEY", "ENGINE_AWS_DEFAULT_REGION")

	opts := awsPlanOpts(t, "networking")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.Contains(t, plan.ResourceChangesMap, `aws_vpc.grauss["vpc-grauss-networking"]`)
	assert.Contains(t, plan.ResourceChangesMap, `aws_internet_gateway.grauss["igw-grauss-networking"]`)
	assert.Contains(t, plan.ResourceChangesMap, `aws_security_group.grauss["sg-default-grauss-networking"]`)
}

func TestAwsVmSimplePlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t, "ENGINE_AWS_ACCESS_KEY_ID", "ENGINE_AWS_SECRET_ACCESS_KEY", "ENGINE_AWS_DEFAULT_REGION")

	opts := awsPlanOpts(t, "vm-simple")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.Contains(t, plan.ResourceChangesMap, `aws_instance.grauss["vm-grauss-simple"]`)
	assert.Contains(t, plan.ResourceChangesMap, `aws_eip.grauss["eip-grauss-vm-simple"]`)
}

func TestAwsWebAppPlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t, "ENGINE_AWS_ACCESS_KEY_ID", "ENGINE_AWS_SECRET_ACCESS_KEY", "ENGINE_AWS_DEFAULT_REGION")

	opts := awsPlanOpts(t, "web-app")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.Contains(t, plan.ResourceChangesMap, `aws_alb.grauss["alb-grauss-web-app"]`)
	assert.Contains(t, plan.ResourceChangesMap, `aws_instance.grauss["vm-app-grauss-web-app"]`)
	assert.Contains(t, plan.ResourceChangesMap, `aws_db_instance.grauss["db-grauss-web-app"]`)
}

func TestAwsKubernetesPlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t, "ENGINE_AWS_ACCESS_KEY_ID", "ENGINE_AWS_SECRET_ACCESS_KEY", "ENGINE_AWS_DEFAULT_REGION")

	opts := awsPlanOpts(t, "kubernetes")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.Contains(t, plan.ResourceChangesMap, `aws_eks_cluster.grauss["eks-grauss-kubernetes"]`)
	assert.Contains(t, plan.ResourceChangesMap, `aws_eks_node_group.grauss["ng-grauss-kubernetes"]`)
}

func TestAwsDataLakePlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t, "ENGINE_AWS_ACCESS_KEY_ID", "ENGINE_AWS_SECRET_ACCESS_KEY", "ENGINE_AWS_DEFAULT_REGION")

	opts := awsPlanOpts(t, "data-lake")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.Contains(t, plan.ResourceChangesMap, `aws_s3_bucket.grauss["bkt-grauss-data-lake-raw"]`)
	assert.Contains(t, plan.ResourceChangesMap, `aws_s3_bucket.grauss["bkt-grauss-data-lake-processed"]`)
	assert.Contains(t, plan.ResourceChangesMap, `aws_db_instance.grauss["db-grauss-data-lake"]`)
}

func TestAwsMessagingPlan(t *testing.T) {
	t.Parallel()
	skipIfMissingEnv(t, "ENGINE_AWS_ACCESS_KEY_ID", "ENGINE_AWS_SECRET_ACCESS_KEY", "ENGINE_AWS_DEFAULT_REGION")

	opts := awsPlanOpts(t, "messaging")

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, opts)

	assert.Contains(t, plan.ResourceChangesMap, `aws_sqs_queue.grauss["queue-grauss-main"]`)
	assert.Contains(t, plan.ResourceChangesMap, `aws_sqs_queue.grauss["queue-grauss-dlq"]`)
}
