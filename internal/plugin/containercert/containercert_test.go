package containercert

import "testing"

func TestContainerCert(t *testing.T) {
	plugin := &plug{}
	if err := plugin.Run(); err != nil {
		t.Error("failed to run the plugin", err)
	}
}
