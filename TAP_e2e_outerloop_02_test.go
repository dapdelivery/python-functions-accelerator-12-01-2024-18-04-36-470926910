//go:build all || multicluster_outerloop || multicluster_outerloop_test || multicluster_airgap
// +build all multicluster_outerloop multicluster_outerloop_test multicluster_airgap

package multicluster_outerloop_test

import (
	"fmt"
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
)

func TestOuterloopTestSupplychainGitSource(t *testing.T) {
	t.Log("************** TestCase START: TestOuterloopTestSupplychainGitSource **************")
	testenv.Test(t,
		common_features.CreateGithubRepo(t, whereForDinnerConfig.Project.Name, whereForDinnerConfig.Project.RepoTemplate, whereForDinnerConfig.Project.AccessToken),

		//build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.UpdateTapProfileSupplyChain(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, "build", "testing", suiteConfig.Tap.Namespace, suiteConfig.Tap.PollTimeout),

		common_features.VerifyPipelineStatus(t, whereForDinnerConfig.AppTestsPipeline.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.TanzuDeployListOfWorkloads(t, whereForDinnerConfig.Workload.TestYamlFilePath, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfTanzuWorkloadStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace, whereForDinnerConfig.Workload.TestYamlFilePath),
		common_features.VerifyListOfGitRepoStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListPipelineRunStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfBuildStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.BuildNameSuffix, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyPodIntentStatus(t, whereForDinnerConfig.Workload.Name[1], whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfTaskRunStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.TaskRunInfix, whereForDinnerConfig.Workload.Namespace),

		// run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),

		// copying deliverable from build to run context
		common_features.ProcessListOfDeliverables(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext, "", ""),

		//run context
		common_features.VerifyListOfRevisionStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfKsvcStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyServiceBindingsStatusPreformatted(t, fmt.Sprintf("%s%s", whereForDinnerConfig.Workload.Name[1], whereForDinnerConfig.Workload.ServiceBindingSuffix), whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyAppResponseUsingCurl(t, appURL, whereForDinnerConfig.Workload.EndPoint, whereForDinnerConfig.Workload.OriginalVerificationString),

		//build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.UpdateGitRepository(t, whereForDinnerConfig.Project.Username, whereForDinnerConfig.Project.Email, whereForDinnerConfig.Project.Repository, whereForDinnerConfig.Project.Name, whereForDinnerConfig.Project.AccessToken, whereForDinnerConfig.Workload.ApplicationFilePath, whereForDinnerConfig.Workload.OriginalString, whereForDinnerConfig.Workload.NewString, whereForDinnerConfig.Project.CommitMessage),
		common_features.VerifyBuildStatusAfterUpdate(t, whereForDinnerConfig.Workload.Name[1], whereForDinnerConfig.Workload.Namespace),

		//run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.VerifyRevisionStatusAfterUpdate(t, whereForDinnerConfig.Workload.Name[1], whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyKsvcStatusAfterUpdate(t, whereForDinnerConfig.Workload.Name[1], whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyAppResponseUsingCurl(t, appURL, whereForDinnerConfig.Workload.EndPoint, whereForDinnerConfig.Workload.NewVerificationString),

		common_features.DeleteGithubRepo(t, whereForDinnerConfig.Project.Name, whereForDinnerConfig.Project.AccessToken),
		//run cluster cleaup
		common_features.DeliverableCleanup(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		//build cluster cleanup
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.DeleteWorkloadAll(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.DeletePipeline(t, whereForDinnerConfig.AppTestsPipeline.Name, whereForDinnerConfig.Workload.Namespace),
	)

	t.Log("************** TestCase END: TestOuterloopTestSupplychainGitSource **************")
}
