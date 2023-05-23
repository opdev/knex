package operatorcert

import (
	"fmt"

	"github.com/redhat-openshift-ecosystem/knex/plugin"
)

func init() {
	plugin.Register("check-operator", NewPlugin())
}

type operatorCertificationPlugin struct{}

func NewPlugin() *operatorCertificationPlugin {
	return &operatorCertificationPlugin{}
}

func (p *operatorCertificationPlugin) Register() error {
	return nil
}

func (p *operatorCertificationPlugin) Run() error {
	fmt.Printf("%s is Running\n", p.Name())
	return nil
}

func (p *operatorCertificationPlugin) Name() string {
	return "Operator Certification"
}
