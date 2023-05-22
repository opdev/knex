package listplugins

import (
	"context"
	"fmt"

	"github.com/redhat-openshift-ecosystem/knex/plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCommand(
	ctx context.Context,
	config *viper.Viper,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-plugins",
		Short: "list the plugins that have been registered",
		RunE:  func(cmd *cobra.Command, args []string) error { return listPlugins() },
	}

	return cmd
}

func listPlugins() error {
	fmt.Println("listing plugins")
	registered := plugin.RegisteredPlugins()
	for k, v := range registered {
		fmt.Printf(`Plugin %s is registered at entrypoint "%s"`, v.Name(), k)
	}

	return nil
}