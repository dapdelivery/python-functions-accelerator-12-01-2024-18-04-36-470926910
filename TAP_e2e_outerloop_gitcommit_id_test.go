//go:build all || multicluster_outerloop || multicluster_outerloop_workload_with_gitcommit
// +build all multicluster_outerloop multicluster_outerloop_workload_with_gitcommit

package multicluster_outerloop_test

import (
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/models"
)

func TestWorkloadWithGitCommitID(t *testing.T) {
	t.Log("************** TestCase START: TestWorkloadWithGitCommitID **************")
	currentSuiteConfig := models.GetSuiteConfig()
	wURL := "http://" + suiteConfig.Innerloop.Workload.Name + "." + outerloopConfig.Namespace + ".tap." + currentSuiteConfig.Multicluster.RunClusterName + "." + suiteConfig.Domain

	testenv.Test(t,

		common_features.CreateGithubRepo(t, suiteConfig.Innerloop.Workload.RepositoryName, suiteConfig.Innerloop.Workload.RepoTemplate, outerloopConfig.Project.AccessToken),

		// //build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.UpdateTapProfileSupplyChain(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, "build", "testing", suiteConfig.Tap.Namespace, suiteConfig.Tap.PollTimeout),
		common_features.TanzuCreateWorkloadUsingGitCommitId(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Repository, whereForDinnerConfig.Workload.Namespace, suiteConfig.Innerloop.Workload.Branch, "true", "web", suiteConfig.Innerloop.Workload.CommitURL, outerloopConfig.Project.AccessToken, outerloopConfig.Project.Username, outerloopConfig.Project.Email),
		common_features.VerifyListOfTanzuWorkloadStatus(t, []string{suiteConfig.Innerloop.Workload.Name}, whereForDinnerConfig.Workload.Namespace, ""),
		common_features.VerifyListOfGitRepoStatus(t, []string{suiteConfig.Innerloop.Workload.Name}, whereForDinnerConfig.Workload.Namespace),

		common_features.VerifyListOfBuildStatus(t, []string{suiteConfig.Innerloop.Workload.Name}, suiteConfig.Innerloop.Workload.BuildNameSuffix, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyListOfTestTaskRunStatus(t, []string{suiteConfig.Innerloop.Workload.Name}, whereForDinnerConfig.Workload.TaskRunTestSuffix, suiteConfig.Innerloop.Workload.Namespace),

		// //run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),

		// //copying deliverable from build to run context
		common_features.ProcessListOfDeliverables(t, []string{suiteConfig.Innerloop.Workload.Name}, suiteConfig.Innerloop.Workload.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext, "", ""),

		// //run context
		common_features.VerifyListOfRevisionStatus(t, []string{suiteConfig.Innerloop.Workload.Name}, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyListOfKsvcStatus(t, []string{suiteConfig.Innerloop.Workload.Name}, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyAppResponseUsingCurl(t, wURL, "", suiteConfig.Innerloop.Workload.OriginalString),
		// // build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.UpdateGitRepository(t, outerloopConfig.Project.Username, outerloopConfig.Project.Email, suiteConfig.Innerloop.Workload.Repository, suiteConfig.Innerloop.Workload.RepositoryName, outerloopConfig.Project.AccessToken, suiteConfig.Innerloop.Workload.ApplicationFilePath2, suiteConfig.Innerloop.Workload.OriginalString, suiteConfig.Innerloop.Workload.NewString, "update message"),
		common_features.VerifyListOfTanzuWorkloadStatus(t, []string{suiteConfig.Innerloop.Workload.Name}, whereForDinnerConfig.Workload.Namespace, ""),
		common_features.TanzuCreateWorkloadUsingGitCommitId(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Repository, whereForDinnerConfig.Workload.Namespace, suiteConfig.Innerloop.Workload.Branch, "true", "web", suiteConfig.Innerloop.Workload.CommitURL, outerloopConfig.Project.AccessToken, outerloopConfig.Project.Username, outerloopConfig.Project.Email),
		common_features.VerifyListOfTanzuWorkloadStatus(t, []string{suiteConfig.Innerloop.Workload.Name}, whereForDinnerConfig.Workload.Namespace, ""),
		common_features.VerifyBuildStatusAfterUpdate(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),

		// //run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.VerifyRevisionStatusAfterUpdate(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyKsvcStatusAfterUpdate(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyAppResponseUsingCurl(t, wURL, "", suiteConfig.Innerloop.Workload.NewString),

		//run cluster cleanup
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.DeleteGithubRepo(t, suiteConfig.Innerloop.Workload.RepositoryName, outerloopConfig.Project.AccessToken),
		common_features.DeliverableCleanup(t, []string{suiteConfig.Innerloop.Workload.Name}, suiteConfig.Innerloop.Workload.Namespace),

		// //build cluster cleanup
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.DeleteWorkloadAll(t, []string{suiteConfig.Innerloop.Workload.Name}, suiteConfig.Innerloop.Workload.Namespace),
		common_features.UpdateTapProfileSupplyChain(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, "build", "basic", suiteConfig.Tap.Namespace, suiteConfig.Tap.PollTimeout),
	)
	t.Log("************** TestCase END: TestWorkloadWithGitCommitID **************")
}
