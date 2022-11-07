package utils_test

import (
	"reflect"
	"testing"

	"google.golang.org/genproto/protobuf/field_mask"

	"protocol-adapter/internal/utils"
)

func Test_mask(t *testing.T) {
	type args struct {
		fm   *field_mask.FieldMask
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Mask",
			args: args{
				fm: &field_mask.FieldMask{Paths: []string{"a1.b1.c2", "a1.b3", "a3"}},
				data: []byte(
					`{
						"a1": {
							"b1": {
								"c1": ["c1", "c1"],
								"c2": "cn1"
							},
							"b2": "b2",
							"b3": ["b3"]
						}, 
						"a2": "a2",
						"a3": "a3"
					}`),
			},
			want:    []byte(`{"a1":{"b1":{"c2":"cn1"},"b3":["b3"]},"a3":"a3"}`),
			wantErr: false,
		},
		{
			name: "Invalid json",
			args: args{
				&field_mask.FieldMask{Paths: []string{}},
				[]byte(`{"key"}`),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.Mask(tt.args.fm, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("mask() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mask() got = %s, want %s", got, tt.want)
			}
		})
	}
}
