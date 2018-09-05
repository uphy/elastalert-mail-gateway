package elastalert

import (
	"testing"
)

func TestParseMailBody(t *testing.T) {
	type args struct {
		body string
	}
	tests := []struct {
		name    string
		args    args
		want    MailBody
		wantErr bool
	}{
		{
			name: "simple",
			args: args{
				body: `alert text1
alert text2

@timestamp: 2018-08-31T01:23:45Z
_id: 1
_index: test
_type: _doc
a: 1
b: foo
c: [
	0,
	1
]
d: [
	{
		"bar": 2,
		"foo": 1
	},
	{
		"bar": 4,
		"foo": 3
	}
]
num_hits: 3
num_matches: 3

`,
			},
			want: MailBody{
				&MailBodyEntry{
					Body: `alert text1
alert text2`,
					Doc: map[string]interface{}{
						"@timestamp": "2018-08-31T01:23:45Z",
						"_id":        1,
						"_index":     "test",
						"_type":      "_doc",
						"a":          1,
						"b":          "foo",
						"c":          []interface{}{0, 1},
						"d": []interface{}{
							map[string]interface{}{
								"bar": 2,
								"foo": 1,
							},
							map[string]interface{}{
								"bar": 4,
								"foo": 3,
							},
						},
						"num_hits":    3,
						"num_matches": 3,
					},
				},
			},
		},
		{
			name: "two alerts",
			args: args{
				body: `alert1 text1
alert1 text2

@timestamp: 2018-08-31T01:23:45Z
_id: 1
_index: test
_type: _doc
a: 1
b: foo
c: [
	0,
	1
]
d: [
	{
		"bar": 2,
		"foo": 1
	},
	{
		"bar": 4,
		"foo": 3
	}
]
num_hits: 3
num_matches: 3

----------------------------------------
alert2 text1
alert2 text2

@timestamp: 2018-08-31T01:23:46Z
_id: 1
_index: test
_type: _doc
a: 1
b: foo
c: [
	0,
	1
]
d: [
	{
		"bar": 2,
		"foo": 1
	},
	{
		"bar": 4,
		"foo": 3
	}
]
num_hits: 3
num_matches: 3

----------------------------------------

`,
			},
			want: MailBody{
				&MailBodyEntry{
					Body: `alert1 text1
alert1 text2`,
					Doc: map[string]interface{}{
						"@timestamp": "2018-08-31T01:23:45Z",
						"_id":        1,
						"_index":     "test",
						"_type":      "_doc",
						"a":          1,
						"b":          "foo",
						"c":          []interface{}{0, 1},
						"d": []interface{}{
							map[string]interface{}{
								"bar": 2,
								"foo": 1,
							},
							map[string]interface{}{
								"bar": 4,
								"foo": 3,
							},
						},
						"num_hits":    3,
						"num_matches": 3,
					},
				},
				&MailBodyEntry{
					Body: `alert2 text1
alert2 text2`,
					Doc: map[string]interface{}{
						"@timestamp": "2018-08-31T01:23:46Z",
						"_id":        1,
						"_index":     "test",
						"_type":      "_doc",
						"a":          1,
						"b":          "foo",
						"c":          []interface{}{0, 1},
						"d": []interface{}{
							map[string]interface{}{
								"bar": 2,
								"foo": 1,
							},
							map[string]interface{}{
								"bar": 4,
								"foo": 3,
							},
						},
						"num_hits":    3,
						"num_matches": 3,
					},
				},
			},
		},
		{
			name: "empty body",
			args: args{
				body: `@timestamp: 2018-08-31T01:23:45Z
_id: 1
_index: test
_type: _doc
a: 1
b: foo
c: [
	0,
	1
]
d: [
	{
		"bar": 2,
		"foo": 1
	},
	{
		"bar": 4,
		"foo": 3
	}
]
num_hits: 3
num_matches: 3

`,
			},
			want: MailBody{
				&MailBodyEntry{
					Body: ``,
					Doc: map[string]interface{}{
						"@timestamp": "2018-08-31T01:23:45Z",
						"_id":        1,
						"_index":     "test",
						"_type":      "_doc",
						"a":          1,
						"b":          "foo",
						"c":          []interface{}{0, 1},
						"d": []interface{}{
							map[string]interface{}{
								"bar": 2,
								"foo": 1,
							},
							map[string]interface{}{
								"bar": 4,
								"foo": 3,
							},
						},
						"num_hits":    3,
						"num_matches": 3,
					},
				},
			},
		},
		{
			name: "no doc",
			args: args{
				body: `abcdefg`,
			},
			want: MailBody{
				&MailBodyEntry{
					Body: `abcdefg`,
					Doc:  map[string]interface{}{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMailBody(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMailBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			eq, desc := got.equalsTo(tt.want)
			if !eq {
				t.Error(desc)
			}
		})
	}
}
