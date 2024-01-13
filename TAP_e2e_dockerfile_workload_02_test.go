//go:build all || multicluster_outerloop || multicluster_dockerimage_mavenbuild
// +build all multicluster_outerloop multicluster_dockerimage_mavenbuild

package multicluster_outerloop_test

import (
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/models"
)

func TestDockerBasedWorkloadWithMavenBuildImage(t *testing.T) {
	t.Log("************** TestCase START: TestDockerBasedWorkloadWithMavenBuildImage **************")
	currentSuiteConfig := models.GetSuiteConfig()
	wURL := "http://" + suiteConfig.Innerloop.Workload.Name + "." + outerloopConfig.Namespace + ".tap." + currentSuiteConfig.Multicluster.RunClusterName + "." + suiteConfig.Domain

	testenv.Test(t,
		//build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),

		common_features.GitClone(t, suiteConfig.GitCredentials.Username, suiteConfig.GitCredentials.Email, suiteConfig.Innerloop.Workload.Gitrepository),
		common_features.DockerServiceStart(t),
		common_features.BuildImagewithMaven(t, suiteConfig.Innerloop.Workload.Name, outerloopConfig.DockerFileBasedWorkload.MavenBuildCommand),
		common_features.DockerImageTag(t, outerloopConfig.DockerFileBasedWorkload.Image, outerloopConfig.DockerFileBasedWorkload.DestImageTag),
		common_features.DockerPushImage(t, outerloopConfig.DockerFileBasedWorkload.DestImageTag),

		common_features.TanzuCreateWorkloadWithBuiltImagePath(t, suiteConfig.Innerloop.Workload.Name, outerloopConfig.DockerFileBasedWorkload.DestImageTag, outerloopConfig.Namespace),
		common_features.VerifyPodIntent(t, suiteConfig.Innerloop.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyTanzuWorkloadStatus(t, suiteConfig.Innerloop.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyTaskRunStatus(t, suiteConfig.Innerloop.Workload.Name, outerloopConfig.Workload.TaskRunInfix, outerloopConfig.Namespace),

		//run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),

		//copying deliverable from build to run context
		common_features.ProcessDeliverable(t, suiteConfig.Innerloop.Workload.Name, outerloopConfig.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext, ""),

		//run context
		common_features.VerifyRevisionStatus(t, suiteConfig.Innerloop.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyKsvcStatus(t, suiteConfig.Innerloop.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyAppResponseUsingCurl(t, wURL, "", suiteConfig.Innerloop.Workload.TanzuAppsString),

		common_features.MulticlusterOuterloopCleanup(t, suiteConfig.Innerloop.Workload.Name, "", outerloopConfig.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext),
		common_features.Removedir(t, suiteConfig.Innerloop.Workload.Name),
	)
	t.Log("************** TestCase END: TestDockerBasedWorkloadWithMavenBuildImage **************")
}
