// Package types contains all types relevant to this PoC.
//
// This is organized into a single place just for PoC purposes.
// These are copied from preflight because preflight contains these in
// an internal package.
package types

import (
	"context"
	"io"
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// CheckEngine defines the functionality necessary to run all checks for a policy,
// and return the results of that check execution.
type CheckEngine interface {
	// ExecuteChecks should execute all checks in a policy and internally
	// store the results. Errors returned by ExecuteChecks should reflect
	// errors in pre-validation tasks, and not errors in individual check
	// execution itself.
	ExecuteChecks(context.Context) error
	// Results returns the outcome of executing all checks.
	Results(context.Context) Results
}

// Check as an interface containing all methods necessary
// to use and identify a given check.
type Check interface {
	// Validate will test the provided image and determine whether the
	// image complies with the check's requirements.
	Validate(ctx context.Context, imageReference ImageReference) (result bool, err error)
	// Name returns the name of the check.
	Name() string
	// Metadata returns the check's metadata.
	Metadata() Metadata
	// Help return the check's help information
	Help() HelpText
}

// ImageReference holds all things image-related
type ImageReference struct {
	ImageURI        string
	ImageFSPath     string
	ImageInfo       v1.Image
	ImageRepository string
	ImageRegistry   string
	ImageTagOrSha   string
}

type Result struct {
	Check
	ElapsedTime time.Duration
}

type Results struct {
	TestedImage       string
	PassedOverall     bool
	TestedOn          OpenshiftClusterVersion
	CertificationHash string
	Passed            []Result
	Failed            []Result
	Errors            []Result
}

// Metadata contains useful information regarding the check.
type Metadata struct {
	// Description contains a brief text detailing the overall goal of the check.
	Description string `json:"description" xml:"description"`
	// Level describes the certification level associated with the given check.
	//
	// TODO: define this more explicitly when requirements surrounding this metadata
	// text.
	Level string `json:"level" xml:"level"`
	// KnowledgeBaseURL is a URL detailing how to resolve a check failure.
	KnowledgeBaseURL string `json:"knowledge_base_url,omitempty" xml:"knowledgeBaseURL"`
	// CheckURL is a URL pointing to the official policy documentation from Red Hat, containing
	// information on exactly what is being tested and why.
	CheckURL string `json:"check_url,omitempty" xml:"checkURL"`
}

// HelpText is the help message associated with any given check
type HelpText struct {
	// Message is text provided to the user indicating where they should look
	// to find out why they failed or encountered an error in validation.
	Message string `json:"message" xml:"message"`
	// Suggestion is text provided to the user indicating what might need to
	// change in order to pass a check.
	Suggestion string `json:"suggestion" xml:"suggestion"`
}

type OpenshiftClusterVersion struct {
	Name    string
	Version string
}

func UnknownOpenshiftClusterVersion() OpenshiftClusterVersion {
	return OpenshiftClusterVersion{
		Name:    "unknown",
		Version: "unknown",
	}
}

// ResponseFormatter describes the expected methods a formatter
// must implement.
type ResponseFormatter interface {
	// PrettyName is the name used to represent this formatter.
	PrettyName() string
	// FileExtension represents the file extension one might use when creating
	// a file with the contents of this formatter.
	FileExtension() string
	// Format takes Results, formats it as needed, and returns the formatted
	// results ready to write as a byte slice.
	Format(context.Context, Results) (response []byte, formattingError error)
}

// ResultWriter defines methods associated with writing check results.
type ResultWriter interface {
	OpenFile(name string) (io.WriteCloser, error)
	io.WriteCloser
}

// ResultSubmitter defines methods associated with submitting results to Red HAt.
type ResultSubmitter interface {
	Submit(context.Context) error
}
