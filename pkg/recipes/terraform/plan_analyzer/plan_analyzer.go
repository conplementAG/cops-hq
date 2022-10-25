package plan_analyzer

import (
	"errors"
	"github.com/conplementag/cops-hq/v2/internal"
	"github.com/conplementag/cops-hq/v2/pkg/recipes/terraform"
	"os"
	"path/filepath"
	"strings"
)

type PlanAnalyzer interface {
	// IsDeployPlanDirty returns true, if deploy plan shows required changes in the Terraform deployment
	IsDeployPlanDirty() (bool, error)

	// IsDestroyPlanDirty returns true, if destroy plan shows required changes in the Terraform deployment
	IsDestroyPlanDirty() (bool, error)
}

type planAnalyzer struct {
	projectName        string
	terraformDirectory string
}

func New(projectName string, terraformDirectory string) PlanAnalyzer {
	return &planAnalyzer{
		projectName:        projectName,
		terraformDirectory: terraformDirectory,
	}
}

func (pa *planAnalyzer) IsDeployPlanDirty() (bool, error) {
	result, err := pa.isPlanDirty(false)
	return result, internal.ReturnErrorOrPanic(err)
}

func (pa *planAnalyzer) IsDestroyPlanDirty() (bool, error) {
	result, err := pa.isPlanDirty(true)
	return result, internal.ReturnErrorOrPanic(err)
}

// isPlanDirty we will check the plan by using the simple text file which is always created via plan flows. This is much less work, and
// "probably" more reliable than analyzing the json file, since the JSON file has no single dirty property to check.  Parsing
// the JSON file would mean checking each resource for required action, which is also very format dependent (and could easily break
// with new terraform versions).
func (pa *planAnalyzer) isPlanDirty(checkForDestroy bool) (bool, error) {
	localTerraformRelativePlanFilePath, err := terraform.GetLocalTerraformRelativePlanFilePath(pa.projectName, pa.terraformDirectory, checkForDestroy)

	if err != nil {
		return true, err
	}

	planTextVersionFilePath := filepath.Join(pa.terraformDirectory, localTerraformRelativePlanFilePath+".txt")
	contents, err := os.ReadFile(planTextVersionFilePath)
	if err != nil {
		return true, err
	}

	// in the context of terraform plans, true is here a much better default value to prevent accidental changes due to bugs in this code
	var isDirty = true

	if strings.Contains(string(contents), "Your infrastructure matches the configuration.") &&
		strings.Contains(string(contents), "found no differences, so no changes are needed") {
		isDirty = false
	} else if strings.Contains(string(contents), "Terraform will perform the following actions") &&
		strings.Contains(string(contents), "To perform exactly these actions, run the following command to apply") {
		isDirty = true
	} else {
		return true, errors.New("we could not determine if the plan was dirty or not, because expected strings for " +
			"both checks were not found in the plan text file")
	}

	return isDirty, nil
}
