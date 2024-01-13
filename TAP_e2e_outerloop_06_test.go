//go:build all || multicluster_outerloop_scan_multiapps_functions
// +build all multicluster_outerloop_scan_multiapps_functions

package multicluster_outerloop_test

import (
	"path/filepath"
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/models"
)

func TestOuterloopScanSupplychainMultipleAppsWithFunctionsWorkload(t *testing.T) {
	t.Log("************** TestCase START: TestOuterloopScanSupplychainMultipleAppsWithFunctionsWorkload **************")
	projectDir := filepath.Join(utils.GetFileDir(), "../../", suiteConfig.Accelerators.RepoName)
	currentSuiteConfig := models.GetSuiteConfig()
	appURLPythonSample := suiteConfig.Accelerators.Name + "." + outerloopConfig.Namespace + ".tap." + currentSuiteConfig.Multicluster.RunClusterName + "." + suiteConfig.Domain

	testenv.Test(t,
		// /*********************************************************view cluster*******************************************************************
		common_features.ChangeContext(t, suiteConfig.Multicluster.ViewClusterContext),
		common_features.AddAcceleratorsFromZip(t, suiteConfig.Accelerators.Name, suiteConfig.Accelerators.ZipUrl, suiteConfig.Accelerators.ZipFile, suiteConfig.Accelerators.ZipDir),
		common_features.UpdateTapValuesAccelerator(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, "view", "testing_scanning", suiteConfig.Tap.Namespace, "LoadBalancer"),
		common_features.GenerateAcceleratorProject(t, suiteConfig.Tap.Namespace, suiteConfig.Accelerators.Name, suiteConfig.Accelerators.RepoName),
		common_features.CreateGithubAcceleratorRepo(t, suiteConfig.Accelerators.GitRepository, suiteConfig.Accelerators.RepoName, outerloopConfig.Project.AccessToken, outerloopConfig.Project.Username, outerloopConfig.Project.Email),

		// /*********************************************************build cluster*******************************************************************
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.TanzuCreateWorkloadWithGitSource(t, suiteConfig.Accelerators.Name, suiteConfig.Accelerators.GitRepository, "master", suiteConfig.Accelerators.BuildEnv, outerloopConfig.Namespace),
		common_features.VerifyTanzuWorkloadStatus(t, suiteConfig.Accelerators.Name, outerloopConfig.Namespace),
		common_features.VerifyBuildStatus(t, suiteConfig.Accelerators.Name, suiteConfig.Innerloop.Workload.BuildNameSuffix, outerloopConfig.Namespace),
		common_features.VerifyTaskRunStatus(t, suiteConfig.Accelerators.Name, outerloopConfig.Workload.TaskRunInfix, outerloopConfig.Namespace),
		common_features.VerifyPodIntent(t, suiteConfig.Accelerators.Name, outerloopConfig.Namespace),

		//run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.ProcessDeliverable(t, suiteConfig.Accelerators.Name, outerloopConfig.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext, ""),
		common_features.VerifyRevisionStatus(t, suiteConfig.Accelerators.Name, outerloopConfig.Namespace),
		common_features.VerifyKsvcStatus(t, suiteConfig.Accelerators.Name, outerloopConfig.Namespace),
		common_features.VerifyAppResponseUsingCurlPost(t, appURLPythonSample, "", suiteConfig.Accelerators.OriginalString),

		// build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.ReplaceStringInFileGitPush(t, suiteConfig.Accelerators.OriginalString, suiteConfig.Accelerators.NewString, suiteConfig.Accelerators.ApplicationFilePath, suiteConfig.Innerloop.Workload.Name, projectDir),
		common_features.VerifyBuildStatusAfterUpdate(t, suiteConfig.Accelerators.Name, outerloopConfig.Namespace),

		// run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.VerifyRevisionStatusAfterUpdate(t, suiteConfig.Accelerators.Name, outerloopConfig.Namespace),
		common_features.VerifyKsvcStatusAfterUpdate(t, suiteConfig.Accelerators.Name, outerloopConfig.Namespace),
		common_features.VerifyAppResponseUsingCurlPost(t, appURLPythonSample, "", suiteConfig.Accelerators.NewString),

		// cleanup
		common_features.Removedir(t, suiteConfig.Accelerators.RepoName),
		common_features.DeleteGithubRepo(t, suiteConfig.Accelerators.RepoName, outerloopConfig.Project.AccessToken),
		common_features.MulticlusterOuterloopCleanup(t, suiteConfig.Accelerators.Name, outerloopConfig.Project.Name, outerloopConfig.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext),
	)
	t.Log("************** TestCase END: TestOuterloopScanSupplychainMultipleAppsWithFunctionsWorkload **************")

}
