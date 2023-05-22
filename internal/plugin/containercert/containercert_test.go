package containercert

import "testing"

func TestContainerCert(t *testing.T) {
	plugin := &containerCertificationPlugin{}
	if err := plugin.Run("quay.io/opdev/simple-demo-operator:latest"); err != nil {
		t.Error("failed to run the plugin", err)
	}
}
