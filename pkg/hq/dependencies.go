package hq

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	semver "github.com/Masterminds/semver/v3"
	"github.com/conplementag/cops-hq/v2/internal"
	"github.com/conplementag/cops-hq/v2/pkg/error_handling"
	"github.com/sirupsen/logrus"
)

// expected min versions should be maintained individual dependent to the dependency
// azurecli -> can be set to latest version
// terraform -> releases should be checked for breaking changes. latest version not required. settled release should be preferred
// helm -> releases should be checked for breaking changes. latest version not required. settled release should be preferred
// kubectl -> should follow the guidelines of kubectl. You must use a kubectl version that is within one minor version difference of your cops cluster. For example, a v1.24 client can communicate with v1.23, v1.24, and v1.25 control planes.
// kubelogin -> can be set to latest version
// copsctl -> can be set to latest version
// sops -> can be set to latest version
const (
	ExpectedMinAzureCliVersion  = "2.81.0" // 01/2026
	ExpectedMinTerraformVersion = "1.13.0" // 01/2026 - allow previous version
	ExpectedMinHelmVersion      = "3.18.0" // 01/2026 - allow previous version
	ExpectedMinKubectlVersion   = "1.33.0" // 01/2026
	ExpectedMinKubeloginVersion = "0.2.12" // 01/2026
	ExpectedMinCopsctlVersion   = "0.16.0" // 01/2026
	ExpectedMinSopsVersion      = "3.11.0" // 01/2026
)

func (hq *hqContainer) CheckToolingDependencies() error {
	logrus.Info("Checking tooling dependencies...")

	// mandatory dependencies
	err1 := hq.checkAzureCli()
	err2 := hq.checkHelm()
	err3 := hq.checkTerraform()
	err4 := hq.checkKubectl()
	err5 := hq.checkKubelogin()
	err6 := hq.checkCopsctl()

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil {
		compositeErr := fmt.Errorf("mandatory tooling dependencies check failed: %v %v %v %v %v %v", err1, err2, err3, err4, err5, err6)
		return internal.ReturnErrorOrPanic(compositeErr)
	}

	// optional but recommended dependencies
	warn1 := hq.checkSops()

	if warn1 != nil {
		logrus.Warnf("Sops - optional dependency (recommended to be installed) not met: %v", warn1)
		logrus.Warn("Sops is a useful tool for source version configuration management.")
	}

	warn2 := hq.checkVim()

	if warn2 != nil {
		logrus.Warnf("Vim - optional dependency (recommended to be installed) not met: %v", warn2)
		logrus.Warn("Vim is used as the default editor for some cops-hq functionality, like fixing MAC versions " +
			"of Sops managed config files.")
	}

	return nil
}

func (hq *hqContainer) checkAzureCli() error {
	logrus.Info("Checking azure cli...")
	azureCliVersion, err := hq.Executor.Execute("az version -o json")

	if err != nil {
		return err
	}

	var response azureCliVersionResponse
	err = json.Unmarshal([]byte(azureCliVersion), &response)

	if err != nil {
		return err
	}

	versionConstraint, _ := semver.NewConstraint(">=" + ExpectedMinAzureCliVersion)
	installedVersion, _ := semver.NewVersion(response.AzureCli)

	if !versionConstraint.Check(installedVersion) {
		return fmt.Errorf("azure cli version mismatch. expected >= %v, got %v", ExpectedMinAzureCliVersion, installedVersion)
	}

	logrus.Info("...ok.")
	return nil
}

func (hq *hqContainer) checkHelm() error {
	logrus.Info("Checking helm...")
	helmVersion, err := hq.Executor.Execute("helm version --template={{.Version}}")

	if err != nil {
		return err
	}

	helmVersion = strings.TrimSuffix(helmVersion, "%") // some systems add % to the output

	versionConstraint, _ := semver.NewConstraint(">=" + ExpectedMinHelmVersion)
	installedVersion, _ := semver.NewVersion(helmVersion)

	if !versionConstraint.Check(installedVersion) {
		return fmt.Errorf("helm version mismatch. expected >= %v, got %v", ExpectedMinHelmVersion, installedVersion)
	}

	logrus.Info("...ok.")
	return nil
}

func (hq *hqContainer) checkTerraform() error {
	logrus.Info("Checking terraform...")
	terraformVersion, err := hq.Executor.Execute("terraform --version -json")

	if err != nil {
		return err
	}

	var terraformResponse terraformVersionResponse
	err = json.Unmarshal([]byte(terraformVersion), &terraformResponse)

	if err != nil {
		return err
	}

	versionConstraint, _ := semver.NewConstraint(">=" + ExpectedMinTerraformVersion)
	installedVersion, _ := semver.NewVersion(terraformResponse.TerraformVersion)

	if !versionConstraint.Check(installedVersion) {
		return fmt.Errorf("terraform version mismatch. expected >= %v, got %v", ExpectedMinTerraformVersion, installedVersion)
	}

	logrus.Info("...ok.")
	return nil
}

func (hq *hqContainer) checkKubectl() error {
	logrus.Info("Checking kubectl...")
	kubectlVersion, err := hq.Executor.Execute("kubectl version --client=true -o json")

	if err != nil {
		return err
	}

	var kubectlResponse kubectlVersionResponse
	err = json.Unmarshal([]byte(kubectlVersion), &kubectlResponse)

	if err != nil {
		return err
	}

	versionConstraint, _ := semver.NewConstraint(">=" + ExpectedMinKubectlVersion)
	installedVersion, err := semver.NewVersion(kubectlResponse.ClientVersion.GitVersion)

	if !versionConstraint.Check(installedVersion) {
		return fmt.Errorf("kubectl version mismatch. expected >= %v, got %v", ExpectedMinKubectlVersion, installedVersion)
	}

	logrus.Info("...ok.")
	return nil
}

func (hq *hqContainer) checkKubelogin() error {
	logrus.Info("Checking kubelogin...")
	kubeloginVersion, err := hq.Executor.Execute("kubelogin --version")

	if err != nil {
		return err
	}

	kubeloginRegex, _ := regexp.Compile(".*(v\\d+\\.\\d+\\.\\d+).*")

	if kubeloginRegex.MatchString(kubeloginVersion) {
		matches := kubeloginRegex.FindStringSubmatch(kubeloginVersion)
		versionConstraint, _ := semver.NewConstraint(">=" + ExpectedMinKubeloginVersion)
		installedVersion, _ := semver.NewVersion(matches[1])

		if !versionConstraint.Check(installedVersion) {
			return fmt.Errorf("kubelogin version mismatch. expected >= %v, got %v", ExpectedMinKubeloginVersion, installedVersion)
		}
	} else {
		return fmt.Errorf("kubelogin version could not be parsed from this output: %s", kubeloginVersion)
	}

	logrus.Info("...ok.")
	return nil
}

func (hq *hqContainer) checkCopsctl() error {
	logrus.Info("Checking copsctl...")
	copsctlVersion, err := hq.Executor.Execute("copsctl --version")

	if err != nil {
		return err
	}

	copsctlVersion = strings.TrimPrefix(copsctlVersion, "copsctl version ")

	versionConstraint, _ := semver.NewConstraint(">=" + ExpectedMinCopsctlVersion)
	installedVersion, _ := semver.NewVersion(copsctlVersion)

	if !versionConstraint.Check(installedVersion) {
		return fmt.Errorf("copsctl version mismatch. expected >= %v, got %v", ExpectedMinCopsctlVersion, installedVersion)
	}

	logrus.Info("...ok.")
	return nil
}

func (hq *hqContainer) checkSops() error {
	logrus.Info("Checking sops...")

	// sops is an optional dependency, so in case we are in panic mode, we should survive it
	previousPanicSetting := error_handling.PanicOnAnyError
	error_handling.PanicOnAnyError = false

	sopsVersion, err := hq.Executor.Execute("sops --version")

	error_handling.PanicOnAnyError = previousPanicSetting

	if err != nil {
		return err
	}

	sopsRegex, _ := regexp.Compile(".*(\\d+\\.\\d+\\.\\d+).*")

	if sopsRegex.MatchString(sopsVersion) {
		matches := sopsRegex.FindStringSubmatch(sopsVersion)
		versionConstraint, _ := semver.NewConstraint(">=" + ExpectedMinSopsVersion)
		installedVersion, _ := semver.NewVersion(matches[1])

		if installedVersion == nil || !versionConstraint.Check(installedVersion) {
			return fmt.Errorf("sops version mismatch. expected >= %v, got %v", ExpectedMinSopsVersion, installedVersion)
		}
	} else {
		return fmt.Errorf("sops version could not be parsed from this output: %s", sopsVersion)
	}

	logrus.Info("...ok.")
	return nil
}

func (hq *hqContainer) checkVim() error {
	logrus.Info("Checking vim...")

	// Vim is an optional dependency, so in case we are in panic mode, we should survive it
	previousPanicSetting := error_handling.PanicOnAnyError
	error_handling.PanicOnAnyError = false

	// Result is ignored, because we simply need to check if installed, which should return no errors.
	// Checking for correct version like for other dependencies is not required here.
	_, err := hq.Executor.Execute("vim --version")

	error_handling.PanicOnAnyError = previousPanicSetting

	if err != nil {
		logrus.Info("...ok.")
	}

	return err
}

type azureCliVersionResponse struct {
	AzureCli string `json:"azure-cli"`
}

type terraformVersionResponse struct {
	TerraformVersion string `json:"terraform_version"`
}

type kubectlVersionResponse struct {
	ClientVersion struct {
		GitVersion string `json:"gitVersion"`
	} `json:"clientVersion"`
}
