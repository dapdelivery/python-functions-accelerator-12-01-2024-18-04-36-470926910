//go:build all || multicluster_outerloop || multicluster_dockerfile_gitsource
// +build all multicluster_outerloop multicluster_dockerfile_gitsource

package multicluster_outerloop_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_libs"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/models"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestDockerFileBasedWorkloadBasicSupplychainGitSource(t *testing.T) {
	t.Log("************** TestCase START: TestDockerFileBasedWorkloadBasicSupplychainGitSource **************")
	currentSuiteConfig := models.GetSuiteConfig()
	wURL := "http://" + suiteConfig.Innerloop.Workload.Name + "." + outerloopConfig.Namespace + ".tap." + currentSuiteConfig.Multicluster.RunClusterName + "." + suiteConfig.Domain

	openshiftDeveloperNamespaceFile := features.New("update-openshift-developer-namespace-for-dockerfilebased").
		Assess("update-namespace-for-dockerfilebase", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("updating developer namespace for openshift dockerfilebased test")
			distribution := os.Getenv("DISTRIBUTION")
			// adding extra RoleBinding for openshift for build (https://jira.eng.vmware.com/browse/TANZUSC-2004)
			if distribution == "Openshift" {
				dockerFileBasedDeveloperNamespaceFile := filepath.Join(suiteResourcesDir, "developer-namespace-mc-openshift.yaml")
				err := kubectl_libs.KubectlApplyConfiguration(dockerFileBasedDeveloperNamespaceFile, suiteConfig.CreateNamespaces[0])
				if err != nil {
					t.Errorf("error while deploying %s", dockerFileBasedDeveloperNamespaceFile)
					t.FailNow()
				} else {
					t.Logf("deployed %s", dockerFileBasedDeveloperNamespaceFile)
				}
			}
			return ctx
		}).
		Feature()

	testenv.Test(t,
		//build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		openshiftDeveloperNamespaceFile,
		common_features.TanzuCreateWorkloadWithDockerFile(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Gitrepository, suiteConfig.Innerloop.Workload.Branch, outerloopConfig.Namespace),
		common_features.VerifyGitRepoStatus(t, suiteConfig.Innerloop.Workload.Name, outerloopConfig.Namespace),
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
	)
	t.Log("************** TestCase END: TestDockerFileBasedWorkloadBasicSupplychainGitSource **************")
}
