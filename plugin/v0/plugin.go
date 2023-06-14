package plugin

import (
	"context"
	"fmt"
	"regexp"

	"github.com/Masterminds/semver/v3"
	"github.com/opdev/knex/types"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	registeredPlugins map[string]Plugin = make(map[string]Plugin)
	isNormalized                        = regexp.MustCompile(`^[a-z][a-z\-]+[a-z]$`).MatchString
)

func RegisteredPlugins() map[string]Plugin {
	rpCopy := map[string]Plugin{}
	for k, v := range registeredPlugins {
		// TODO(jose): double check this isn't a reference error in the map
		// assignment caused by the loop?
		rpCopy[k] = v
	}

	return rpCopy
}

// Register registers Plugin with name. If the plugin does not conform to the
// expected standards, this panics.
func Register(name string, plugin Plugin) {
	if err := ensurePluginNameMeetsStandards(name); err != nil {
		panic(err)
	}

	if err := ensurePluginNameIsUnique(name); err != nil {
		panic(err)
	}

	registeredPlugins[name] = plugin
}

type Plugin interface {
	// Init is called before all Execution, allowing a plugin to
	// configure itself informed by the configuration passed in by
	// the user, and by preflight itself.
	//
	// Args here will represent remaining positional arguments to be parsed
	// by the plugin. It may be a good practice to use the cobra PositionalArgs
	// function definition as a guideline for how to treat this slice.
	Init(
		// ctx will contain an initialized logger and artifacts writer at the time
		// of Init.
		ctx context.Context,
		// config contains the parsed flag/environment, including those flag values
		// provided by the plugin's BindFlags call.
		config *viper.Viper,
		// args contains positional arguments being passed to the plugin itself.
		args []string) error
	// Name identifies the plugin. Should be a formal definition
	// (e.g. "My Plugin")
	Name() string
	// Version is a semantic version representation for your plugin.
	Version() semver.Version
	// BindFlags binds a plugin's flags to the runtime execution.
	BindFlags(*pflag.FlagSet) *pflag.FlagSet
	// Run executes the plugin. Leaving commented for now. Using an arbitrary "run" method like this
	// may be worth considering if existing structured lements like the Check Engine don't work for this use case.
	// Run() error

	// Plumbing, allowing for standardized execution of a plugin.
	types.CheckEngine
	// Invoked if the user requested submissions.
	types.ResultSubmitter
	types.ResultWriter
}

func ensurePluginNameMeetsStandards(name string) error {
	// This is just an example validation.
	if !isNormalized(name) {
		return fmt.Errorf("invalid plugin name")
	}

	return nil
}

func ensurePluginNameIsUnique(name string) error {
	if _, exists := registeredPlugins[name]; exists {
		return fmt.Errorf("plugin already exists with name")
	}

	return nil
}

// Note(Jose): Plugin registration would be better if we could do it at compile
// time instead of at runtime to prevent shipping a binary that has plugin
// conflicts.
//
// Either that, or a definitive way to ensure that we're registering plugins
// to prevent release boo-boos.
