package deserialization

import (
	"reflect"
	"testing"
)

func TestProtoDeserilizer(t *testing.T) {
	tests := []struct {
		name string
		want Deserializer
	}{
		{
			name: "Create new proto Deserializer",
			want: DeserializeFunc(func(b []byte, i interface{}) error {
				return nil
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProtoDeserilizer(); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("ProtoDeserilizer() = %v, want %v", got, tt.want)
			}
		})
	}
}
