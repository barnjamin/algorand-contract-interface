package generator

import "testing"

func TestNew(t *testing.T) {
	tc, err := NewTemplateContract("../contracts/simple/test.tmpl.teal")
	if err != nil {
		t.Fatalf("Faied to parse contract template: %+v", err)
	}

	t.Logf("%+v", tc)
}
