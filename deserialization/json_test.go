package deserialization

import "testing"

func TestJSONDeserializer_Deserialize(t *testing.T) {
	type args struct {
		b []byte
		i interface{}
	}
	tests := []struct {
		name    string
		j       *JSONDeserializer
		args    args
		wantErr bool
	}{
		{
			name: "Use JSON Deserializer",
			j:    &JSONDeserializer{},
			args: args{
				b: []byte(`{"A": "a"}`),
				i: &struct {
					A string
				}{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JSONDeserializer{}
			if err := j.Deserialize(tt.args.b, tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("JSONDeserializer.Deserialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
