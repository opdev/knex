package plugin

import (
	"context"

	"github.com/spf13/cobra"
)

func NewCommand(
	ctx context.Context,
	invocation string,
	pl Plugin,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:     invocation,
		Short:   pl.Name() + pl.Version().String(),
		Version: pl.Version().String(),
	}

	pl.BindFlags(cmd.Flags())

	return cmd
}
