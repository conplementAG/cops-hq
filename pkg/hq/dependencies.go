package hq

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/sirupsen/logrus"
	"strings"
)

const (
	ExpectedMinAzureCliVersion  = "2.36.0"
	ExpectedMinTerraformVersion = "1.1.9"
	ExpectedMinHelmVersion      = "3.8.2"
	ExpectedMinKubectlVersion   = "1.23.5"
	ExpectedMinCopsctlVersion   = "0.8.0"
	ExpectedMinSopsVersion      = "3.7.3"
)

func (hq *hqContainer) CheckToolingDependencies() error {
	logrus.Info("Checking tooling dependencies...")

	// mandatory dependencies
	err1 := hq.checkAzureCli()
	err2 := hq.checkHelm()
	err3 := hq.checkTerraform()
	err4 := hq.checkKubectl()
	err5 := hq.checkCopsctl()

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		return fmt.Errorf("mandatory tooling dependencies check failed: %v %v %v %v %v", err1, err2, err3, err4, err5)
	}

	// optional but recommended dependencies
	warn1 := hq.checkSops()

	if warn1 != nil {
		logrus.Warnf("optional dependency (recommended to be used) not met: %v", warn1)
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
	sopsVersion, err := hq.Executor.Execute("sops --version")

	if err != nil {
		return err
	}

	sopsVersion = strings.TrimPrefix(sopsVersion, "sops ")
	sopsVersion = strings.TrimSuffix(sopsVersion, " (latest)")

	versionConstraint, _ := semver.NewConstraint(">=" + ExpectedMinSopsVersion)
	installedVersion, _ := semver.NewVersion(sopsVersion)

	if !versionConstraint.Check(installedVersion) {
		return fmt.Errorf("sops version mismatch. expected >= %v, got %v", ExpectedMinSopsVersion, installedVersion)
	}

	logrus.Info("...ok.")
	return nil
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
