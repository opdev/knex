package run

import (
	"context"

	"github.com/opdev/knex/plugin/v0"
	"github.com/opdev/knex/types"
	"github.com/redhat-openshift-ecosystem/openshift-preflight/artifacts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	spfviper "github.com/spf13/viper"
)

// NewBackwardsCompatCheckCommand is a compatibilty command bridging Preflight's
// Pluggable design to Preflight's legacy design. It's the equivalent of the
// `preflight check` command, but will run the corresponding plugins instead. It
// is expected that this subcommand will be removed near-future after this
// redesign is published.
func NewBackwardsCompatCheckCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Run checks for an operator or container. This is subcommand exists for backwards compatibility and will be removed in a future release.",
		Long:  "This command will allow you to execute the Red Hat Certification tests for an operator or a container.",
	}

	cmd.PersistentFlags().String("logfile", "", "Where the execution logfile will be written. (env: PFLT_LOGFILE)")
	cmd.PersistentFlags().String("loglevel", "", "The verbosity of the preflight tool itself. Ex. warn, debug, trace, info, error. (env: PFLT_LOGLEVEL)")
	cmd.PersistentFlags().String("artifacts", "", "Where check-specific artifacts will be written. (env: PFLT_ARTIFACTS)")

	containerConfig := spfviper.New()
	// Configure the parent command's config bindings after the plugin has bound its flagset.
	_ = containerConfig.BindPFlag("logfile", cmd.PersistentFlags().Lookup("logfile"))
	_ = containerConfig.BindPFlag("loglevel", cmd.PersistentFlags().Lookup("loglevel"))
	_ = containerConfig.BindPFlag("artifacts", cmd.PersistentFlags().Lookup("artifacts"))
	containerConfig.SetDefault("logfile", DefaultLogFile)
	containerConfig.SetDefault("loglevel", DefaultLogLevel)
	containerConfig.SetDefault("artifacts", artifacts.DefaultArtifactsDir)

	// Build out the Container Plugin
	cmd.AddCommand(containerPlugin(ctx, containerConfig))
	// cmd.Hidden = true
	return cmd
}

// containerPlugin explicitly calls the check-container plugin. This should only
// be used for backwards compatibility purposes.
func containerPlugin(ctx context.Context, config *viper.Viper) *cobra.Command {
	// TODO(Jose): This is hard coded to depend on the name of the container check to be check-container
	plug := plugin.RegisteredPlugins()["check-container"]
	plcmd := plugin.NewCommand(ctx, config, "check-container", plug)
	plcmd.RunE = func(cmd *cobra.Command, args []string) error {
		return run(ctx, args, "check-container", config, &types.ResultWriterFile{})
	}
	plcmd.Use = "container"
	return plcmd
}
