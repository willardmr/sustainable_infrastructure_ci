//go:build gcp
// +build gcp

// NOTE: We use build tags to differentiate GCP testing for better isolation and parallelism when executing our tests.

package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	terraform "github.com/gruntwork-io/terratest/modules/terraform"

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

	evaluatePlan(t, plan, 10)

}

func evaluatePlan(t *testing.T, plan *terraform.PlanStruct, maximumCostIncreasePercentage float64) {
	supportedResources := map[string]string{
		"google_storage_bucket": "location",
	}

	for resourceName, resourceChanges := range plan.ResourceChangesMap {
		if resourceChanges.Change.Actions.Create() {
			regionKey, isSupported := supportedResources[resourceChanges.Type]
			if isSupported {
				plannedValues, _ := plan.ResourcePlannedValuesMap[resourceName]
				region := plannedValues.AttributeValues[regionKey].(string)
				bestRegion := getBestRegion(region, maximumCostIncreasePercentage)
				if strings.ToLower(bestRegion.Region) != strings.ToLower(region) {
					assert := assert.New(t)
					assert.Equal(bestRegion.Region, region,
						fmt.Sprintf("Resource of type %s with name %s is being created in region %s, but it should be created in region %s to reduce carbon emissions.", resourceChanges.Type, plannedValues.AttributeValues["name"].(string), strings.ToLower(region), bestRegion.Region))
				}
			}
		}
	}
}

func getCostIncreasePercentage(currentRegion string, newRegion string) float64 {
	cloudCosts := getCloudCosts()
	return 100 * (cloudCosts[strings.ToLower(newRegion)].Cost - cloudCosts[strings.ToLower(currentRegion)].Cost) / cloudCosts[strings.ToLower(newRegion)].Cost
}

/* Best means lowest carbon intensity without too high of a price increase */
func getBestRegion(currentRegion string, maximumCostIncreasePercentage float64) CloudIntensity {
	cloudIntensities := getCloudIntensities()

	lowestIntensity := 100000.0
	var bestRegion CloudIntensity
	generalRegion := getGcpGeneralRegion(currentRegion)
	for _, intensity := range cloudIntensities {
		if strings.ToLower(intensity.GeneralRegion) == strings.ToLower(generalRegion) && intensity.Impact < lowestIntensity && getCostIncreasePercentage(currentRegion, intensity.Region) <= maximumCostIncreasePercentage {
			bestRegion = intensity
			lowestIntensity = intensity.Impact
		}
	}
	return bestRegion
}
