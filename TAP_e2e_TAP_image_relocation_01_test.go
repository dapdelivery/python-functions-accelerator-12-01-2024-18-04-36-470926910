//go:build all || multicluster_outerloop || multicluster_image_relocation
// +build all multicluster_outerloop multicluster_image_relocation

package multicluster_outerloop_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
)

func TestImageRelocation(t *testing.T) {
	t.Log("************** TestCase START: TestImageRelocation **************")
	rand.Seed(time.Now().UnixNano())
	randUniq := strconv.Itoa(rand.Intn(1000000000000))

	secretName := fmt.Sprintf("%s-2", suiteConfig.TapRegistrySecret.Name)
	testenv.Test(t,
		common_features.CreateGithubRepo(t, whereForDinnerConfig.Project.Name, whereForDinnerConfig.Project.RepoTemplate, whereForDinnerConfig.Project.AccessToken),

		//build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.TanzuDeployListOfWorkloads(t, whereForDinnerConfig.Workload.YamlFilePath, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfTanzuWorkloadStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace, whereForDinnerConfig.Workload.YamlFilePath),
		common_features.VerifyListOfBuildStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.BuildNameSuffix, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyPodIntentStatus(t, whereForDinnerConfig.Workload.Name[1], whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyListOfTaskRunStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.TaskRunInfix, whereForDinnerConfig.Workload.Namespace),

		//imgpkg copy to target repo
		common_features.DockerLogin(t, suiteConfig.AzureRepository.Server, suiteConfig.AzureRepository.Username, suiteConfig.AzureRepository.Password),
		common_features.ImageCopyFromDeliverableListToRepo(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace, whereForDinnerConfig.TestTargetRepo, randUniq),

		//run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.CreateSecret(t, secretName, suiteConfig.AzureRepository.Server, suiteConfig.AzureRepository.Username, suiteConfig.AzureRepository.Password, "", whereForDinnerConfig.Workload.Namespace, true),

		//copying deliverable from build to run context
		common_features.ProcessListOfDeliverables(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext, whereForDinnerConfig.TestTargetRepo, randUniq),

		//run context
		common_features.VerifyListOfKsvcStatus(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyAppResponseUsingCurl(t, appURL, whereForDinnerConfig.Workload.EndPoint, whereForDinnerConfig.Workload.OriginalVerificationString),

		//build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.UpdateGitRepository(t, whereForDinnerConfig.Project.Username, whereForDinnerConfig.Project.Email, whereForDinnerConfig.Project.Repository, whereForDinnerConfig.Project.Name, whereForDinnerConfig.Project.AccessToken, whereForDinnerConfig.Workload.ApplicationFilePath, whereForDinnerConfig.Workload.OriginalString, whereForDinnerConfig.Workload.NewString, whereForDinnerConfig.Project.CommitMessage),
		common_features.VerifyBuildStatusAfterUpdate(t, whereForDinnerConfig.Workload.Name[1], whereForDinnerConfig.Workload.Namespace),
		common_features.ImageCopyFromDeliverableListToRepo(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace, whereForDinnerConfig.TestTargetRepo, randUniq),

		//run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.VerifyKsvcStatusAfterUpdate(t, whereForDinnerConfig.Workload.Name[1], whereForDinnerConfig.Workload.Namespace),
		common_features.VerifyAppResponseUsingCurl(t, appURL, whereForDinnerConfig.Workload.EndPoint, whereForDinnerConfig.Workload.NewVerificationString),

		common_features.DeleteGithubRepo(t, whereForDinnerConfig.Project.Name, whereForDinnerConfig.Project.AccessToken),
		common_features.DeliverableCleanup(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
		//build cluster cleanup
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.DeleteWorkloadAll(t, whereForDinnerConfig.Workload.Name, whereForDinnerConfig.Workload.Namespace),
	)

	t.Log("************** TestCase END: TestImageRelocation **************")
}
