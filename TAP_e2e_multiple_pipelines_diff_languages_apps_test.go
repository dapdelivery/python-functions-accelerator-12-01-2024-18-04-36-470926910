//go:build all || multicluster_outerloop || multicluster_multiple_pipelines_diff_language_apps
// +build all multicluster_outerloop multicluster_multiple_pipelines_diff_language_apps

package multicluster_outerloop_test

import (
	"os"
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/models"
)

func TestMultiplePipelinesDiffLanguageApp(t *testing.T) {
	distribution := os.Getenv("DISTRIBUTION")
	if distribution == "Openshift" {
		t.Log("TestCase END: --- BLOCKED: TestMultiplePipelinesDiffLanguageApp (TANZUSC-5200) ")
		return
	}
	t.Log("************** TestCase START: TestMultiplePipelinesDiffLanguageApp **************")
	currentSuiteConfig := models.GetSuiteConfig()
	yamlParamJava := "testing_pipeline_matching_labels=$'apps.tanzu.vmware.com/pipeline: test\\napps.tanzu.vmware.com/language: java'"
	yamlParamGo := "testing_pipeline_matching_labels=$'apps.tanzu.vmware.com/pipeline: test\\napps.tanzu.vmware.com/language: golang'"
	wURL := "http://" + suiteConfig.Innerloop.Workload.Name + "." + outerloopConfig.Namespace + ".tap." + currentSuiteConfig.Multicluster.RunClusterName + "." + suiteConfig.Domain
	goURL := "http://" + outerloopConfig.GolangWorkload.Name + "." + outerloopConfig.Namespace + ".tap." + currentSuiteConfig.Multicluster.RunClusterName + "." + suiteConfig.Domain

	testenv.Test(t,

		//build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.UpdateTapProfileSupplyChain(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, "build", "testing", suiteConfig.Tap.Namespace, suiteConfig.Tap.PollTimeout),
		common_features.CreateAndApplyTestPipeline(t, whereForDinnerConfig.AppTestsPipeline.YamlFile, outerloopConfig.WorkloadAppSSO.PipelineName, "java", whereForDinnerConfig.Workload.Namespace),
		common_features.CreateAndApplyTestPipeline(t, whereForDinnerConfig.AppTestsPipeline.YamlFile, outerloopConfig.GolangWorkload.Gopipelinename, "golang", whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyPipelineStatus(t, outerloopConfig.WorkloadAppSSO.PipelineName, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyPipelineStatus(t, outerloopConfig.GolangWorkload.Gopipelinename, whereForDinnerConfig.Workload.Namespace),
		common_features.TanzuCreateWorkloadUsingParam(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Gitrepository, whereForDinnerConfig.Workload.Namespace, suiteConfig.Innerloop.Workload.Branch, "true", "web", yamlParamJava),
		common_features.TanzuCreateWorkloadUsingParam(t, outerloopConfig.GolangWorkload.Name, outerloopConfig.GolangWorkload.Gitrepository, whereForDinnerConfig.Workload.Namespace, outerloopConfig.GolangWorkload.Branch, "true", "web", yamlParamGo),

		common_features.VerifyListOfTanzuWorkloadStatus(t, []string{suiteConfig.Innerloop.Workload.Name, outerloopConfig.GolangWorkload.Name}, whereForDinnerConfig.Workload.Namespace, ""),
		common_features.VerifyListOfGitRepoStatus(t, []string{suiteConfig.Innerloop.Workload.Name, outerloopConfig.GolangWorkload.Name}, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListPipelineRunStatus(t, []string{suiteConfig.Innerloop.Workload.Name, outerloopConfig.GolangWorkload.Name}, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyPipelineLabels(t, outerloopConfig.GolangWorkload.Name, whereForDinnerConfig.Workload.Namespace, outerloopConfig.GolangWorkload.Gopipelinename),
		common_features.VerifyPipelineLabels(t, suiteConfig.Innerloop.Workload.Name, whereForDinnerConfig.Workload.Namespace, outerloopConfig.WorkloadAppSSO.PipelineName),
		common_features.VerifyListOfBuildStatus(t, []string{suiteConfig.Innerloop.Workload.Name, outerloopConfig.GolangWorkload.Name}, outerloopConfig.Workload.BuildNameSuffix, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfTestTaskRunStatus(t, []string{suiteConfig.Innerloop.Workload.Name, outerloopConfig.GolangWorkload.Name}, whereForDinnerConfig.Workload.TaskRunTestSuffix, whereForDinnerConfig.Workload.Namespace),

		// // run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),

		// // copying deliverable from build to run context
		common_features.ProcessListOfDeliverables(t, []string{suiteConfig.Innerloop.Workload.Name, outerloopConfig.GolangWorkload.Name}, whereForDinnerConfig.Workload.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext, "", ""),

		// //run context
		common_features.VerifyListOfRevisionStatus(t, []string{suiteConfig.Innerloop.Workload.Name, outerloopConfig.GolangWorkload.Name}, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfKsvcStatus(t, []string{suiteConfig.Innerloop.Workload.Name, outerloopConfig.GolangWorkload.Name}, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyAppResponseUsingCurl(t, goURL, outerloopConfig.GolangWorkload.Endpoint, outerloopConfig.GolangWorkload.OriginalString),
		common_features.VerifyAppResponseUsingCurl(t, wURL, "", suiteConfig.Innerloop.Workload.TanzuAppsString),

		// //run cluster cleaup
		common_features.DeliverableCleanup(t, []string{outerloopConfig.GolangWorkload.Name, suiteConfig.Innerloop.Workload.Name}, whereForDinnerConfig.Workload.Namespace),

		// //build cluster cleanup
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.DeleteWorkloadAll(t, []string{suiteConfig.Innerloop.Workload.Name, outerloopConfig.GolangWorkload.Name}, whereForDinnerConfig.Workload.Namespace),
		common_features.DeletePipeline(t, outerloopConfig.WorkloadAppSSO.PipelineName, whereForDinnerConfig.Workload.Namespace),
		common_features.DeletePipeline(t, outerloopConfig.GolangWorkload.Gopipelinename, whereForDinnerConfig.Workload.Namespace),
		common_features.UpdateTapProfileSupplyChain(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, "build", "basic", suiteConfig.Tap.Namespace, suiteConfig.Tap.PollTimeout),
	)
	t.Log("************** TestCase END: TestMultiplePipelinesDiffLanguageApp **************")
}
