package serialization

import (
	"reflect"
	"testing"
)

func TestProtoDeserilizer(t *testing.T) {
	tests := []struct {
		name string
		want Serializer
	}{
		{
			name: "initializing Proto Desrializer",
			want: SerializeFunc(func(m interface{}) ([]byte, error) {
				return nil, nil
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProtoSerilizer(); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("ProtoDeserilizer() = %v, want %v", got, tt.want)
			}
		})
	}
}
