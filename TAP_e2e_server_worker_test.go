//go:build all || multicluster_outerloop || server_worker_basic_supply_chain
// +build all multicluster_outerloop server_worker_basic_supply_chain

package multicluster_outerloop_test

import (
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
)

func TestServerWorkerOuterloopBasicSupplyChain(t *testing.T) {
	t.Log("************** TestCase START: TestServerWorkerOuterloopBasicSupplyChain **************")
	testenv.Test(t,
		//view context
		common_features.ChangeContext(t, suiteConfig.Multicluster.ViewClusterContext),
		common_features.UpdateTapValuesAccelerator(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, "view", "basic", suiteConfig.Tap.Namespace, "LoadBalancer"),
		common_features.GenerateSMTPAcceleratorProject(t, suiteConfig.Tap.Namespace, serverWorkerConfig.Workload.AcceleratorName, serverWorkerConfig.Workload.AcceleratorName, serverWorkerConfig.Workload.RabbitmqClusterNodes),

		// build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.UpdateTapProfileSupplyChain(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, "build", "basic", suiteConfig.Tap.Namespace, suiteConfig.Tap.PollTimeout),
		common_features.ReplaceStringInFileWithoutCompile(t, serverWorkerConfig.Workload.DefaultNamespace, outerloopConfig.Namespace, serverWorkerConfig.Workload.Path, serverWorkerConfig.Workload.AcceleratorName),
		common_features.ReplaceStringInFileWithoutCompile(t, serverWorkerConfig.Workload.DefaultUrl, serverWorkerConfig.Workload.Url, serverWorkerConfig.Workload.Path, serverWorkerConfig.Workload.AcceleratorName),

		common_features.DeployWorkloadFromAcceleratorProject(t, serverWorkerConfig.Workload.AcceleratorName, serverWorkerConfig.Workload.Path, outerloopConfig.Namespace),

		// server workload verification
		common_features.VerifyBuildStatus(t, serverWorkerConfig.Workload.Server, suiteConfig.Innerloop.Workload.BuildNameSuffix, outerloopConfig.Namespace),
		common_features.VerifyTanzuWorkloadStatus(t, serverWorkerConfig.Workload.Server, outerloopConfig.Namespace),

		// worker workload verification
		common_features.VerifyBuildStatus(t, serverWorkerConfig.Workload.Worker, suiteConfig.Innerloop.Workload.BuildNameSuffix, outerloopConfig.Namespace),
		common_features.VerifyTanzuWorkloadStatus(t, serverWorkerConfig.Workload.Worker, outerloopConfig.Namespace),

		// run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.UpdateTapProfileSupplyChain(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, "run", "basic", suiteConfig.Tap.Namespace, suiteConfig.Tap.PollTimeout),
		common_features.ApplyKubectlConfigurationFile(t, suiteConfig.ServiceToolkit.Gitrepository, ""),
		common_features.CreateNamespacesIfNotExist(t, serverWorkerConfig.Workload.ServiceInstanceNamespace),
		common_features.ApplyKubectlConfigurationFile(t, serverWorkerConfig.Workload.RmqInstanceFile, serverWorkerConfig.Workload.ServiceInstanceNamespace),
		common_features.ApplyKubectlConfigurationFile(t, serverWorkerConfig.Workload.RmqResourceClaimFile, ""),

		// copying deliverable from build to run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.ProcessDeliverable(t, serverWorkerConfig.Workload.Server, outerloopConfig.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext, ""),
		common_features.ProcessDeliverable(t, serverWorkerConfig.Workload.Worker, outerloopConfig.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext, ""),

		// run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.ValidateDeploymentExists(t, serverWorkerConfig.Workload.Server, outerloopConfig.Namespace),
		common_features.PortForward(t, serverWorkerConfig.Workload.Server, serverWorkerConfig.Workload.PortNumber, outerloopConfig.Namespace),

		common_features.ExecuteTelnetMessage(t),
		common_features.WorkloadResponseFromLogs(t, serverWorkerConfig.Workload.Worker, outerloopConfig.Namespace, serverWorkerConfig.Workload.VerificationString),
		common_features.KillPortForwardingProcess(t),

		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.MulticlusterOuterloopCleanup(t, serverWorkerConfig.Workload.Server, serverWorkerConfig.Workload.AcceleratorName, outerloopConfig.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext),
		common_features.MulticlusterOuterloopCleanup(t, serverWorkerConfig.Workload.Worker, serverWorkerConfig.Workload.AcceleratorName, outerloopConfig.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext),
		common_features.DeleteNamespaceIfExists(t, serverWorkerConfig.Workload.ServiceInstanceNamespace, suiteConfig.Multicluster.RunClusterContext),
		common_features.Removedir(t, serverWorkerConfig.Workload.AcceleratorName),
	)

	t.Log("************** TestCase END: TestServerWorkerOuterloopBasicSupplyChain **************")
}
