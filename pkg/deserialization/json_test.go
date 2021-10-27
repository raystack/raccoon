package deserialization

import (
	"reflect"
	"testing"
)

func TestJSONDeserializer(t *testing.T) {
	tests := []struct {
		name string
		want Deserializer
	}{
		{
			name: "Creating new JSON Deserializer",
			want: DeserializeFunc(func(b []byte, i interface{}) error {
				return nil
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JSONDeserializer(); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("JSONDeserializer() = %v, want %v", got, tt.want)
			}
		})
	}
}
