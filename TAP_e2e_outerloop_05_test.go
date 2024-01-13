//go:build all || multicluster_outerloop || multicluster_outerloop_basic_maven_build
// +build all multicluster_outerloop multicluster_outerloop_basic_maven_build

package multicluster_outerloop_test

import (
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/models"
)

// Sumit - This case adaptation as of it is not possible. All where-for-dinner workload don't have mvn build system
func TestOuterloopBasicSupplyMavenBuild(t *testing.T) {
	t.Log("************** TestCase START: TestOuterloopBasicSupplyMavenBuild **************")
	currentSuiteConfig := models.GetSuiteConfig()
	appURLMaven := "http://" + outerloopConfig.Workload.Name + "." + outerloopConfig.Namespace + ".tap." + currentSuiteConfig.Multicluster.RunClusterName + "." + suiteConfig.Domain

	testenv.Test(t,
		common_features.CreateGithubRepo(t, outerloopConfig.Project.Name, outerloopConfig.Project.RepoTemplate, outerloopConfig.Project.AccessToken),
		//build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.ApplyKubectlConfigurationFile(t, outerloopConfig.Workload.GitBasicSecretYamlFile, outerloopConfig.Namespace),
		common_features.ApplyKubectlConfigurationFile(t, outerloopConfig.Workload.GitBasicSecretYamlFile, suiteConfig.PackageRepository.Namespace),
		common_features.UpdateTapValuesWithMaven(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, "build", "basic", suiteConfig.Tap.Namespace, suiteConfig.Tap.PollTimeout, suiteConfig.Multicluster.ViewClusterContext, suiteConfig.Multicluster.BuildClusterContext, outerloopConfig.Project.MavenInstallLink),
		common_features.CreateMavenBuild(t, outerloopConfig.Project.Name, outerloopConfig.Project.Repository),
		common_features.CreateMavenBuildAndPush(t, outerloopConfig.Project.Name, outerloopConfig.Project.Repository, outerloopConfig.MavenBuild.GroupID, outerloopConfig.MavenBuild.ArtifactID, outerloopConfig.MavenBuild.Version, outerloopConfig.Project.Username, outerloopConfig.Project.Email, outerloopConfig.Project.AccessToken, outerloopConfig.Project.MavenBranch),
		common_features.TanzuDeployWorkload(t, outerloopConfig.Workload.MavenYamlFile, outerloopConfig.Namespace),
		common_features.VerifyPodIntentStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyTanzuWorkloadStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyBuildStatus(t, outerloopConfig.Workload.Name, suiteConfig.Innerloop.Workload.BuildNameSuffix, outerloopConfig.Namespace),
		common_features.VerifyTaskRunStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Workload.TaskRunInfix, outerloopConfig.Namespace),
		common_features.VerifyMavenArtifcats(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),

		//run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.ApplyKubectlConfigurationFile(t, outerloopConfig.Mysql.YamlFile, outerloopConfig.Namespace),

		//copying deliverable from build to run context
		common_features.ProcessDeliverable(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext, ""),

		//run context
		common_features.VerifyRevisionStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		//common_features.VerifyKsvcStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyListOfKsvcStatus(t, []string{outerloopConfig.Workload.Name}, outerloopConfig.Namespace),
		common_features.VerifyServiceBindingsStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Workload.ServiceBindingSuffix, outerloopConfig.Namespace),
		common_features.VerifyAppResponseUsingCurl(t, appURLMaven, "/vets.html", outerloopConfig.Project.OriginalString),

		//build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.UpdateMavenGitRepository(t, outerloopConfig.Project.Username, outerloopConfig.Project.Email, outerloopConfig.Project.Repository, outerloopConfig.Project.Name, outerloopConfig.Project.AccessToken, outerloopConfig.Project.File, outerloopConfig.Project.OriginalString, outerloopConfig.Project.NewString, outerloopConfig.Project.CommitMessage, outerloopConfig.MavenBuild.GroupID, outerloopConfig.MavenBuild.ArtifactID, outerloopConfig.Project.MavenBranch),
		common_features.VerifyBuildStatusAfterUpdate(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),

		//run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.VerifyRevisionStatusAfterUpdate(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyKsvcStatusAfterUpdate(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyAppResponseUsingCurl(t, appURLMaven, "/vets.html", outerloopConfig.Project.NewString),

		common_features.DeleteGithubRepo(t, outerloopConfig.Project.Name, outerloopConfig.Project.AccessToken),
		common_features.MulticlusterOuterloopCleanup(t, outerloopConfig.Workload.Name, outerloopConfig.Project.Name, outerloopConfig.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext),
	)

	t.Log("************** TestCase END: TestOuterloopBasicSupplyMavenBuild **************")
}
