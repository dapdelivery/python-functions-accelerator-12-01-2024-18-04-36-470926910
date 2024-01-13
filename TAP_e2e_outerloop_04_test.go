//go:build all || multicluster_outerloop || multicluster_outerloop_gitops
// +build all multicluster_outerloop multicluster_outerloop_gitops

package multicluster_outerloop_test

import (
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/models"
)

func TestOuterloopBasicSupplychainGitopsDelivery(t *testing.T) {
	t.Log("************** TestCase START: TestOuterloopBasicSupplychainGitopsDelivery **************")

	currentSuiteConfig := models.GetSuiteConfig()
	appURLGitOpsDelivery := "http://" + outerloopConfig.Workload.Name + "." + outerloopConfig.Namespace + ".tap." + currentSuiteConfig.Multicluster.RunClusterName + "." + suiteConfig.Domain

	testenv.Test(t,
		common_features.CreateGithubRepo(t, outerloopConfig.Project.Name, outerloopConfig.Project.RepoTemplate, outerloopConfig.Project.AccessToken),
		common_features.CreateGithubRepo(t, outerloopConfig.Project.DestName, outerloopConfig.Project.DestRepoTemplate, outerloopConfig.Project.AccessToken),

		//build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.UpdateTapProfileGitopsSsh(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, "build", "basic", "git-ssh", suiteConfig.Tap.Namespace, suiteConfig.Tap.PollTimeout),
		common_features.ApplyKubectlConfigurationFile(t, outerloopConfig.Workload.GitSSHSecretYamlFile, outerloopConfig.Namespace),
		common_features.PatchServiceAccountSecrets(t, "default", outerloopConfig.Namespace, []string{"git-ssh"}, []string{}),
		common_features.TanzuDeployWorkload(t, outerloopConfig.Workload.GitopsYamlFile, outerloopConfig.Namespace),
		common_features.VerifyGitRepoStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyImageskpac(t, outerloopConfig.Namespace),
		common_features.VerifyBuildStatus(t, outerloopConfig.Workload.Name, suiteConfig.Innerloop.Workload.BuildNameSuffix, outerloopConfig.Namespace),
		common_features.VerifyPodIntentStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyTanzuWorkloadStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyTaskRunStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Workload.TaskRunInfix, outerloopConfig.Namespace),

		//run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.ApplyKubectlConfigurationFile(t, outerloopConfig.Mysql.YamlFile, outerloopConfig.Namespace),
		common_features.ApplyKubectlConfigurationFile(t, outerloopConfig.Workload.GitSSHSecretYamlFile, outerloopConfig.Namespace),

		//copying deliverable from build to run context
		common_features.ProcessDeliverable(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext, ""),

		//run context
		common_features.VerifyRevisionStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyKsvcStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyServiceBindingsStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Workload.ServiceBindingSuffix, outerloopConfig.Namespace),
		common_features.VerifyAppResponseUsingCurl(t, appURLGitOpsDelivery, "/vets.html", outerloopConfig.Project.OriginalString),

		//build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.UpdateGitRepository(t, outerloopConfig.Project.Username, outerloopConfig.Project.Email, outerloopConfig.Project.Repository, outerloopConfig.Project.Name, outerloopConfig.Project.AccessToken, outerloopConfig.Project.File, outerloopConfig.Project.OriginalString, outerloopConfig.Project.NewString, outerloopConfig.Project.CommitMessage),
		common_features.VerifyBuildStatusAfterUpdate(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),

		//run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.VerifyRevisionStatusAfterUpdate(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyKsvcStatusAfterUpdate(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyAppResponseUsingCurl(t, appURLGitOpsDelivery, "/vets.html", outerloopConfig.Project.NewString),

		common_features.DeleteGithubRepo(t, outerloopConfig.Project.Name, outerloopConfig.Project.AccessToken),
		common_features.DeleteGithubRepo(t, outerloopConfig.Project.DestName, outerloopConfig.Project.AccessToken),
		common_features.MulticlusterOuterloopCleanup(t, outerloopConfig.Workload.Name, outerloopConfig.Project.Name, outerloopConfig.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext),
	)

	t.Log("************** TestCase END: TestOuterloopBasicSupplychainGitopsDelivery **************")
}
