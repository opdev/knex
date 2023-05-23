package run

import (
	"context"
	"fmt"

	"github.com/redhat-openshift-ecosystem/knex/plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Newcommand(
	ctx context.Context,
	config *viper.Viper,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:       "run",
		Short:     "Run a Certification Plugin",
		ValidArgs: validArgs(),
		Args: cobra.MatchAll(
			cobra.ExactArgs(1),
			cobra.OnlyValidArgs,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(args[0])
		},
	}

	return cmd
}

func validArgs() []string {
	registered := plugin.RegisteredPlugins()
	validArgs := make([]string, 0, len(registered))
	for k, _ := range registered {
		validArgs = append(validArgs, k)
	}
	return validArgs
}

func run(pluginName string) error {
	fmt.Println("Run invoked with plugin:", pluginName)
	plugin := plugin.RegisteredPlugins()[pluginName]
	return plugin.Run()
}
