package rootcommand

import (
	"fmt"
	"os"

	"github.com/pivotal/kpack/pkg/logs"
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/kpack-cli/pkg/kpackcompat"

	"github.com/vmware-tanzu/kpack-cli/pkg/commands"
	buildcmds "github.com/vmware-tanzu/kpack-cli/pkg/commands/build"
	buildercmds "github.com/vmware-tanzu/kpack-cli/pkg/commands/builder"
	clusterbuildercmds "github.com/vmware-tanzu/kpack-cli/pkg/commands/clusterbuilder"
	clusterstackcmds "github.com/vmware-tanzu/kpack-cli/pkg/commands/clusterstack"
	clusterstorecmds "github.com/vmware-tanzu/kpack-cli/pkg/commands/clusterstore"
	configcmds "github.com/vmware-tanzu/kpack-cli/pkg/commands/config"
	imgcmds "github.com/vmware-tanzu/kpack-cli/pkg/commands/image"
	importcmds "github.com/vmware-tanzu/kpack-cli/pkg/commands/import"
	"github.com/vmware-tanzu/kpack-cli/pkg/commands/lifecycle"
	secretcmds "github.com/vmware-tanzu/kpack-cli/pkg/commands/secret"
	importpkg "github.com/vmware-tanzu/kpack-cli/pkg/import"
	"github.com/vmware-tanzu/kpack-cli/pkg/k8s"
	"github.com/vmware-tanzu/kpack-cli/pkg/registry"
	"github.com/vmware-tanzu/kpack-cli/pkg/secret"
)

var (
	Version   = "dev"
	CommitSHA = ""
)

func GetRootCommand() *cobra.Command {
	var clientSetProvider k8s.DefaultClientSetProvider

	rootCmd := &cobra.Command{
		Use: "kp",
		Long: `kp controls the kpack installation on Kubernetes.

kpack extends Kubernetes and utilizes unprivileged kubernetes primitives to provide 
builds of OCI images as a platform implementation of Cloud Native Buildpacks (CNB).
Learn more about kpack @ https://github.com/pivotal/kpack`,
	}
	rootCmd.AddCommand(
		getVersionCommand(),
		getImageCommand(clientSetProvider),
		getBuildCommand(clientSetProvider),
		getSecretCommand(clientSetProvider),
		getClusterBuilderCommand(clientSetProvider),
		getBuilderCommand(clientSetProvider),
		getStackCommand(clientSetProvider),
		getStoreCommand(clientSetProvider),
		getLifecycleCommand(clientSetProvider),
		getImportCommand(clientSetProvider),
		getConfigCommand(clientSetProvider),
		getCompletionCommand(),
	)

	return rootCmd
}

func getVersionCommand() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Display kp version",
		Run: func(cmd *cobra.Command, _ []string) {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), Version+" "+CommitSHA)
		},
	}
	return versionCmd
}

func getImageCommand(clientSetProvider k8s.ClientSetProvider) *cobra.Command {
	newImageWaiter := func(clientSet k8s.ClientSet) imgcmds.ImageWaiter {
		return kpackcompat.NewImageWaiterForV1alpha2(logs.NewImageWaiter(clientSet.KpackClient, logs.NewBuildLogsClient(clientSet.K8sClient)))
	}

	imageRootCmd := &cobra.Command{
		Use:     "image",
		Short:   "Image commands",
		Aliases: []string{"images", "imgs", "img"},
	}
	imageRootCmd.AddCommand(
		imgcmds.NewCreateCommand(clientSetProvider, registry.DefaultUtilProvider{}, newImageWaiter),
		imgcmds.NewPatchCommand(clientSetProvider, registry.DefaultUtilProvider{}, newImageWaiter),
		imgcmds.NewSaveCommand(clientSetProvider, registry.DefaultUtilProvider{}, newImageWaiter),
		imgcmds.NewListCommand(clientSetProvider),
		imgcmds.NewDeleteCommand(clientSetProvider),
		imgcmds.NewTriggerCommand(clientSetProvider),
		imgcmds.NewStatusCommand(clientSetProvider),
	)
	return imageRootCmd
}

func getBuildCommand(clientSetProvider k8s.ClientSetProvider) *cobra.Command {
	buildRootCmd := &cobra.Command{
		Use:     "build",
		Short:   "Build Commands",
		Aliases: []string{"builds", "blds", "bld"},
	}
	buildRootCmd.AddCommand(
		buildcmds.NewListCommand(clientSetProvider),
		buildcmds.NewStatusCommand(clientSetProvider),
		buildcmds.NewLogsCommand(clientSetProvider),
	)
	return buildRootCmd
}

func getSecretCommand(clientSetProvider k8s.ClientSetProvider) *cobra.Command {
	credentialFetcher := &commands.CredentialFetcher{}
	secretFactory := &secret.Factory{
		CredentialFetcher: credentialFetcher,
	}

	secretRootCmd := &cobra.Command{
		Use:     "secret",
		Short:   "Secret Commands",
		Aliases: []string{"secrets"},
	}
	secretRootCmd.AddCommand(
		secretcmds.NewCreateCommand(clientSetProvider, secretFactory),
		secretcmds.NewDeleteCommand(clientSetProvider),
		secretcmds.NewListCommand(clientSetProvider),
	)
	return secretRootCmd
}

func getClusterBuilderCommand(clientSetProvider k8s.ClientSetProvider) *cobra.Command {
	clusterBuilderRootCmd := &cobra.Command{
		Use:     "clusterbuilder",
		Short:   "ClusterBuilder Commands",
		Aliases: []string{"clusterbuilders", "clstrbldrs", "clstrbldr", "cbldrs", "cbldr", "cbs", "cb"},
	}
	clusterBuilderRootCmd.AddCommand(
		clusterbuildercmds.NewCreateCommand(clientSetProvider, commands.NewResourceWaiter),
		clusterbuildercmds.NewPatchCommand(clientSetProvider, commands.NewResourceWaiter),
		clusterbuildercmds.NewSaveCommand(clientSetProvider, commands.NewResourceWaiter),
		clusterbuildercmds.NewListCommand(clientSetProvider),
		clusterbuildercmds.NewStatusCommand(clientSetProvider),
		clusterbuildercmds.NewDeleteCommand(clientSetProvider),
	)
	return clusterBuilderRootCmd
}

func getBuilderCommand(clientSetProvider k8s.ClientSetProvider) *cobra.Command {
	builderRootCmd := &cobra.Command{
		Use:     "builder",
		Short:   "Builder Commands",
		Aliases: []string{"builders", "bldrs", "bldr"},
	}
	builderRootCmd.AddCommand(
		buildercmds.NewCreateCommand(clientSetProvider, commands.NewResourceWaiter),
		buildercmds.NewPatchCommand(clientSetProvider, commands.NewResourceWaiter),
		buildercmds.NewSaveCommand(clientSetProvider, commands.NewResourceWaiter),
		buildercmds.NewListCommand(clientSetProvider),
		buildercmds.NewDeleteCommand(clientSetProvider),
		buildercmds.NewStatusCommand(clientSetProvider),
	)
	return builderRootCmd
}

func getStackCommand(clientSetProvider k8s.ClientSetProvider) *cobra.Command {
	stackRootCmd := &cobra.Command{
		Use:     "clusterstack",
		Aliases: []string{"clusterstacks", "clstrcsks", "clstrcsk", "cstacks", "cstack", "cstks", "cstk", "csks", "csk"},
		Short:   "ClusterStack Commands",
	}
	stackRootCmd.AddCommand(
		clusterstackcmds.NewCreateCommand(clientSetProvider, registry.DefaultUtilProvider{}, commands.NewResourceWaiter),
		clusterstackcmds.NewPatchCommand(clientSetProvider, registry.DefaultUtilProvider{}, commands.NewResourceWaiter),
		clusterstackcmds.NewSaveCommand(clientSetProvider, registry.DefaultUtilProvider{}, commands.NewResourceWaiter),
		clusterstackcmds.NewListCommand(clientSetProvider),
		clusterstackcmds.NewStatusCommand(clientSetProvider),
		clusterstackcmds.NewDeleteCommand(clientSetProvider),
	)
	return stackRootCmd
}

func getStoreCommand(clientSetProvider k8s.ClientSetProvider) *cobra.Command {
	storeRootCommand := &cobra.Command{
		Use:     "clusterstore",
		Aliases: []string{"clusterstores", "clstrcsrs", "clstrcsr", "cstores", "cstore", "cstrs", "cstr", "csrs", "csr"},
		Short:   "ClusterStore Commands",
	}
	storeRootCommand.AddCommand(
		clusterstorecmds.NewCreateCommand(clientSetProvider, registry.DefaultUtilProvider{}, commands.NewResourceWaiter),
		clusterstorecmds.NewAddCommand(clientSetProvider, registry.DefaultUtilProvider{}, commands.NewResourceWaiter),
		clusterstorecmds.NewSaveCommand(clientSetProvider, registry.DefaultUtilProvider{}, commands.NewResourceWaiter),
		clusterstorecmds.NewDeleteCommand(clientSetProvider, commands.NewConfirmationProvider()),
		clusterstorecmds.NewStatusCommand(clientSetProvider),
		clusterstorecmds.NewRemoveCommand(clientSetProvider, commands.NewResourceWaiter),
		clusterstorecmds.NewListCommand(clientSetProvider),
	)

	return storeRootCommand
}

func getLifecycleCommand(clientSetProvider k8s.ClientSetProvider) *cobra.Command {
	lifecycleRootCommand := &cobra.Command{
		Use:   "lifecycle",
		Short: "Lifecycle Commands",
	}
	lifecycleRootCommand.AddCommand(
		lifecycle.NewUpdateCommand(clientSetProvider, registry.DefaultUtilProvider{}),
	)
	return lifecycleRootCommand
}

func getImportCommand(clientSetProvider k8s.ClientSetProvider) *cobra.Command {
	return importcmds.NewImportCommand(
		commands.Differ{},
		clientSetProvider,
		registry.DefaultUtilProvider{},
		importpkg.DefaultTimestampProvider(),
		commands.NewConfirmationProvider(),
		commands.NewResourceWaiter,
	)
}

func getConfigCommand(clientSetProvider k8s.ClientSetProvider) *cobra.Command {
	configRootCmd := &cobra.Command{
		Use:     "config",
		Short:   "Config commands",
		Aliases: []string{"config", "cfg"},
	}
	configRootCmd.AddCommand(
		configcmds.NewDefaultRepositoryCommand(clientSetProvider),
		configcmds.NewDefaultServiceAccountCommand(clientSetProvider),
	)

	return configRootCmd
}

func getCompletionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:

$ source <(kp completion bash)

# To load completions for each session, execute once:
Linux:
  $ kp completion bash > /etc/bash_completion.d/kp
MacOS:
  $ kp completion bash > /usr/local/etc/bash_completion.d/kp

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ kp completion zsh > "${fpath[1]}/_kp"

# You will need to start a new shell for this setup to take effect.

Fish:

$ kp completion fish | source

# To load completions for each session, execute once:
$ kp completion fish > ~/.config/fish/completions/kp.fish
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletion(os.Stdout)
			}
		},
	}
}
