package containercert

import "testing"

func TestContainerCert(t *testing.T) {
	plugin := &containerCertificationPlugin{}
	if err := plugin.Run(); err != nil {
		t.Error("failed to run the plugin", err)
	}
}
