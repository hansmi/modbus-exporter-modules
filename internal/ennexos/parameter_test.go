package ennexos

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestReadParameterFile(t *testing.T) {
	for _, tc := range []struct {
		name    string
		content string
		path    string
		want    ParameterFile
		wantErr error
	}{
		{
			name:    "missing file",
			path:    filepath.Join(t.TempDir(), "missing"),
			wantErr: os.ErrNotExist,
		},
		{
			name:    "empty file",
			wantErr: cmpopts.AnyError,
		},
		{
			name:    "empty payload",
			content: "{}",
		},
		{
			name: "valid",
			content: `
{
  "header": {
    "description": "Header description"
  },
  "parameters": [
    {
      "access_read": true,
      "access_write": false,
      "address": 12345,
      "channel": "modbus.profile.revision",
      "data_format": {
        "fixed": {
          "decimal_shift": 1,
          "default": null
        }
      },
      "data_type": "u32",
      "group": "Type Label",
      "name": "Modbus profile revision"
    }
  ]
}
`,
			want: ParameterFile{
				Header: ParameterHeader{Description: "Header description"},
				Parameters: []Parameter{
					{
						Address:    12345,
						Name:       "Modbus profile revision",
						AccessRead: true,
						DataType:   Uint32,
						DataFormat: ParameterFormat{
							Fixed: &ParameterFormatFixed{
								DecimalShift: 1,
							},
						},
						Group:   "Type Label",
						Channel: "modbus.profile.revision",
					},
				},
			},
		},
		{
			name:    "bad data type",
			content: `{ "parameters": [ { "data_type": "bad invalid" } ] }`,
			wantErr: os.ErrInvalid,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := tc.path

			if path == "" {
				path = filepath.Join(t.TempDir(), "params")

				if err := os.WriteFile(path, []byte(tc.content), 0644); err != nil {
					t.Fatalf("WriteFile(%q) failed: %v", path, err)
				}
			}

			got, err := ReadParameterFile(path)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, *got); diff != "" {
					t.Errorf("Result diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
