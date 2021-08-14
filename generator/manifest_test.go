package generator

import (
	"encoding/json"
	"testing"
)

const simple_path = "../contracts/simple"

func TestNewManifest(t *testing.T) {
	_, err := NewManifest(simple_path)
	if err != nil {
		t.Fatalf("Failed to create manifest: %+v", err)
	}
}

func TestGenerateInterface(t *testing.T) {
	m, err := NewManifest(simple_path)
	if err != nil {
		t.Fatalf("Failed to create manifest: %+v", err)
	}

	ci, err := m.GenerateInterface()
	if err != nil {
		t.Fatalf("Failed to generate interface: %+v", err)
	}

	b, err := json.MarshalIndent(ci, "", " ")
	if err != nil {
		t.Fatalf("Failed to marshal message: %+v", err)
	}

	t.Logf("%s", b)
}
