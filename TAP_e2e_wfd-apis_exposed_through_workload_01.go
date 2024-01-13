//go:build multicluster_outerloop_apiswfd
// +build multicluster_outerloop_apiswfd

package multicluster_outerloop_test

import (
	"fmt"
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
)

// Sumit - This is not supported in where-for-dinner as of yet. Support to come is still TBD. Update the go build tag
// Slack disucssion with (Greg Meyer): https://vmware.slack.com/archives/C02D60T1ZDJ/p1679909845325949
func TestOuterloopApisRegistrationTestWFD(t *testing.T) {
	t.Log("************** TestCase START: TestOuterloopApisRegistrationTest **************")
	testenv.Test(t,
		common_features.CreateGithubRepo(t, whereForDinnerConfig.Project.Name, whereForDinnerConfig.Project.RepoTemplate, whereForDinnerConfig.Project.AccessToken),

		//view context
		common_features.ChangeContext(t, suiteConfig.Multicluster.ViewClusterContext),
		common_features.UpdateDomainRecords(t, suiteConfig.Multicluster.ViewClusterName, suiteConfig.Bind9DNSServer.UpdateDnsFilepath, suiteConfig.Bind9DNSServer.HostIP, suiteConfig.Bind9DNSServer.Password),

		//run  context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.UpdateDomainRecords(t, suiteConfig.Multicluster.RunClusterName, suiteConfig.Bind9DNSServer.UpdateDnsFilepath, suiteConfig.Bind9DNSServer.HostIP, suiteConfig.Bind9DNSServer.Password),

		//build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.UpdateMetadataStoreScanning(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, "build", "basic", suiteConfig.Tap.Namespace, suiteConfig.Tap.PollTimeout, whereForDinnerConfig.MetadataStore.Domain, suiteConfig.Multicluster.ViewClusterContext, suiteConfig.Multicluster.BuildClusterContext, whereForDinnerConfig.MetadataStore.Namespace, false),
		common_features.TanzuDeployListOfWorkloads(t, whereForDinnerConfig.Workload.ApiRegistrationYamlFilePath, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfTanzuWorkloadStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace, whereForDinnerConfig.Workload.ApiRegistrationYamlFilePath),
		common_features.VerifyListOfGitRepoStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfBuildStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.BuildNameSuffix, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyPodIntentStatus(t, whereForDinnerConfig.Workload.Name[1], whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfTaskRunStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.TaskRunInfix, whereForDinnerConfig.Workload.Namespace),

		//run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),

		//copying deliverable from build to run context
		common_features.ProcessListOfDeliverables(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext, "", ""),

		//run context
		common_features.UpdatePackage(t, "api-auto-registration", "apis.apps.tanzu.vmware.com", "", suiteConfig.Tap.Namespace, whereForDinnerConfig.ApiRegisrationValuesFile, "30m"),
		common_features.VerifyListOfRevisionStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfKsvcStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyServiceBindingsStatusPreformatted(t, fmt.Sprintf("%s%s", whereForDinnerConfig.Workload.Name[1], whereForDinnerConfig.Workload.ServiceBindingSuffix), whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyAppResponseUsingCurl(t, appURL, whereForDinnerConfig.Workload.EndPoint, whereForDinnerConfig.Workload.OriginalVerificationString),
		common_features.VerifyApiDescriptor(t, whereForDinnerConfig.Workload.Name[0], whereForDinnerConfig.Workload.Namespace, suiteConfig.Multicluster.ViewClusterName, suiteConfig.Multicluster.RunClusterName, suiteConfig.Domain),
		common_features.MulticlusterListOuterloopCleanup(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext),

		//build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.TanzuDeployListOfWorkloads(t, whereForDinnerConfig.Workload.TestYamlFilePath, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfTanzuWorkloadStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace, whereForDinnerConfig.Workload.TestYamlFilePath),
		common_features.VerifyListOfGitRepoStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfBuildStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.BuildNameSuffix, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyPodIntentStatus(t, whereForDinnerConfig.Workload.Name[1], whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfTaskRunStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.TaskRunInfix, whereForDinnerConfig.Workload.Namespace),

		//run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),

		//copying deliverable from build to run context
		common_features.ProcessListOfDeliverables(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext, "", ""),

		//run context
		common_features.UpdatePackage(t, "api-auto-registration", "", "", suiteConfig.Tap.Namespace, whereForDinnerConfig.ApiRegisrationValuesFile, "30m"),
		common_features.VerifyListOfRevisionStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfKsvcStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyServiceBindingsStatusPreformatted(t, fmt.Sprintf("%s%s", whereForDinnerConfig.Workload.Name[1], whereForDinnerConfig.Workload.ServiceBindingSuffix), whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyAppResponseUsingCurl(t, appURL, whereForDinnerConfig.Workload.EndPoint, whereForDinnerConfig.Workload.OriginalVerificationString),
		common_features.VerifyApiDescriptor(t, whereForDinnerConfig.Workload.Name[1], whereForDinnerConfig.Workload.Namespace, suiteConfig.Multicluster.ViewClusterName, suiteConfig.Multicluster.RunClusterName, suiteConfig.Domain),
		common_features.MulticlusterListOuterloopCleanup(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext),

		common_features.DeleteGithubRepo(t, whereForDinnerConfig.Project.Name, whereForDinnerConfig.Project.AccessToken),
	)

	t.Log("************** TestCase END: TestOuterloopApisRegistrationTest **************")
}
