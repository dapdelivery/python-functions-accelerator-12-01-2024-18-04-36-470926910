package multicluster_outerloop_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/envfuncs"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/models"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

var testenv env.Environment
var suiteConfig = models.SuiteConfig{}
var outerloopConfig = models.OuterloopConfig{}
var whereForDinnerConfig = models.WhereForDinnerConfig{}
var suiteResourcesDir = filepath.Join(utils.GetFileDir(), "../../resources/suite")

var airgapResourcesDir = filepath.Join(utils.GetFileDir(), "../../resources/airgap")
var serverWorkerConfig = models.ServerWorkerConfig{}
var appURL string
var airgapConfig = models.AirgapConfig{}
var ViewTapValuesFile string
var BuildTapValuesFile string
var RunTapValuesFile string
var appTestPipelineFile string

func TestMain(m *testing.M) {
	// set logger
	logFile, err := utils.SetLogger(filepath.Join(utils.GetFileDir(), "logs"))
	if err != nil {
		log.Fatal(fmt.Errorf("error while setting log file %s: %w", logFile, err))
	}
	home, _ := os.UserHomeDir()
	cfg, _ := envconf.NewFromFlags()
	cfg.WithKubeconfigFile(filepath.Join(home, ".kube", "config"))
	testenv = env.NewWithConfig(cfg)

	// read suite config
	suiteConfig = models.GetSuiteConfig()
	outerloopConfig, _ = models.GetOuterloopConfig()
	whereForDinnerConfig, _ = models.GetWhereForDinnerConfig()
	serverWorkerConfig, _ = models.GetServerWorkerConfig()
	nsp_namespaces := []string{"my-apps2"}
	//crearte appURL - check if we get it from some place
	currentSuiteConfig := models.GetSuiteConfig()
	appURL = "http://where-for-dinner" + "." + whereForDinnerConfig.Workload.Namespace + ".tap." + currentSuiteConfig.Multicluster.RunClusterName + "." + whereForDinnerConfig.Domain.Name
	airgapConfig = models.GetAirgapConfig()

	developerNamespaceFile := filepath.Join(suiteResourcesDir, "developer-namespace.yaml")
	secretFiles := []string{filepath.Join(airgapResourcesDir, airgapConfig.Airgap.ReposiliteCaCertificateSecretFile), filepath.Join(airgapResourcesDir, airgapConfig.Airgap.GitlabSecretFile), filepath.Join(airgapResourcesDir, airgapConfig.Airgap.TbsSettingsXMLSecretFile)}
	tbsPackageRepository := airgapConfig.Airgap.TBS.PackageRepository.Image + ":" + airgapConfig.Airgap.TBS.PackageRepository.Version

	if suiteConfig.AirgapEnv {
		ViewTapValuesFile = airgapConfig.Airgap.ViewTapValuesFile
		BuildTapValuesFile = airgapConfig.Airgap.BuildTapValuesFile
		RunTapValuesFile = airgapConfig.Airgap.RunTapValuesFile
	} else {
		ViewTapValuesFile = suiteConfig.Multicluster.ViewTapValuesFile
		BuildTapValuesFile = suiteConfig.Multicluster.BuildTapValuesFile
		RunTapValuesFile = suiteConfig.Multicluster.RunTapValuesFile
	}
	// set the airgap env varibale.
	os.Setenv("Airgap", strconv.FormatBool(suiteConfig.AirgapEnv))

	distribution := os.Getenv("DISTRIBUTION")
	if distribution == "Openshift" || distribution == "TKGM" || distribution == "TKGS" {
		appTestPipelineFile = whereForDinnerConfig.AppTestsPipeline.YamlFilePvtCloud
	} else {
		appTestPipelineFile = whereForDinnerConfig.AppTestsPipeline.YamlFile
	}

	// setup
	testenv.Setup(
		envfuncs.InstallTanzuCli(suiteConfig.TanzuCli.Version),
		envfuncs.UseContext(suiteConfig.Multicluster.ViewClusterContext),
		envfuncs.InstallClusterEssentials(suiteConfig.TanzuClusterEssentials.TanzunetHost, suiteConfig.TanzuClusterEssentials.TanzunetApiToken, suiteConfig.TanzuClusterEssentials.ProductFileId, suiteConfig.TanzuClusterEssentials.ReleaseVersion, suiteConfig.TanzuClusterEssentials.ProductSlug, suiteConfig.TanzuClusterEssentials.DownloadBundle, suiteConfig.TanzuClusterEssentials.InstallBundle, suiteConfig.TanzuClusterEssentials.InstallRegistryHostname, suiteConfig.TanzuClusterEssentials.InstallRegistryUsername, suiteConfig.TanzuClusterEssentials.InstallRegistryPassword),
		envfuncs.CreateNamespacesIfNotExists(suiteConfig.CreateNamespaces),
		envfuncs.CreateSecret(suiteConfig.TapRegistrySecret.Name, suiteConfig.TapRegistrySecret.Registry, suiteConfig.TapRegistrySecret.Username, suiteConfig.TapRegistrySecret.Password, suiteConfig.TapRegistrySecret.Namespace, suiteConfig.TapRegistrySecret.Export),
		envfuncs.CreateSecret(suiteConfig.RegistryCredentialsSecret.Name, suiteConfig.RegistryCredentialsSecret.Registry, suiteConfig.RegistryCredentialsSecret.Username, suiteConfig.RegistryCredentialsSecret.Password, suiteConfig.RegistryCredentialsSecret.Namespace, suiteConfig.RegistryCredentialsSecret.Export),
		envfuncs.AddPackageRepository(suiteConfig.PackageRepository.Name, suiteConfig.PackageRepository.Image, suiteConfig.PackageRepository.Namespace),
		envfuncs.CheckIfPackageRepositoryReconciled(suiteConfig.PackageRepository.Name, suiteConfig.PackageRepository.Namespace, 10, 60),
		envfuncs.InstallPackage(suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, suiteConfig.Tap.Namespace, ViewTapValuesFile, suiteConfig.Tap.PollTimeout),
		envfuncs.CheckIfPackageInstalled(suiteConfig.Tap.Name, suiteConfig.Tap.Namespace, 10, 60),
		envfuncs.ListInstalledPackages(suiteConfig.Tap.Namespace),
		envfuncs.UpdateDomainRecordsEnvFunc(suiteConfig.Multicluster.ViewClusterName, suiteConfig.Bind9DNSServer.UpdateDnsFilepath, suiteConfig.Bind9DNSServer.HostIP, suiteConfig.Bind9DNSServer.Password),

		envfuncs.UseContext(suiteConfig.Multicluster.BuildClusterContext),
		envfuncs.InstallClusterEssentials(suiteConfig.TanzuClusterEssentials.TanzunetHost, suiteConfig.TanzuClusterEssentials.TanzunetApiToken, suiteConfig.TanzuClusterEssentials.ProductFileId, suiteConfig.TanzuClusterEssentials.ReleaseVersion, suiteConfig.TanzuClusterEssentials.ProductSlug, suiteConfig.TanzuClusterEssentials.DownloadBundle, suiteConfig.TanzuClusterEssentials.InstallBundle, suiteConfig.TanzuClusterEssentials.InstallRegistryHostname, suiteConfig.TanzuClusterEssentials.InstallRegistryUsername, suiteConfig.TanzuClusterEssentials.InstallRegistryPassword),
		envfuncs.CreateNamespacesIfNotExists(suiteConfig.CreateNamespaces),
		envfuncs.CreateNamespacesIfNotExists(nsp_namespaces),
		envfuncs.CreateSecret(suiteConfig.TapRegistrySecret.Name, suiteConfig.TapRegistrySecret.Registry, suiteConfig.TapRegistrySecret.Username, suiteConfig.TapRegistrySecret.Password, suiteConfig.TapRegistrySecret.Namespace, suiteConfig.TapRegistrySecret.Export),
		envfuncs.CreateSecret(suiteConfig.RegistryCredentialsSecret.Name, suiteConfig.RegistryCredentialsSecret.Registry, suiteConfig.RegistryCredentialsSecret.Username, suiteConfig.RegistryCredentialsSecret.Password, suiteConfig.RegistryCredentialsSecret.Namespace, suiteConfig.RegistryCredentialsSecret.Export),
		envfuncs.CreateSecret(suiteConfig.RegistryCredentialsSecret.Name, suiteConfig.RegistryCredentialsSecret.Registry, suiteConfig.RegistryCredentialsSecret.Username, suiteConfig.RegistryCredentialsSecret.Password, suiteConfig.CreateNamespaces[1], suiteConfig.RegistryCredentialsSecret.Export),
		envfuncs.SetupAirgapEnvSecrets(suiteConfig.AirgapEnv, secretFiles, suiteConfig.CreateNamespaces[0]),
		envfuncs.AddPackageRepository(suiteConfig.PackageRepository.Name, suiteConfig.PackageRepository.Image, suiteConfig.PackageRepository.Namespace),
		envfuncs.CheckIfPackageRepositoryReconciled(suiteConfig.PackageRepository.Name, suiteConfig.PackageRepository.Namespace, 10, 60),
		envfuncs.InstallPackage(suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, suiteConfig.Tap.Namespace, BuildTapValuesFile, suiteConfig.Tap.PollTimeout),
		envfuncs.AddTBSFullDepsRepo(suiteConfig.AirgapEnv, airgapConfig.Airgap.TBS.PackageRepository.Name, tbsPackageRepository, airgapConfig.Airgap.TBS.PackageRepository.Namespace),
		envfuncs.CheckIfTBSPackageRepositoryReconciled(suiteConfig.AirgapEnv, airgapConfig.Airgap.TBS.PackageRepository.Name, airgapConfig.Airgap.TBS.PackageRepository.Namespace, 10, 60),
		envfuncs.InstallTBSFullDepsPkg(suiteConfig.AirgapEnv, airgapConfig.Airgap.TBS.Package.Name, airgapConfig.Airgap.TBS.Package.PackageName, airgapConfig.Airgap.TBS.Package.Version, airgapConfig.Airgap.TBS.Package.Namespace, "", ""),
		envfuncs.CheckIfPackageInstalled(suiteConfig.Tap.Name, suiteConfig.Tap.Namespace, 10, 60),
		envfuncs.CheckIfTBSFullDepsPackageInstalled(suiteConfig.AirgapEnv, airgapConfig.Airgap.TBS.Package.Name, airgapConfig.Airgap.TBS.Package.Namespace, 10, 60),

		envfuncs.CheckForSCCs(),
		envfuncs.ListInstalledPackages(suiteConfig.Tap.Namespace),
		envfuncs.SetupDeveloperNamespace(developerNamespaceFile, suiteConfig.CreateNamespaces[0]),
		envfuncs.UpdateDomainRecordsEnvFunc(suiteConfig.Multicluster.BuildClusterName, suiteConfig.Bind9DNSServer.UpdateDnsFilepath, suiteConfig.Bind9DNSServer.HostIP, suiteConfig.Bind9DNSServer.Password),

		envfuncs.UseContext(suiteConfig.Multicluster.RunClusterContext),
		envfuncs.InstallClusterEssentials(suiteConfig.TanzuClusterEssentials.TanzunetHost, suiteConfig.TanzuClusterEssentials.TanzunetApiToken, suiteConfig.TanzuClusterEssentials.ProductFileId, suiteConfig.TanzuClusterEssentials.ReleaseVersion, suiteConfig.TanzuClusterEssentials.ProductSlug, suiteConfig.TanzuClusterEssentials.DownloadBundle, suiteConfig.TanzuClusterEssentials.InstallBundle, suiteConfig.TanzuClusterEssentials.InstallRegistryHostname, suiteConfig.TanzuClusterEssentials.InstallRegistryUsername, suiteConfig.TanzuClusterEssentials.InstallRegistryPassword),
		envfuncs.CreateNamespacesIfNotExists(suiteConfig.CreateNamespaces),
		envfuncs.CreateNamespacesIfNotExists(nsp_namespaces),
		envfuncs.CreateSecret(suiteConfig.TapRegistrySecret.Name, suiteConfig.TapRegistrySecret.Registry, suiteConfig.TapRegistrySecret.Username, suiteConfig.TapRegistrySecret.Password, suiteConfig.TapRegistrySecret.Namespace, suiteConfig.TapRegistrySecret.Export),
		envfuncs.CreateSecret(suiteConfig.RegistryCredentialsSecret.Name, suiteConfig.RegistryCredentialsSecret.Registry, suiteConfig.RegistryCredentialsSecret.Username, suiteConfig.RegistryCredentialsSecret.Password, suiteConfig.RegistryCredentialsSecret.Namespace, suiteConfig.RegistryCredentialsSecret.Export),
		envfuncs.AddPackageRepository(suiteConfig.PackageRepository.Name, suiteConfig.PackageRepository.Image, suiteConfig.PackageRepository.Namespace),
		envfuncs.CheckIfPackageRepositoryReconciled(suiteConfig.PackageRepository.Name, suiteConfig.PackageRepository.Namespace, 10, 60),
		envfuncs.InstallPackage(suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, suiteConfig.Tap.Namespace, RunTapValuesFile, suiteConfig.Tap.PollTimeout),
		envfuncs.CheckIfPackageInstalled(suiteConfig.Tap.Name, suiteConfig.Tap.Namespace, 10, 60),
		envfuncs.CheckForSCCs(),
		envfuncs.ListInstalledPackages(suiteConfig.Tap.Namespace),
		envfuncs.SetupDeveloperNamespace(developerNamespaceFile, suiteConfig.CreateNamespaces[0]),
		envfuncs.UpdateDomainRecordsEnvFunc(suiteConfig.Multicluster.RunClusterName, suiteConfig.Bind9DNSServer.UpdateDnsFilepath, suiteConfig.Bind9DNSServer.HostIP, suiteConfig.Bind9DNSServer.Password),
		envfuncs.PreReqForWhereForDinnerApp(&suiteConfig, &whereForDinnerConfig),
	)

	// finish
	testenv.Finish(
		envfuncs.UseContext(suiteConfig.Multicluster.ViewClusterContext),
		envfuncs.UninstallPackage(suiteConfig.Tap.Name, suiteConfig.Tap.Namespace),
		envfuncs.DeletePackageRepository(suiteConfig.PackageRepository.Name, suiteConfig.PackageRepository.Namespace),
		envfuncs.DeleteSecret(suiteConfig.RegistryCredentialsSecret.Name, suiteConfig.RegistryCredentialsSecret.Namespace),
		envfuncs.DeleteSecret(suiteConfig.TapRegistrySecret.Name, suiteConfig.TapRegistrySecret.Namespace),
		envfuncs.DeleteNamespaces(suiteConfig.CreateNamespaces),

		envfuncs.UseContext(suiteConfig.Multicluster.BuildClusterContext),
		envfuncs.DeleteDeveloperNamespace(developerNamespaceFile, suiteConfig.CreateNamespaces[0]),
		envfuncs.UninstallPackage(suiteConfig.Tap.Name, suiteConfig.Tap.Namespace),
		envfuncs.DeletePackageRepository(suiteConfig.PackageRepository.Name, suiteConfig.PackageRepository.Namespace),
		envfuncs.DeleteSecret(suiteConfig.RegistryCredentialsSecret.Name, suiteConfig.RegistryCredentialsSecret.Namespace),
		envfuncs.DeleteSecret(suiteConfig.TapRegistrySecret.Name, suiteConfig.TapRegistrySecret.Namespace),
		envfuncs.DeleteSecret(suiteConfig.TapRegistrySecret.Name, suiteConfig.CreateNamespaces[1]),
		envfuncs.DeleteAirgapEnvSecrets(suiteConfig.AirgapEnv, secretFiles, suiteConfig.CreateNamespaces[0]),
		envfuncs.UninstallTBSFullDepsPkg(suiteConfig.AirgapEnv, airgapConfig.Airgap.TBS.Package.Name, airgapConfig.Airgap.TBS.Package.Namespace),
		envfuncs.DeleteTBSFullDepsRepo(suiteConfig.AirgapEnv, airgapConfig.Airgap.TBS.PackageRepository.Name, airgapConfig.Airgap.TBS.PackageRepository.Namespace),
		envfuncs.DeleteNamespaces(suiteConfig.CreateNamespaces),
		envfuncs.DeleteNamespaces(nsp_namespaces),

		envfuncs.UseContext(suiteConfig.Multicluster.RunClusterContext),
		envfuncs.WhereForDinneApprCleanUp(&suiteConfig, &whereForDinnerConfig),
		envfuncs.DeleteDeveloperNamespace(developerNamespaceFile, suiteConfig.CreateNamespaces[0]),
		envfuncs.UninstallPackage(suiteConfig.Tap.Name, suiteConfig.Tap.Namespace),
		envfuncs.DeletePackageRepository(suiteConfig.PackageRepository.Name, suiteConfig.PackageRepository.Namespace),
		envfuncs.DeleteSecret(suiteConfig.RegistryCredentialsSecret.Name, suiteConfig.RegistryCredentialsSecret.Namespace),
		envfuncs.DeleteSecret(suiteConfig.TapRegistrySecret.Name, suiteConfig.TapRegistrySecret.Namespace),
		envfuncs.DeleteNamespaces(suiteConfig.CreateNamespaces),
		envfuncs.DeleteNamespaces(nsp_namespaces),
	)

	os.Exit(testenv.Run(m))
}
