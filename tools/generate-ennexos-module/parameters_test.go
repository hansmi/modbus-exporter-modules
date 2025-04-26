package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParameterLabelsToMap(t *testing.T) {
	zero := 0

	for _, tc := range []struct {
		name   string
		labels parameterLabels
		want   map[string]string
	}{
		{name: "empty"},
		{
			name: "string",
			labels: parameterLabels{
				String: &zero,
			},
			want: map[string]string{
				"string": "0",
			},
		},
		{
			name: "phase",
			labels: parameterLabels{
				Phase: 1,
			},
			want: map[string]string{
				"phase": "1",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.labels.toMap()

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Result diff (-want +got):\n%s", diff)
			}
		})
	}
}
