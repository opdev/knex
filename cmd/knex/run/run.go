package run

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/bombsimon/logrusr/v4"
	"github.com/go-logr/logr"
	"github.com/opdev/knex/formatters"
	"github.com/opdev/knex/plugin/v0"
	"github.com/opdev/knex/types"
	"github.com/redhat-openshift-ecosystem/openshift-preflight/artifacts"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	spfviper "github.com/spf13/viper"
)

const (
	DefaultLogFile  = "preflight.log"
	DefaultLogLevel = "info"
)

func NewCommand(
	ctx context.Context,
) *cobra.Command {
	cmd := &cobra.Command{
		Use: "run",
	}

	cmd.PersistentFlags().String("logfile", "", "Where the execution logfile will be written. (env: PFLT_LOGFILE)")
	cmd.PersistentFlags().String("loglevel", "", "The verbosity of the preflight tool itself. Ex. warn, debug, trace, info, error. (env: PFLT_LOGLEVEL)")
	cmd.PersistentFlags().String("artifacts", "", "Where check-specific artifacts will be written. (env: PFLT_ARTIFACTS)")
	cmd.PersistentFlags().BoolP("submit", "s", false, "Submit results to Red Hat if the called plugin supports it automated submission through this tool.")

	for i, p := range plugin.RegisteredPlugins() {
		invocation := i
		plug := p
		config := spfviper.New()
		plcmd := plugin.NewCommand(ctx, config, invocation, plug)
		plcmd.RunE = func(cmd *cobra.Command, args []string) error {
			return run(ctx, args, invocation, config, &types.ResultWriterFile{})
		}

		// Configure the parent command's config bindings after the plugin has bound its flagset.
		_ = config.BindPFlag("logfile", cmd.PersistentFlags().Lookup("logfile"))
		_ = config.BindPFlag("loglevel", cmd.PersistentFlags().Lookup("loglevel"))
		_ = config.BindPFlag("artifacts", cmd.PersistentFlags().Lookup("artifacts"))
		_ = config.BindPFlag("submit", cmd.PersistentFlags().Lookup("submit"))

		config.SetDefault("logfile", DefaultLogFile)
		config.SetDefault("loglevel", DefaultLogLevel)
		config.SetDefault("artifacts", artifacts.DefaultArtifactsDir)
		config.SetDefault("submit", false)

		cmd.AddCommand(plcmd)
	}

	return cmd
}

func run(
	ctx context.Context,
	args []string,
	pluginName string,
	config *spfviper.Viper,
	resultWriter types.ResultWriter,
) error {
	// Manage outputs on behalf of the plugin. This must happen before the
	// plugin init is called to prevent modifications to the viper configuration
	// that's passed to it from bubbling upward to preflight's scope.
	//
	// This is borrowed from preflight's check PreRunE with the intention of
	// stuffing the logger and artifacts writer in the context to maintain
	// compatibility with the existing container/operator certification.
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{DisableColors: true})

	logname := config.GetString("logfile")
	logFile, err := os.OpenFile(logname, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err == nil {
		mw := io.MultiWriter(os.Stderr, logFile)
		l.SetOutput(mw)
		defer logFile.Close()
	} else {
		l.Infof("Failed to log to file, using default stderr")
	}
	if ll, err := logrus.ParseLevel(config.GetString("loglevel")); err == nil {
		l.SetLevel(ll)
	}

	logger := logrusr.New(l)

	// Pass in a logger with addition keys/values to the plugin so we know the plugin emitted a log line.
	ctx = logr.NewContext(ctx, logger.WithValues("emitter", "plugin"))

	artifactsWriter, err := artifacts.NewFilesystemWriter(artifacts.WithDirectory(config.GetString("artifacts")))
	if err != nil {
		return err
	}
	ctx = artifacts.ContextWithWriter(ctx, artifactsWriter)

	// Make the configuration look preflight-ish
	config.SetEnvPrefix("pflt")
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(`-`, `_`))

	// Writing Results, also borrowed from Preflight (RunPreflight, specifically)
	// Fail early if we cannot write to the results path.
	// TODO(Jose): The preflight version of this handles formatters, etc. Stubbed this out to .txt for PoC
	resultsFilePath, err := artifactsWriter.WriteFile("results.json", strings.NewReader(""))
	if err != nil {
		return err
	}

	resultsFile, err := resultWriter.OpenFile(resultsFilePath)
	if err != nil {
		return err
	}

	defer resultsFile.Close()
	resultsOutputTarget := io.MultiWriter(os.Stdout, resultsFile)

	// Run the plugin
	plugin := plugin.RegisteredPlugins()[pluginName]
	logger.Info("Calling plugin", "name", plugin.Name(), "version", plugin.Version())

	if err := plugin.Init(ctx, config, args); err != nil {
		logger.Error(err, "unable to initialize plugin")
		return err
	}

	if err := plugin.ExecuteChecks(ctx); err != nil {
		logger.Error(err, "unable to execute checks")
		return err
	}

	results := plugin.Results(ctx)
	textResults, err := formatAsJSON(ctx, results)
	if err != nil {
		logger.Error(err, "unable to format results")
		return err
	}

	_, err = resultsOutputTarget.Write(textResults)
	if err != nil {
		logger.Error(err, "unable to write results")
	}

	if config.GetBool("submit") {
		if err := plugin.Submit(ctx); err != nil {
			logger.Error(err, "unable to call plugin submission")
			return err
		}
	}

	return nil
}

var formatAsJSON formatters.FormatterFunc = func(ctx context.Context, r types.Results) (response []byte, formattingError error) {
	f, _ := formatters.NewByName("json")
	return f.Format(ctx, r)
}
