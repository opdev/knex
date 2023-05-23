package containercert

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/redhat-openshift-ecosystem/knex/plugin"
	"github.com/redhat-openshift-ecosystem/knex/types"
	"github.com/spf13/viper"
)

func init() {
	plugin.Register("check-container", NewPlugin())
}

type plug struct {
	fileWriter
}

func NewPlugin() *plug {
	return &plug{}
}

func (p *plug) Register() error {
	return nil
}

func (p *plug) Run() error {
	fmt.Printf("%s is Running\n", p.Name())
	return nil
}

func (p *plug) Name() string {
	return "Container Certification"
}

func (p *plug) Init(cfg *viper.Viper) error {
	return nil
}

var vers = semver.MustParse("0.0.1")

func (p *plug) Version() semver.Version {
	return *vers
}

func (p *plug) ExecuteChecks(_ context.Context) error {
	fmt.Println("Execute Checks Called")
	return nil
}

func (p *plug) Results(_ context.Context) types.Results {
	return types.Results{}
}

func (p *plug) Submit(_ context.Context) error {
	fmt.Println("Submit called")
	return nil
}
