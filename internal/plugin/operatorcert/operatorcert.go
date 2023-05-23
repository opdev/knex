package operatorcert

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/redhat-openshift-ecosystem/knex/plugin"
	"github.com/spf13/viper"
)

func init() {
	plugin.Register("check-operator", NewPlugin())
}

type plug struct{}

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
	return "Operator Certification"
}

func (p *plug) Init(cfg *viper.Viper) error {
	return nil
}

var vers = semver.MustParse("0.0.1")

func (p *plug) Version() semver.Version {
	return *vers
}
