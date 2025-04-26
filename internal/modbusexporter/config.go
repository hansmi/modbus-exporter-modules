package modbusexporter

import (
	"fmt"
	"regexp"

	"github.com/goccy/go-yaml"
)

const TCPIP = "tcp/ip"

var endOfLineSpace = regexp.MustCompile(`(?m)[[:blank:]]+$`)

type Config struct {
	Modules []ModuleDef `yaml:"modules"`
}

type ModuleDef struct {
	Name     string            `yaml:"name"`
	Protocol string            `yaml:"protocol"`
	Metrics  []ModuleMetricDef `yaml:"metrics"`
}

type ModuleMetricDef struct {
	Name       string            `yaml:"name"`
	Help       string            `yaml:"help,omitempty"`
	Labels     map[string]string `yaml:"labels,omitempty"`
	Address    uint32            `yaml:"address"`
	DataType   string            `yaml:"dataType"`
	Endianness string            `yaml:"endianness,omitempty"`
	MetricType string            `yaml:"metricType"`
	Factor     float64           `yaml:"factor,omitempty"`
}

func MarshalConfig(c Config) ([]byte, error) {
	buf, err := yaml.MarshalWithOptions(
		c,
		yaml.AutoInt(),
		yaml.Flow(false),
		yaml.Indent(2),
		yaml.IndentSequence(false),
		yaml.UseLiteralStyleIfMultiline(true),
	)
	if err != nil {
		return nil, fmt.Errorf("marshaling config: %w", err)
	}

	return endOfLineSpace.ReplaceAllLiteral(buf, nil), nil
}
