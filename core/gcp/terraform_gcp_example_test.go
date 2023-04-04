//go:build gcp
// +build gcp

// NOTE: We use build tags to differentiate GCP testing for better isolation and parallelism when executing our tests.

package test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

func TestTerraformGcpExample(t *testing.T) {
	t.Parallel()

	exampleDir := test_structure.CopyTerraformFolderToTemp(t, "../../", "examples/gcp/terraform-gcp-example")

	// Get the Project Id to use
	projectId := gcp.GetGoogleProjectIDFromEnvVar(t)

	// Create all resources in the following zone
	zone := "us-east1-b"

	// Give the example bucket a unique name so we can distinguish it from any other bucket in your GCP account
	expectedBucketName := fmt.Sprintf("terratest-gcp-example-%s", strings.ToLower(random.UniqueId()))

	// Also give the example instance a unique name
	expectedInstanceName := fmt.Sprintf("terratest-gcp-example-%s", strings.ToLower(random.UniqueId()))

	// website::tag::1::Configure Terraform setting path to Terraform code, bucket name, and instance name. Construct
	// the terraform options with default retryable errors to handle the most common retryable errors in terraform
	// testing.
	planFilePath := filepath.Join(exampleDir, "plan.out")
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: exampleDir,

		// Configure a plan file path so we can introspect the plan and make assertions about it.
		PlanFilePath: planFilePath,

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"gcp_project_id": projectId,
			"zone":           zone,
			"instance_name":  expectedInstanceName,
			"bucket_name":    expectedBucketName,
		},
	})

	// website::tag::5::At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// website::tag::2::This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	supportedResources := map[string]string{
		"google_storage_bucket":   "location",
		"google_compute_instance": "zone",
	}

	cloudIntensities := getCloudIntensities()

	for resourceName, resourceChanges := range plan.ResourceChangesMap {
		// storage:::
		// Get the best nearby region for it
		// If it is not the best region fail the test with a useful message
		// Need to handle multi region storage and allow opt-out

		// computer::
		// Need to figure out load balancing
		// Maybe also need to check that the desired compute instance type can be provisioned in the recommended location
		if resourceChanges.Change.Actions.Create() {
			regionKey, isSupported := supportedResources[resourceChanges.Type]
			if isSupported {
				plannedValues, _ := plan.ResourcePlannedValuesMap[resourceName]
				region := plannedValues.AttributeValues[regionKey].(string)
				for _, intensity := range cloudIntensities {
					t.Log(intensity.GeneralRegion, getGcpGeneralRegion(region))
				}

			}
		}
	}
}
