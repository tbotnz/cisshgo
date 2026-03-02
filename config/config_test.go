package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadInventory(t *testing.T) {
	content := `---
devices:
  - platform: csr1000v
    count: 10
  - platform: iosxr
    count: 5
`
	tmpFile := filepath.Join(t.TempDir(), "inventory.yaml")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	inv, err := LoadInventory(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	if len(inv.Devices) != 2 {
		t.Fatalf("Devices len = %d, want 2", len(inv.Devices))
	}
	if inv.Devices[0].Platform != "csr1000v" || inv.Devices[0].Count != 10 {
		t.Errorf("Devices[0] = %+v, want {csr1000v 10}", inv.Devices[0])
	}
}

func TestLoadInventory_MissingFile(t *testing.T) {
	_, err := LoadInventory("/nonexistent/inventory.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadInventory_InvalidYAML(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "bad.yaml")
	os.WriteFile(tmpFile, []byte(`not: valid: yaml: [[[`), 0644)
	_, err := LoadInventory(tmpFile)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoadInventory_BothPlatformAndScenario(t *testing.T) {
	content := `---
devices:
  - platform: csr1000v
    scenario: csr1000v-add-interface
    count: 1
`
	tmpFile := filepath.Join(t.TempDir(), "inventory.yaml")
	os.WriteFile(tmpFile, []byte(content), 0644)
	_, err := LoadInventory(tmpFile)
	if err == nil {
		t.Error("expected error when both platform and scenario are set")
	}
}

func TestLoadInventory_NeitherPlatformNorScenario(t *testing.T) {
	content := `---
devices:
  - count: 1
`
	tmpFile := filepath.Join(t.TempDir(), "inventory.yaml")
	os.WriteFile(tmpFile, []byte(content), 0644)
	_, err := LoadInventory(tmpFile)
	if err == nil {
		t.Error("expected error when neither platform nor scenario is set")
	}
}

func TestLoadInventory_NegativeCount(t *testing.T) {
	content := `---
devices:
  - platform: csr1000v
    count: -1
`
	tmpFile := filepath.Join(t.TempDir(), "inventory.yaml")
	os.WriteFile(tmpFile, []byte(content), 0644)
	_, err := LoadInventory(tmpFile)
	if err == nil {
		t.Error("expected error for negative count")
	}
}
