// Package registration contains all plugins to register.
// Plugin developers should blank-initialize their plugins here
package registration

import (
	// Plugin initialization

	_ "github.com/opdev/container-certification/plugin"
	_ "github.com/opdev/container-certification/plugin/rootexception"
	_ "github.com/opdev/container-certification/plugin/scratchexception"

	_ "github.com/opdev/plugin-template/plugin"
)
