package serialization

import (
	"reflect"
	"testing"
)

func TestJSONSerializer(t *testing.T) {
	tests := []struct {
		name string
		want Serializer
	}{
		{
			name: "Initilizing JSON serializer",
			want: SerializeFunc(func(m interface{}) ([]byte, error) {
				return nil, nil
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JSONSerializer(); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("JSONSerializer() = %v, want %v", got, tt.want)
			}
		})
	}
}
