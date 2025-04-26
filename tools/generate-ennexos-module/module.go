package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"math"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/hansmi/modbus-exporter-modules/internal/ennexos"
	"github.com/hansmi/modbus-exporter-modules/internal/modbusexporter"
)

var dataTypeToModule = map[ennexos.DataType]string{
	ennexos.Int16:  "int16",
	ennexos.Int32:  "int32",
	ennexos.Int64:  "int64",
	ennexos.Uint16: "uint16",
	ennexos.Uint32: "uint32",
	ennexos.Uint64: "uint64",
}

func prependFunctionCode(address uint32) uint32 {
	const functionCode = 4

	digits := int(math.Floor(1 + math.Log10(float64(address))))

	return (functionCode * uint32(math.Pow10(digits))) + address
}

func formatParameterHelp(params []parameter) (string, error) {
	type kv struct {
		key   string
		value string
	}

	var extra []kv

	first := params[0]

	if value := first.Def.Group; value != "" {
		extra = append(extra, kv{"Group", value})
	}

	if value := first.Def.Channel; value != "" {
		extra = append(extra, kv{"Channel", value})
	}

	if fixed := first.Def.DataFormat.Fixed; fixed != nil {
		scale := func(v int64) float64 {
			return float64(v)
		}

		if fixed.DecimalShift != 0 {
			factor := math.Pow10(-fixed.DecimalShift)

			scale = func(v int64) float64 {
				return factor * float64(v)
			}
		}

		if fixed.Unit != nil {
			extra = append(extra, kv{"Unit", *fixed.Unit})
		}

		if fixed.Default != nil {
			extra = append(extra, kv{
				"Default",
				fmt.Sprintf("%.6g", scale(*fixed.Default)),
			})
		}

		if !(fixed.LowerBound == nil || fixed.UpperBound == nil) {
			extra = append(extra, kv{
				"Range",
				fmt.Sprintf("%.6g to %.6g", scale(*fixed.LowerBound), scale(*fixed.UpperBound)),
			})
		}
	} else if taglist := first.Def.DataFormat.Taglist; taglist != nil {
		if taglist.Default != nil {
			extra = append(extra, kv{
				"Default",
				strconv.FormatInt(*taglist.Default, 10),
			})
		}

		if len(taglist.Values) > 0 {
			for _, value := range slices.Sorted(maps.Keys(taglist.Values)) {
				extra = append(extra, kv{
					fmt.Sprintf("Value %d", value),
					taglist.Values[value],
				})
			}
		}
	}

	var buf strings.Builder

	name := strings.TrimRightFunc(first.Def.Name, func(r rune) bool {
		return unicode.IsSpace(r) || r == '.'
	})

	buf.WriteString(name)
	buf.WriteRune('.')

	if len(extra) > 0 {
		buf.WriteString("\n\n")

		for _, i := range extra {
			buf.WriteString(i.key)
			buf.WriteString(": ")
			buf.WriteString(i.value)
			buf.WriteRune('\n')
		}
	}

	return buf.String(), nil
}

func parametersToMetricDefs(cfg config, namePrefix string, params []parameter) ([]modbusexporter.ModuleMetricDef, error) {
	var result []modbusexporter.ModuleMetricDef

	name := modbusexporter.SanitizeMetricName(fmt.Sprintf("%s_%s",
		namePrefix, strings.ToLower(params[0].Name)))

	metricType := cfg.metricType(params[0])

	labelsToParam := map[string]parameter{}

	for idx, p := range params {
		if serializedLabels, err := json.Marshal(p.Labels); err != nil {
			return nil, err
		} else if earlier, ok := labelsToParam[string(serializedLabels)]; ok {
			return nil, fmt.Errorf("labels not unique: %q already used by %#v", serializedLabels, earlier)
		} else {
			labelsToParam[string(serializedLabels)] = p
		}

		def := modbusexporter.ModuleMetricDef{
			Name:       name,
			Labels:     p.Labels.toMap(),
			Address:    prependFunctionCode(p.Def.Address),
			Endianness: "big",
			MetricType: metricType,
		}

		if dt, ok := dataTypeToModule[p.Def.DataType]; !ok {
			return nil, fmt.Errorf("unsupported data type %q", p.Def.DataType)
		} else {
			def.DataType = dt
		}

		if fixed := p.Def.DataFormat.Fixed; fixed != nil && fixed.DecimalShift != 0 {
			def.Factor = math.Pow10(-fixed.DecimalShift)
		}

		if idx == 0 {
			if help, err := formatParameterHelp(params); err != nil {
				return nil, err
			} else {
				def.Help = help
			}
		}

		result = append(result, def)
	}

	return result, nil
}

func parametersToMetrics(cfg config, prefix string, params []parameter) ([]modbusexporter.ModuleMetricDef, error) {
	byName := map[string][]parameter{}

	filtered, err := cfg.filter(params)
	if err != nil {
		return nil, err
	}

	for _, p := range filtered {
		byName[p.Name] = append(byName[p.Name], p)
	}

	var result []modbusexporter.ModuleMetricDef

	for _, name := range slices.Sorted(maps.Keys(byName)) {
		defs, err := parametersToMetricDefs(cfg, prefix, byName[name])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}

		result = append(result, defs...)
	}

	return result, nil
}
