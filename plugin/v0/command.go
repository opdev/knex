package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCommand(
	config *viper.Viper,
	invocation string,
	pl Plugin,
) *cobra.Command {
	cmd := cobra.Command{
		Use:     invocation,
		Short:   fmt.Sprintf("%s at version %s", pl.Name(), pl.Version().String()),
		Version: pl.Version().String(),
	}

	cmd.Flags().AddFlagSet(pl.Flags())
	if err := config.BindPFlags(cmd.LocalFlags()); err != nil {
		// Note(komish): This panics to help preflight detect if flag binding will actually work. This is still
		// a runtime, check, though, and doesn't happen until we call the run subcommna.d
		panic(fmt.Sprintf("fatal error attempting to bind plugin flags for plugin %s: %s", pl.Name(), err))
	}

	return &cmd
}
