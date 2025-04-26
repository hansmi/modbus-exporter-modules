package modbusexporter

import "regexp"

var metricNameSanitizeRe = regexp.MustCompile(`[^_a-zA-Z0-9]+`)

func SanitizeMetricName(name string) string {
	return metricNameSanitizeRe.ReplaceAllLiteralString(name, "_")
}
