package terraform

import (
	"github.com/conplementag/cops-hq/v2/pkg/recipes/terraform/plan_analyzer"
	"os"
	"path/filepath"
)

// persistPlanInAdditionalFormatsOnDisk - we also persist the plan output to disk in both human-readable and json formats,
// which can be later be processed, without requiring terraform init & terraform show separately to achieve the same result.
func (tf *terraformWrapper) persistPlanInAdditionalFormatsOnDisk(planAsPlaintext string, terraformRelativePlanFilePath string) error {
	// to persist the plan in other file formats, we need to convert the terraformRelativePlanFilePath to a path
	// resolvable from where we are running at the moment (e.g. cmd/example-cli).
	planFullFilePath := filepath.Join(tf.terraformDirectory, terraformRelativePlanFilePath)

	// 1. human-readable form we already have from planOutput, so we just persist it
	textFile, err := os.Create(planFullFilePath + ".txt")
	defer textFile.Close()
	if err != nil {
		return err
	}

	_, err = textFile.WriteString(planAsPlaintext)
	if err != nil {
		return err
	}

	// 2. json form we need to get with an extra terraform call. Since init is already done, this will work
	// also, we use the terraformRelativePlanFilePath since this is a terraform command, initialized with -chdir
	jsonPlanOutput, err := tf.executor.Execute("terraform -chdir=" + tf.terraformDirectory + " show -json " + terraformRelativePlanFilePath)
	if err != nil {
		return err
	}

	jsonFile, err := os.Create(planFullFilePath + ".json")
	defer jsonFile.Close()
	if err != nil {
		return err
	}

	_, err = jsonFile.WriteString(jsonPlanOutput)
	if err != nil {
		return err
	}

	return nil
}

// persistAnalysisResultOnDisk - we also run the plan analyzer and persist the result as a file, in case plan contains no changes.
func (tf *terraformWrapper) persistAnalysisResultOnDisk(terraformRelativePlanFilePath string) error {
	analyzer := plan_analyzer.New(tf.projectName, tf.terraformDirectory)

	isPlanDirty, err := analyzer.IsDeployPlanDirty()
	if err != nil {
		return err
	}

	if !isPlanDirty {
		// we need to convert the terraformRelativePlanFilePath to a path resolvable from where we are running at the
		//moment (e.g. cmd/example-cli).
		planFullFilePath := filepath.Join(tf.terraformDirectory, terraformRelativePlanFilePath)

		noChangesMarker, err := os.Create(planFullFilePath + ".plan-has-no-changes")
		defer noChangesMarker.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
