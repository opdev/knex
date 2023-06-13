package run

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/opdev/knex/plugin/v0"
	"github.com/opdev/knex/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCommand(
	ctx context.Context,
	config *viper.Viper,
) *cobra.Command {
	cmd := &cobra.Command{
		Use: "run",
	}

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
	config *viper.Viper,
) error {
	plugin := plugin.RegisteredPlugins()[pluginName]
	fmt.Println("Running Plugin =>", plugin.Name(), plugin.Version())

	// Make the configuration look preflight-ish
	config.SetEnvPrefix("pflt")
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(`-`, `_`))

	if err := plugin.Init(config, args); err != nil {
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
var formatAsText FormatterFunc = func(ctx context.Context, r types.Results) (response []byte, formattingError error) {
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
