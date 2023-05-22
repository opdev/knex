package containercert

import (
	"fmt"

	"github.com/redhat-openshift-ecosystem/knex/plugin"
)

func init() {
	plugin.Register("check-container", NewPlugin())
}

type containerCertificationPlugin struct{}

func NewPlugin() *containerCertificationPlugin {
	return &containerCertificationPlugin{}
}

func (p *containerCertificationPlugin) Register() error {
	return nil
}

func (p *containerCertificationPlugin) Run(image string) error {
	fmt.Println("The image is:", image)
	return nil
}

func (p *containerCertificationPlugin) Name() string {
	return "container-certification"
}
