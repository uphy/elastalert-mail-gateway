package condition

import (
	"testing"

	"github.com/uphy/elastalert-mail-gateway/pkgs/jsonutil"
)

func TestConditions_Eval(t *testing.T) {
	type fields struct {
		Match      *Match
		conditions []Condition
	}
	type args struct {
		scope Scope
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "match true",
			fields: fields{
				Match: &Match{
					KeyValue: jsonutil.KeyValue{
						Key:   "a",
						Value: jsonutil.NewValue(1),
					},
				},
			},
			args: args{
				scope: jsonutil.Object(map[string]interface{}{
					"a": 1,
				}),
			},
			want: true,
		},
		{
			name: "match false",
			fields: fields{
				Match: &Match{
					KeyValue: jsonutil.KeyValue{
						Key:   "a",
						Value: jsonutil.NewValue(1),
					},
				},
			},
			args: args{
				scope: jsonutil.Object(map[string]interface{}{
					"a": 2,
				}),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conditions{
				Match:      tt.fields.Match,
				conditions: tt.fields.conditions,
			}
			got, err := c.Eval(tt.args.scope)
			if (err != nil) != tt.wantErr {
				t.Errorf("Conditions.Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Conditions.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}
