//go:build multicluster_outerloop_gitops_ignored
// +build multicluster_outerloop_gitops_ignored

package multicluster_outerloop_test

import (
	"fmt"
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
)

func TestOuterloopBasicSupplychainGitopsDelivery(t *testing.T) {
	t.Log("************** TestCase START: TestOuterloopBasicSupplychainGitopsDelivery **************")
	testenv.Test(t,
		common_features.CreateGithubRepo(t, whereForDinnerConfig.Project.Name, whereForDinnerConfig.Project.RepoTemplate, whereForDinnerConfig.Project.AccessToken),
		common_features.CreateGithubRepo(t, whereForDinnerConfig.Project.DestName, whereForDinnerConfig.Project.RepoTemplate, whereForDinnerConfig.Project.AccessToken),

		//build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.UpdateTapProfileGitopsSsh(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, "full", "basic", "git-ssh", suiteConfig.Tap.Namespace, suiteConfig.Tap.PollTimeout),
		common_features.ApplyKubectlConfigurationFile(t, whereForDinnerConfig.Workload.GitSSHSecretYamlFile, whereForDinnerConfig.Workload.Namespace),
		common_features.PatchServiceAccountSecrets(t, "default", whereForDinnerConfig.Workload.Namespace, []string{"git-ssh"}, []string{}),
		common_features.TanzuDeployListOfWorkloadsGitOps(t, whereForDinnerConfig.Workload.GitOpsYamlFilePath, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfGitRepoStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfImagesKpacStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfBuildStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.BuildNameSuffix, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfPodIntentStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyTanzuWorkloadStatus(t, whereForDinnerConfig.Workload.Name[1], whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyTaskRunStatus(t, whereForDinnerConfig.Workload.Name[1], whereForDinnerConfig.Workload.TaskRunInfix, whereForDinnerConfig.Workload.Namespace),

		//run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.ApplyKubectlConfigurationFile(t, whereForDinnerConfig.Workload.GitSSHSecretYamlFile, whereForDinnerConfig.Workload.Namespace),

		//copying deliverable from build to run context
		common_features.ProcessListOfDeliverables(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext, ""),

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
		common_features.DeleteGithubRepo(t, whereForDinnerConfig.Project.DestName, whereForDinnerConfig.Project.AccessToken),

		//run cluster cleaup
		common_features.DeliverableCleanup(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		//build cluster cleanup
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.DeleteWorkloadAll(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
	)

	t.Log("************** TestCase END: TestOuterloopBasicSupplychainGitopsDelivery **************")
}
