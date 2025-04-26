package main

import (
	"bytes"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/hansmi/modbus-exporter-modules/internal/ennexos"
	"github.com/hansmi/modbus-exporter-modules/internal/modbusexporter"
)

func TestTransform(t *testing.T) {
	var buf bytes.Buffer

	var cfg config
	var pf ennexos.ParameterFile

	if err := transform(&buf, "test", cfg, pf); err != nil {
		t.Errorf("transform() failed: %v", err)
	}

	var got modbusexporter.Config

	if err := yaml.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Errorf("Unmarshaling module failed: %v", err)
	}
}
