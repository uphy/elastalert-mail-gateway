package jsonutil

import (
	"reflect"
	"testing"
)

func TestObject_Get(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		d    Object
		args args
		want *Value
	}{
		{
			name: "simple string",
			d: Object{
				"a": "AAA",
			},
			args: args{
				name: "a",
			},
			want: NewValue("AAA"),
		},
		{
			name: "simple int",
			d: Object{
				"a": 1,
			},
			args: args{
				name: "a",
			},
			want: NewValue(1),
		},
		{
			name: "not found",
			d: Object{
				"a": "AAA",
			},
			args: args{
				name: "b",
			},
			want: nil,
		},
		{
			name: "dotted-name found",
			d: Object{
				"a": map[string]interface{}{
					"b": "AAA",
				},
			},
			args: args{
				name: "a.b",
			},
			want: NewValue("AAA"),
		},
		{
			name: "dotted-name found int",
			d: Object{
				"a": map[string]interface{}{
					"b": 1,
				},
			},
			args: args{
				name: "a.b",
			},
			want: NewValue(1),
		},
		{
			name: "dotted-name found 2",
			d: Object{
				"a": Object{
					"b": "AAA",
				},
			},
			args: args{
				name: "a.b",
			},
			want: NewValue("AAA"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Get(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Object.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
