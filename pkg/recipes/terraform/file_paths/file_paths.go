package file_paths

import (
	"github.com/conplementag/cops-hq/v2/internal"
	"os"
	"path/filepath"
)

var PlansDirectory = ".plans"

// GetLocalTerraformRelativePlanFilePath gets the plan file path, relative to the given terraform directory. Output of this method
// is usually used with terraform commands, which already have the root directory set with -chdir
func GetLocalTerraformRelativePlanFilePath(projectName string, terraformDirectory string, getForDestroyPlan bool) (string, error) {
	// This is the place where plans directory is always ensured it exists.
	// Saving the plan to separate directory that terraform directory (which would be the default normal choice),
	// so that if the directory is mounted somewhere (like in Dockerfile / CD process), it will only have access to
	// the plan files, and not the whole local state cache (.e.g mounting the directory above would also expose the
	// contents of .terraform directory, and all the terraform files as well).
	fullPlansDirectoryPath := filepath.Join(terraformDirectory, PlansDirectory)
	err := os.MkdirAll(fullPlansDirectoryPath, os.ModePerm)
	if err != nil {
		return "", internal.ReturnErrorOrPanic(err)
	}

	// terraform file paths are always relative to terraform path set on -chdir, which we set to root where the sources are located
	return filepath.Join(PlansDirectory, GetPlanFileName(projectName, getForDestroyPlan)), nil
}

// GetPlanFileName return the Terraform plan file name
func GetPlanFileName(projectName string, getForDestroyPlan bool) string {
	var planFileName string

	if !getForDestroyPlan {
		planFileName = projectName + ".deploy.tfplan"
	} else {
		planFileName = projectName + ".destroy.tfplan"
	}

	return planFileName
}
