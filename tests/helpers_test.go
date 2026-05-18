package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// repoRoot returns the absolute path of the repository root (parent of tests/).
func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs("..")
	require.NoError(t, err)
	return root
}

// pluginCacheDir returns the shared provider plugin cache path.
func pluginCacheDir(t *testing.T) string {
	t.Helper()
	dir := filepath.Join(repoRoot(t), "tmp", "tf-plugin-cache")
	require.NoError(t, os.MkdirAll(dir, 0755))
	return dir
}

// prepareWorkspace copies providers/<cloud>/ into a temp dir, injects version
// placeholders, and copies data/<cloud>/<project>/ into <tmp>/data/.
// The temp dir is automatically removed when the test ends (t.TempDir).
//
// The resulting directory is a self-contained Terraform root ready for init+plan.
func prepareWorkspace(t *testing.T, cloud, project string) string {
	t.Helper()
	root := repoRoot(t)

	providersDir := filepath.Join(root, "providers", cloud)
	dataDir := filepath.Join(root, "data", cloud, project)

	_, err := os.Stat(filepath.Join(providersDir, "_locals.tf"))
	require.NoError(t, err, "providers/%s not generated — run: task generate:%s", cloud, cloud)

	_, err = os.Stat(dataDir)
	require.NoError(t, err, "data/%s/%s not found", cloud, project)

	tmpDir := t.TempDir()

	// copy providers → tmpDir (exclude transient terraform dirs)
	rsync(t, providersDir+"/", tmpDir+"/",
		"--exclude=.terraform/",
		"--exclude=.terraform.lock.hcl",
		"--exclude=backend.remote.tf.json",
	)

	// mark engine binary so _locals.tf resolves data from the workspace dir
	require.NoError(t, os.WriteFile(
		filepath.Join(tmpDir, ".iac_engine_bin"), []byte("terraform"), 0644,
	))

	// inject __*_VERSION__ placeholders
	injectScript := filepath.Join(root, "scripts", "inject-provider-versions.sh")
	run(t, "bash", injectScript, tmpDir)

	// copy project data → tmpDir/data/
	// _locals.tf falls back to ${path.module}/data when data/<cloud>/ is absent
	dataDestDir := filepath.Join(tmpDir, "data")
	require.NoError(t, os.MkdirAll(dataDestDir, 0755))
	rsync(t, dataDir+"/", dataDestDir+"/")

	return tmpDir
}

// skipIfMissingEnv skips the test when any of the required environment variables are empty.
func skipIfMissingEnv(t *testing.T, vars ...string) {
	t.Helper()
	for _, v := range vars {
		if os.Getenv(v) == "" {
			t.Skipf("skipping: %s is not set", v)
		}
	}
}

// azureEnvVars maps ENGINE_ARM_* vars to the ARM_* names the AzureRM provider expects.
func azureEnvVars() map[string]string {
	return map[string]string{
		"ARM_CLIENT_ID":       os.Getenv("ENGINE_ARM_CLIENT_ID"),
		"ARM_CLIENT_SECRET":   os.Getenv("ENGINE_ARM_CLIENT_SECRET"),
		"ARM_TENANT_ID":       os.Getenv("ENGINE_ARM_TENANT_ID"),
		"ARM_SUBSCRIPTION_ID": os.Getenv("ENGINE_ARM_SUBSCRIPTION_ID"),
		"ARM_USE_AZUREAD":     "false",
	}
}

// awsEnvVars maps ENGINE_AWS_* vars to the AWS_* names the AWS provider expects.
func awsEnvVars() map[string]string {
	return map[string]string{
		"AWS_ACCESS_KEY_ID":     os.Getenv("ENGINE_AWS_ACCESS_KEY_ID"),
		"AWS_SECRET_ACCESS_KEY": os.Getenv("ENGINE_AWS_SECRET_ACCESS_KEY"),
		"AWS_SESSION_TOKEN":     os.Getenv("ENGINE_AWS_SESSION_TOKEN"),
		"AWS_DEFAULT_REGION":    os.Getenv("ENGINE_AWS_DEFAULT_REGION"),
	}
}

// gcpEnvVars maps ENGINE_GOOGLE_* vars to the GOOGLE_* names the Google provider expects.
func gcpEnvVars() map[string]string {
	return map[string]string{
		"GOOGLE_PROJECT":                  os.Getenv("ENGINE_GOOGLE_PROJECT"),
		"GOOGLE_REGION":                   os.Getenv("ENGINE_GOOGLE_REGION"),
		"GOOGLE_CREDENTIALS":              os.Getenv("ENGINE_GOOGLE_CREDENTIALS"),
		"GOOGLE_APPLICATION_CREDENTIALS":  os.Getenv("ENGINE_GOOGLE_APPLICATION_CREDENTIALS"),
	}
}

// ociEnvVars maps ENGINE_OCI_* vars to the OCI_* names the OCI provider expects.
func ociEnvVars() map[string]string {
	return map[string]string{
		"TF_VAR_tenancy_ocid":     os.Getenv("ENGINE_OCI_TENANCY_OCID"),
		"TF_VAR_user_ocid":        os.Getenv("ENGINE_OCI_USER_OCID"),
		"TF_VAR_fingerprint":      os.Getenv("ENGINE_OCI_FINGERPRINT"),
		"TF_VAR_private_key_path": os.Getenv("ENGINE_OCI_PRIVATE_KEY_FILE"),
		"TF_VAR_region":           os.Getenv("ENGINE_OCI_REGION"),
	}
}

// mergeEnvVars merges multiple env var maps into one, adding TF_PLUGIN_CACHE_DIR.
func mergeEnvVars(t *testing.T, maps ...map[string]string) map[string]string {
	t.Helper()
	out := map[string]string{
		"TF_PLUGIN_CACHE_DIR": pluginCacheDir(t),
	}
	for _, m := range maps {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}

// rsync copies src → dst using rsync, passing optional extra flags.
func rsync(t *testing.T, src, dst string, extraFlags ...string) {
	t.Helper()
	args := append([]string{"-a"}, extraFlags...)
	args = append(args, src, dst)
	run(t, "rsync", args...)
}

// run executes a command, failing the test on error.
func run(t *testing.T, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	require.NoError(t, cmd.Run(), "command failed: %s %v", name, args)
}
