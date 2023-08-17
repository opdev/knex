# knex

Pluggable Certification.

## Building and Registering a Plugin

- Implement the [Plugin](./plugin/v0/plugin.go#L44)
- Create an `init` function somewhere in your plugin codebase that calls the [Register](./plugin/aliases.go#L8) function
- Submit a PR to the repository adding a blank-initialization of your plugin code (ex. [here](./plugin/registration/add_plugins_here.go#L8))
- Ensure the go.mod value for your plugin points to your version. The repository
  encourages semantic versioning, and will represent plugin versions to users.
  You are encouraged to ensure version changes are associated with new behaviors.

## Writing Logs and Artifacts

- Knex will pass a logger and an artifact writer to your plugin via the
  `context`.
- For the logger, utilize the logr helper function
  [FromContextOrDiscard](https://pkg.go.dev/github.com/go-logr/logr#FromContextOrDiscard)
  (or equivalents).
- For ArtifactsWriter, utilize the helper function
  [WriterFromContext](https://pkg.go.dev/github.com/redhat-openshift-ecosystem/openshift-preflight/artifacts#WriterFromContext)
- Plugins are generally discouraged from reconfiguring the included logger or artifact writer.

## Binding Environment/Flags

- A plugin will be passed a
  [pflag.FlagSet](https://pkg.go.dev/github.com/spf13/pflag#FlagSet) to its
  `BindFlags` method. It should bind all flags necessary for the plugin to
  operate at this time.
- These flags are converted to environment variables using viper's
  [AutomaticEnv](https://pkg.go.dev/github.com/spf13/viper#AutomaticEnv). Users
  will need to prefix your environment variables with `PFLT_`.
    - Dashes are converted to hyphens.
    - Other non-env special characters are not supported (e.g. period)
