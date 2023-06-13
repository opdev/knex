package plugin

import (
	"github.com/go-logr/logr"
	"github.com/redhat-openshift-ecosystem/openshift-preflight/artifacts"
)

// RuntimeConfiguration includes common IO components to be leveraged by
// plugins.
type RuntimeConfiguration struct {
	// Logger is a preconfigured logger, writing to the expected places
	// at the expected verbosity.
	// TODO(Jose): Should this be using an interface? Should we only expose a limited
	// subset of this functionality?
	Logger *logr.Logger
	// ArtifactWriter is a preconfigured artifact writer for writing arbitrary
	// artifacts.
	ArtifactWriter artifacts.ArtifactWriter
}
