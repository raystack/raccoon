package deserialization

import "testing"

func TestJSONDeserializer_Deserialize(t *testing.T) {
	type args struct {
		b []byte
		i interface{}
	}
	tests := []struct {
		name    string
		j       DeserializeFunc
		args    args
		wantErr bool
	}{
		{
			name: "Use JSON Deserializer",
			j:    DeserializeJSON,
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
			if err := tt.j(tt.args.b, tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("JSONDeserializer.Deserialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
