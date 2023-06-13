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
		Short:   pl.Name() + pl.Version().String(),
		Version: pl.Version().String(),
	}

	pl.BindFlags(cmd.Flags())
	if err := config.BindPFlags(cmd.Flags()); err != nil {
		fmt.Println("unable to bind environment variables", err)
	}

	return cmd
}
