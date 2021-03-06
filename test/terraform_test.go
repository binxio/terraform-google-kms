package test

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/require"
)

// Give it a random string so we're sure it's created this test run
var expectedEnvironment string
var testPreq *testing.T
var terraformOptions *terraform.Options
var tmpSaReaderEmail string
var tmpSaOwnerEmail string
var blacklistRegions []string

func TestMain(m *testing.M) {
	expectedEnvironment = fmt.Sprintf("terratest %s", strings.ToLower(random.UniqueId()))
	blacklistRegions = []string{"asia-east2"}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		<-c
		TestCleanup(testPreq)
		os.Exit(1)
	}()

	result := m.Run()

	Clean()

	os.Exit(result)
}

// -------------------------------------------------------------------------------------------------------- //
// Utility functions
// -------------------------------------------------------------------------------------------------------- //
func setTerraformOptions(dir string, region string, projectId string) {
	terraformOptions = &terraform.Options {
		TerraformDir: dir,
		// Pass the expectedEnvironment for tagging
		Vars: map[string]interface{}{
			"environment": expectedEnvironment,
			"location": region,
			"sa_reader_email": tmpSaReaderEmail,
			"sa_owner_email": tmpSaOwnerEmail,
		},
		EnvVars: map[string]string{
			"GOOGLE_CLOUD_PROJECT": projectId,
		},
	}
}

// A build step that removes temporary build and test files
func Clean() error {
	fmt.Println("Cleaning...")

	return filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		if info.IsDir() && info.Name() == ".terraform" {
			os.RemoveAll(path)
			fmt.Printf("Removed \"%v\"\n", path)
			return filepath.SkipDir
		}
		if !info.IsDir() && (info.Name() == "terraform.tfstate" ||
		info.Name() == "terraform.tfplan" ||
		info.Name() == "terraform.tfstate.backup") {
			os.Remove(path)
			fmt.Printf("Removed \"%v\"\n", path)
		}
		return nil
	})
}

func Test_Prereq(t *testing.T) {
	projectId := gcp.GetGoogleProjectIDFromEnvVar(t)
	region := gcp.GetRandomRegion(t, projectId, nil, blacklistRegions)
	setTerraformOptions(".", region, projectId)
	testPreq = t

	terraform.InitAndApply(t, terraformOptions)

	tmpSaReaderEmail = terraform.OutputRequired(t, terraformOptions, "sa_reader_email")
	tmpSaOwnerEmail = terraform.OutputRequired(t, terraformOptions, "sa_owner_email")
}

// -------------------------------------------------------------------------------------------------------- //
// Unit Tests
// -------------------------------------------------------------------------------------------------------- //
/*
func TestUT_Assertions(t *testing.T) {
	// Pick a random GCP region to test in. This helps ensure your code works in all regions.
	projectId := gcp.GetGoogleProjectIDFromEnvVar(t)
	region := gcp.GetRandomRegion(t, projectId, nil, blacklistRegions)

	expectedAssertUnknownVar := "Unknown kms variable assigned"
	expectedAssertNameTooLong := "'s generated name is too long:"
	expectedAssertNameInvalidChars := "does not match regex"
	//expectedAssertKMSKeyMissing := "KMS Encryption key id is required."

	setTerraformOptions("assertions", region, projectId)

	out, err := terraform.InitAndPlanE(t, terraformOptions)

	require.Error(t, err)
	assert.Contains(t, out, expectedAssertUnknownVar)
	assert.Contains(t, out, expectedAssertNameTooLong)
	assert.Contains(t, out, expectedAssertNameInvalidChars)
	//assert.Contains(t, out, expectedAssertKMSKeyMissing)
}
*/

func TestUT_Defaults(t *testing.T) {
	projectId := gcp.GetGoogleProjectIDFromEnvVar(t)
	region := gcp.GetRandomRegion(t, projectId, nil, blacklistRegions)
	setTerraformOptions("defaults", region, projectId)
	terraform.InitAndPlan(t, terraformOptions)
}

/*
func TestUT_Overrides(t *testing.T) {
	projectId := gcp.GetGoogleProjectIDFromEnvVar(t)
	region := gcp.GetRandomRegion(t, projectId, nil, blacklistRegions)
	setTerraformOptions("overrides", region, projectId)
	terraform.InitAndPlan(t, terraformOptions)
}
*/

// -------------------------------------------------------------------------------------------------------- //
// Integration Tests
// -------------------------------------------------------------------------------------------------------- //

func TestIT_Defaults(t *testing.T) {
	projectId := gcp.GetGoogleProjectIDFromEnvVar(t)
	region := gcp.GetRandomRegion(t, projectId, nil, blacklistRegions)
	setTerraformOptions("defaults", region, projectId)

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	outputs := terraform.OutputAll(t, terraformOptions)

	// Ugly typecasting because Go....
	keyRing := outputs["key_ring_name"].(string)
	cryptoKey := outputs["crypto_key_name"].(string)

	// Make sure our names match with expected output
	fmt.Printf("Checking key ring and crypto key names (%s and %s) matches our expectations...\n", keyRing, cryptoKey)
	assert.Contains(t, keyRing, "testapp")
	assert.Contains(t, cryptoKey, "testapp")
}

/*
func TestIT_Overrides(t *testing.T) {
	projectId := gcp.GetGoogleProjectIDFromEnvVar(t)
	region := gcp.GetRandomRegion(t, projectId, nil, blacklistRegions)
	setTerraformOptions("overrides", region, projectId)

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	outputs := terraform.OutputAll(t, terraformOptions)

	// Ugly typecasting because Go....
	bucketMap := outputs["map"].(map[string]interface{})
	functionsBucket := bucketMap["functions"].(map[string]interface{})
	bucketId := functionsBucket["id"].(string)

	// Make sure our bucket is created
	fmt.Printf("Checking bucket %s...\n", bucketId)
	gcp.AssertStorageBucketExists(t, bucketId)
}
*/

func TestCleanup(t *testing.T) {
	fmt.Println("Cleaning possible lingering resources..")
	terraform.Destroy(t, terraformOptions)

	// Also clean up prereq. resources
	fmt.Println("Cleaning our prereq resources...")
	projectId := gcp.GetGoogleProjectIDFromEnvVar(t)
	region := gcp.GetRandomRegion(t, projectId, nil, blacklistRegions)
	setTerraformOptions(".", region, projectId)
	terraform.Destroy(t, terraformOptions)
}
