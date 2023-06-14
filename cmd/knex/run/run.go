package run

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/bombsimon/logrusr/v4"
	"github.com/go-logr/logr"
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
	config *spfviper.Viper,
) *cobra.Command {
	cmd := &cobra.Command{
		Use: "run",
	}

	cmd.PersistentFlags().String("logfile", "", "Where the execution logfile will be written. (env: PFLT_LOGFILE)")
	_ = config.BindPFlag("logfile", cmd.PersistentFlags().Lookup("logfile"))

	cmd.PersistentFlags().String("loglevel", "", "The verbosity of the preflight tool itself. Ex. warn, debug, trace, info, error. (env: PFLT_LOGLEVEL)")
	_ = config.BindPFlag("loglevel", cmd.PersistentFlags().Lookup("loglevel"))

	cmd.PersistentFlags().String("artifacts", "", "Where check-specific artifacts will be written. (env: PFLT_ARTIFACTS)")
	_ = config.BindPFlag("artifacts", cmd.PersistentFlags().Lookup("artifacts"))

	config.SetDefault("logfile", DefaultLogFile)
	config.SetDefault("loglevel", DefaultLogLevel)
	config.SetDefault("artifacts", artifacts.DefaultArtifactsDir)

	for plinvoke, pl := range plugin.RegisteredPlugins() {
		plcmd := plugin.NewCommand(ctx, config, plinvoke, pl)
		plcmd.RunE = func(cmd *cobra.Command, args []string) error {
			return run(args, ctx, plinvoke, config)
		}
		cmd.AddCommand(plcmd)
	}

	return cmd
}
func run(
	args []string,
	ctx context.Context,
	pluginName string,
	config *spfviper.Viper,
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
	} else {
		l.Infof("Failed to log to file, using default stderr")
	}
	if ll, err := logrus.ParseLevel(config.GetString("loglevel")); err == nil {
		l.SetLevel(ll)
	}

	logger := logrusr.New(l)
	ctx = logr.NewContext(ctx, logger)

	artifactsWriter, err := artifacts.NewFilesystemWriter(artifacts.WithDirectory(config.GetString("artifacts")))
	if err != nil {
		return err
	}
	ctx = artifacts.ContextWithWriter(ctx, artifactsWriter)

	// Make the configuration look preflight-ish
	config.SetEnvPrefix("pflt")
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(`-`, `_`))

	// Run the plugin
	plugin := plugin.RegisteredPlugins()[pluginName]
	fmt.Println("Running Plugin =>", plugin.Name(), plugin.Version())

	if err := plugin.Init(ctx, config, args); err != nil {
		fmt.Println("ERR problem running init", err)
		return err
	}

	if err := plugin.ExecuteChecks(ctx); err != nil {
		fmt.Println("ERR problem running ExecuteChecks", err)
		return err
	}

	results := plugin.Results(ctx)
	f, err := plugin.OpenFile("results.json")
	if err != nil {
		fmt.Println("ERR problem opening results file", err)
		return err
	}
	defer f.Close()
	out := io.MultiWriter(os.Stdout, f)

	textResults, err := formatAsText(ctx, results)
	if err != nil {
		fmt.Println("ERR converting results to text", err)
		return err
	}

	_, err = out.Write(textResults)
	if err != nil {
		fmt.Println("Err couldn't write output")
	}

	if config.GetBool("submit") {
		if err := plugin.Submit(ctx); err != nil {
			log.Println("Err submitting", err)
			return err
		}
	}

	return nil
}

type FormatterFunc = func(context.Context, types.Results) (response []byte, formattingError error)

// Just as poc formatter, borrowed from preflight's library docs
var formatAsText FormatterFunc = func(_ context.Context, r types.Results) (response []byte, formattingError error) {
	b := []byte{}
	for _, v := range r.Passed {
		t := v.ElapsedTime.Milliseconds()
		s := fmt.Sprintf("PASSED  %s in %dms\n", v.Name(), t)
		b = append(b, []byte(s)...)
	}
	for _, v := range r.Failed {
		t := v.ElapsedTime.Milliseconds()
		s := fmt.Sprintf("FAILED  %s in %dms\n", v.Name(), t)
		b = append(b, []byte(s)...)
	}
	for _, v := range r.Errors {
		t := v.ElapsedTime.Milliseconds()
		s := fmt.Sprintf("ERRORED %s in %dms\n", v.Name(), t)
		b = append(b, []byte(s)...)
	}

	return b, nil
}
