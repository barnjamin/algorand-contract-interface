package generator

import "testing"

func TestNew(t *testing.T) {
	_, err := NewTemplateContract(simple_path + "/test.tmpl.teal")
	if err != nil {
		t.Fatalf("Faied to parse contract template: %+v", err)
	}
}
