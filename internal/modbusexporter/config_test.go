package modbusexporter

import (
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestMarshalConfig(t *testing.T) {
	for _, tc := range []struct {
		name    string
		config  Config
		wantErr error
		want    []*regexp.Regexp
	}{
		{
			name: "empty",
			want: []*regexp.Regexp{
				regexp.MustCompile(`^modules:`),
			},
		},
		{
			name: "multiline strings",
			config: Config{
				Modules: []ModuleDef{
					{
						Name: "foo\nbar",
						Metrics: []ModuleMetricDef{
							{
								Help: "first\nsecond\n  indented\nthird\n",
							},
						},
					},
				},
			},
			want: []*regexp.Regexp{
				regexp.MustCompile(`(?m)^\s*-\s+name:\s*\|-$`),
				regexp.MustCompile(`(?m)^\s{2,6}help:\s*\|\n\s{4,8}first\n`),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := MarshalConfig(tc.config)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			for _, expr := range tc.want {
				if !expr.Match(got) {
					t.Errorf("MarshalConfig() result doesn't match %q:\n%s", expr.String(), got)
				}
			}
		})
	}
}
