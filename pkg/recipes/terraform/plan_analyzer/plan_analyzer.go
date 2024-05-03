package plan_analyzer

import (
	"github.com/conplementag/cops-hq/v2/internal"
	"github.com/conplementag/cops-hq/v2/pkg/recipes/terraform/file_paths"
	"os"
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

// isPlanDirty we will check the plan state file using the simple text file which is always created via plan flows.
func (pa *planAnalyzer) isPlanDirty(checkForDestroy bool) (bool, error) {
	if checkForDestroy {
		// mark destroy plans always as dirty - we do not create state during planning
		return true, nil
	}

	planStateFullFilePath, err := file_paths.GetLocalTerraformRelativePlanStateFilePath(pa.projectName, pa.terraformDirectory)
	if err != nil {
		return true, err
	}

	_, err = os.Stat(planStateFullFilePath)
	return os.IsNotExist(err), nil
}
