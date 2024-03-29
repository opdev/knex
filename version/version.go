// Package version contains all identifiable versioning info for
// describing the knex project.
package version

import "fmt"

var (
	projectName = "github.com/opdev/knex"
	version     = "unknown"
	commit      = "unknown"
)

var Version = VersionContext{
	Name:    projectName,
	Version: version,
	Commit:  commit,
}

type VersionContext struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

func (vc *VersionContext) String() string {
	return fmt.Sprintf("%s <commit: %s>", vc.Version, vc.Commit)
}
