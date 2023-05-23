package containercert

import (
	"io"
	"os"
)

// Original Source:
// https://github.com/redhat-openshift-ecosystem/openshift-preflight/blob/main/internal/runtime/result_writer.go
// Note(Jose): I don't think this needs to live here. This is an application
// dependency, so this probably should be passed into the plugin as a dependency, or
// invoked at the top level.

// fileWriter implements a ResultWriter for use at preflight runtime.
type fileWriter struct {
	file *os.File
}

// OpenFile will open the expected file for writing.
func (f *fileWriter) OpenFile(name string) (io.WriteCloser, error) {
	file, err := os.OpenFile(
		name,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0o600)
	if err != nil {
		return nil, err
	}

	f.file = file // so we can close it later.
	return f, nil
}

func (f *fileWriter) Close() error {
	return f.file.Close()
}

func (f *fileWriter) Write(p []byte) (int, error) {
	return f.file.Write(p)
}
