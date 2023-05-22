package plugin

import (
	"fmt"
	"regexp"
)

var registeredPlugins map[string]Plugin = make(map[string]Plugin)
var isNormalized = regexp.MustCompile(`^[a-z][a-z\-]+[a-z]$`).MatchString

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
	Name() string
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