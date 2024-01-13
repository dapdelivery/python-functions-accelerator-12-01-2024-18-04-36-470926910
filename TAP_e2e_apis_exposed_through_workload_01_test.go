//go:build all || multicluster_outerloop || multicluster_outerloop_apis
// +build all multicluster_outerloop multicluster_outerloop_apis

package multicluster_outerloop_test

import (
	"fmt"
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/models"
)

func TestOuterloopApisRegistrationTest(t *testing.T) {
	serviceBindingName := fmt.Sprintf("%s-%s%s", outerloopConfig.ApiRegistrationWorkload.Name, outerloopConfig.Workload.Name, outerloopConfig.Workload.ServiceBindingSuffix)
	t.Log("************** TestCase START: TestOuterloopApisRegistrationTest **************")
	currentSuiteConfig := models.GetSuiteConfig()
	appUrlApiReg := "http://" + outerloopConfig.ApiRegistrationWorkload.Name + "." + outerloopConfig.Namespace + ".tap." + currentSuiteConfig.Multicluster.RunClusterName + "." + suiteConfig.Domain

	testenv.Test(t,
		common_features.CreateGithubRepo(t, outerloopConfig.Project.Name, outerloopConfig.Project.RepoTemplate, outerloopConfig.Project.AccessToken),

		//view context
		common_features.ChangeContext(t, suiteConfig.Multicluster.ViewClusterContext),
		common_features.UpdateDomainRecords(t, suiteConfig.Multicluster.ViewClusterName, suiteConfig.Bind9DNSServer.UpdateDnsFilepath, suiteConfig.Bind9DNSServer.HostIP, suiteConfig.Bind9DNSServer.Password),

		//run  context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.UpdateDomainRecords(t, suiteConfig.Multicluster.RunClusterName, suiteConfig.Bind9DNSServer.UpdateDnsFilepath, suiteConfig.Bind9DNSServer.HostIP, suiteConfig.Bind9DNSServer.Password),

		//build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.UpdateMetadataStoreScanning(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, "build", "basic", suiteConfig.Tap.Namespace, suiteConfig.Tap.PollTimeout, outerloopConfig.MetadataStore.Domain, suiteConfig.Multicluster.ViewClusterContext, suiteConfig.Multicluster.BuildClusterContext, outerloopConfig.MetadataStore.Namespace),
		common_features.TanzuDeployWorkload(t, outerloopConfig.ApiRegistrationWorkload.YamlFile2, outerloopConfig.Namespace),
		common_features.VerifyGitRepoStatus(t, outerloopConfig.ApiRegistrationWorkload.Name, outerloopConfig.Namespace),
		common_features.VerifyBuildStatus(t, outerloopConfig.ApiRegistrationWorkload.Name, suiteConfig.Innerloop.Workload.BuildNameSuffix, outerloopConfig.Namespace),
		common_features.VerifyPodIntentStatus(t, outerloopConfig.ApiRegistrationWorkload.Name, outerloopConfig.Namespace),
		common_features.VerifyTanzuWorkloadStatus(t, outerloopConfig.ApiRegistrationWorkload.Name, outerloopConfig.Namespace),
		common_features.VerifyTaskRunStatus(t, outerloopConfig.ApiRegistrationWorkload.Name, outerloopConfig.Workload.TaskRunInfix, outerloopConfig.Namespace),
		//run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.ApplyKubectlConfigurationFile(t, outerloopConfig.Mysql.YamlFile, outerloopConfig.Namespace),

		//copying deliverable from build to run context
		common_features.ProcessDeliverable(t, outerloopConfig.ApiRegistrationWorkload.Name, outerloopConfig.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext, ""),

		//run context
		common_features.VerifyRevisionStatus(t, outerloopConfig.ApiRegistrationWorkload.Name, outerloopConfig.Namespace),
		common_features.VerifyKsvcStatus(t, outerloopConfig.ApiRegistrationWorkload.Name, outerloopConfig.Namespace),
		common_features.VerifyServiceBindingsStatusPreformatted(t, serviceBindingName, outerloopConfig.Namespace),
		common_features.VerifyAppResponseUsingCurl(t, appUrlApiReg, "/"+outerloopConfig.Project.WebpageRelativePath, outerloopConfig.Project.OriginalString),
		common_features.VerifyApiDescriptor(t, outerloopConfig.ApiRegistrationWorkload.Name, outerloopConfig.Namespace, suiteConfig.Multicluster.ViewClusterName, suiteConfig.Multicluster.RunClusterName, suiteConfig.Domain),
		common_features.MulticlusterOuterloopCleanup(t, outerloopConfig.ApiRegistrationWorkload.Name, outerloopConfig.Project.Name, outerloopConfig.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext),

		//build context
		common_features.TanzuDeployWorkload(t, outerloopConfig.ApiRegistrationWorkload.YamlFile, outerloopConfig.Namespace),
		common_features.VerifyGitRepoStatus(t, outerloopConfig.ApiRegistrationWorkload.Name, outerloopConfig.Namespace),
		common_features.VerifyBuildStatus(t, outerloopConfig.ApiRegistrationWorkload.Name, suiteConfig.Innerloop.Workload.BuildNameSuffix, outerloopConfig.Namespace),
		common_features.VerifyPodIntentStatus(t, outerloopConfig.ApiRegistrationWorkload.Name, outerloopConfig.Namespace),
		common_features.VerifyTanzuWorkloadStatus(t, outerloopConfig.ApiRegistrationWorkload.Name, outerloopConfig.Namespace),
		common_features.VerifyTaskRunStatus(t, outerloopConfig.ApiRegistrationWorkload.Name, outerloopConfig.Workload.TaskRunInfix, outerloopConfig.Namespace),
		//run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.ApplyKubectlConfigurationFile(t, outerloopConfig.Mysql.YamlFile, outerloopConfig.Namespace),

		//copying deliverable from build to run context
		common_features.ProcessDeliverable(t, outerloopConfig.ApiRegistrationWorkload.Name, outerloopConfig.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext, ""),

		//run context
		common_features.VerifyRevisionStatus(t, outerloopConfig.ApiRegistrationWorkload.Name, outerloopConfig.Namespace),
		common_features.VerifyKsvcStatus(t, outerloopConfig.ApiRegistrationWorkload.Name, outerloopConfig.Namespace),
		common_features.VerifyServiceBindingsStatusPreformatted(t, serviceBindingName, outerloopConfig.Namespace),
		common_features.VerifyAppResponseUsingCurl(t, appUrlApiReg, "/"+outerloopConfig.Project.WebpageRelativePath, outerloopConfig.Project.OriginalString),
		common_features.VerifyApiDescriptor(t, outerloopConfig.ApiRegistrationWorkload.Name, outerloopConfig.Namespace, suiteConfig.Multicluster.ViewClusterName, suiteConfig.Multicluster.RunClusterName, suiteConfig.Domain),
		common_features.MulticlusterOuterloopCleanup(t, outerloopConfig.ApiRegistrationWorkload.Name, outerloopConfig.Project.Name, outerloopConfig.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext),
		common_features.DeleteGithubRepo(t, outerloopConfig.Project.Name, outerloopConfig.Project.AccessToken),
	)

	t.Log("************** TestCase END: TestOuterloopApisRegistrationTest **************")
}
