package operatorcert

import "testing"

func TestOperatorCert(t *testing.T) {
	plugin := &plug{}
	if err := plugin.Run(); err != nil {
		t.Error("failed to run the plugin", err)
	}
}
