package modbusexporter

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSanitizeMetricName(t *testing.T) {
	for _, tc := range []struct {
		value string
		want  string
	}{
		{},
		{
			value: "HelloWorld",
			want:  "HelloWorld",
		},
		{
			value: "foo...bar baz",
			want:  "foo_bar_baz",
		},
	} {
		t.Run(tc.value, func(t *testing.T) {
			got := SanitizeMetricName(tc.value)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Result diff (-want +got):\n%s", diff)
			}
		})
	}
}
