//go:build all || multicluster_outerloop || multicluster_outerloop_snyk_scan
// +build all multicluster_outerloop multicluster_outerloop_snyk_scan

package multicluster_outerloop_test

import (
	"fmt"
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
)

func TestOuterloopWithSnykScanSupplychainGitSource(t *testing.T) {
	t.Log("************** TestCase START: TestOuterloopWithSnykScanSupplychainGitSource **************")
	testenv.Test(t,
		common_features.ChangeContext(t, suiteConfig.Multicluster.ViewClusterContext),
		common_features.UpdateDomainRecords(t, suiteConfig.Multicluster.ViewClusterName, suiteConfig.Bind9DNSServer.UpdateDnsFilepath, suiteConfig.Bind9DNSServer.HostIP, suiteConfig.Bind9DNSServer.Password),
		common_features.CreateGithubRepo(t, whereForDinnerConfig.Project.Name, whereForDinnerConfig.Project.RepoTemplate, whereForDinnerConfig.Project.AccessToken),

		//build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.ApplyKubectlConfigurationFile(t, whereForDinnerConfig.SnykSecret.YamlFile, whereForDinnerConfig.Workload.Namespace),
		common_features.UpdateMetadataStoreAndInstallSnykScanning(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, "build", "testing_scanning", suiteConfig.Tap.Namespace, suiteConfig.Tap.PollTimeout, whereForDinnerConfig.MetadataStore.Domain, suiteConfig.Multicluster.ViewClusterContext, suiteConfig.Multicluster.BuildClusterContext, whereForDinnerConfig.MetadataStore.Namespace, whereForDinnerConfig.SnykInstallValues.YamlFile),
		common_features.VerifyPipelineStatus(t, whereForDinnerConfig.AppTestsPipeline.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.ApplyKubectlConfigurationFile(t, whereForDinnerConfig.SnykScanPolicy.YamlFile, whereForDinnerConfig.Workload.Namespace),
		common_features.TanzuDeployListOfWorkloads(t, whereForDinnerConfig.Workload.TestYamlFilePath, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfTanzuWorkloadStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace, whereForDinnerConfig.Workload.TestYamlFilePath),
		common_features.VerifyListOfGitRepoStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfSourceScanStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListPipelineRunStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfBuildStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.BuildNameSuffix, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyPodIntentStatus(t, whereForDinnerConfig.Workload.Name[1], whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfImageScanStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfTestTaskRunStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.TaskRunTestSuffix, whereForDinnerConfig.Workload.Namespace),

		// //run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),

		// //copying deliverable from build to run context
		common_features.ProcessListOfDeliverables(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext, "", ""),

		// //run context
		common_features.VerifyListOfRevisionStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfKsvcStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyServiceBindingsStatusPreformatted(t, fmt.Sprintf("%s%s", whereForDinnerConfig.Workload.Name[1], whereForDinnerConfig.Workload.ServiceBindingSuffix), whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyAppResponseUsingCurl(t, appURL, whereForDinnerConfig.Workload.EndPoint, whereForDinnerConfig.Workload.OriginalVerificationString),

		// //build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.UpdateGitRepository(t, whereForDinnerConfig.Project.Username, whereForDinnerConfig.Project.Email, whereForDinnerConfig.Project.Repository, whereForDinnerConfig.Project.Name, whereForDinnerConfig.Project.AccessToken, whereForDinnerConfig.Workload.ApplicationFilePath, whereForDinnerConfig.Workload.OriginalString, whereForDinnerConfig.Workload.NewString, whereForDinnerConfig.Project.CommitMessage),
		common_features.VerifyBuildStatusAfterUpdate(t, whereForDinnerConfig.Workload.Name[1], whereForDinnerConfig.Workload.Namespace),

		//run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.VerifyRevisionStatusAfterUpdate(t, whereForDinnerConfig.Workload.Name[1], whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyKsvcStatusAfterUpdate(t, whereForDinnerConfig.Workload.Name[1], whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyAppResponseUsingCurl(t, appURL, whereForDinnerConfig.Workload.EndPoint, whereForDinnerConfig.Workload.NewVerificationString),

		//run cluster cleaup
		common_features.DeleteGithubRepo(t, whereForDinnerConfig.Project.Name, whereForDinnerConfig.Project.AccessToken),
		common_features.DeliverableCleanup(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		//build cluster cleanup
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.DeleteNamespace(t, "metadata-store-secrets", ""),
		common_features.DeleteWorkloadAll(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.DeletePipeline(t, whereForDinnerConfig.AppTestsPipeline.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.DeletePackage(t, "snyk-scanner", suiteConfig.Tap.Namespace),
	)
	t.Log("************** TestCase END: TestOuterloopWithSnykScanSupplychainGitSource **************")
}
