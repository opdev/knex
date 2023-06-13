package plugin

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCommand(
	ctx context.Context,
	config *viper.Viper,
	invocation string,
	pl Plugin,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:     invocation,
		Short:   fmt.Sprintf("%s at verions %s", pl.Name(), pl.Version().String()),
		Version: pl.Version().String(),
	}

	pl.BindFlags(cmd.Flags())
	if err := config.BindPFlags(cmd.Flags()); err != nil {
		// TODO(komish): This throwing an error is problematic at the moment because
		// we don't return an error, but we certainly can and just need to do additional parsing.
		// at the point of call.
		fmt.Println("unable to bind environment variables", err)
	}

	return cmd
}
