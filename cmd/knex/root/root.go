// Package r0ot implements the command-line interface.
package root

import (
	"context"

	// Plugin initialization
	_ "github.com/opdev/container-certification/plugin"
	_ "github.com/opdev/helm-certification/plugin"
	_ "github.com/opdev/plugin-template/plugin"

	"github.com/opdev/knex/cmd/knex/listplugins"
	"github.com/opdev/knex/cmd/knex/run"
	"github.com/spf13/cobra"
)

func NewCommand(
	ctx context.Context,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "knex",
		Short:   "Pluggable Certification",
		Version: "0.0.0",
	}

	cmd.AddCommand(listplugins.NewCommand(ctx))
	cmd.AddCommand(run.NewCommand(ctx))

	return cmd
}

/*
// preRunConfig is used by cobra.PreRun in all non-root commands to load all necessary configurations
func preRunConfig(cmd *cobra.Command, args []string) {
	viper := viper.Instance()
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{DisableColors: true})

	// set up logging
	logname := viper.GetString("logfile")
	logFile, err := os.OpenFile(logname, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err == nil {
		mw := io.MultiWriter(os.Stderr, logFile)
		l.SetOutput(mw)
	} else {
		l.Infof("Failed to log to file, using default stderr")
	}
	if ll, err := logrus.ParseLevel(viper.GetString("loglevel")); err == nil {
		l.SetLevel(ll)
	}

	// if we are in the offline flow redirect log file to exist in the directory where all other artifact exist
	if viper.GetBool("offline") {
		// Get the base name of the logfile, in case logfile has a path
		baseLogName := filepath.Base(logname)
		artifacts := viper.GetString("artifacts")

		// ignoring error since OpenFile will error and we'll still have the multiwriter from above
		_ = os.Mkdir(artifacts, 0o777)

		artifactsLogFile, err := os.OpenFile(filepath.Join(artifacts, baseLogName), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
		if err == nil {
			mw := io.MultiWriter(os.Stderr, logFile, artifactsLogFile)
			l.SetOutput(mw)
		}

		// setting log level to trace, to provide the most detailed logs possible
		l.SetLevel(logrus.TraceLevel)
	}

	if !configFileUsed {
		l.Debug("config file not found, proceeding without it")
	}

	logger := logrusr.New(l)
	ctx := logr.NewContext(cmd.Context(), logger)
	cmd.SetContext(ctx)
}
*/
