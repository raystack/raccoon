package serde

import (
	"reflect"
	"testing"

	pb "github.com/raystack/raccoon/proto"
)

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

func TestJSONSerializer_Serialize(t *testing.T) {
	type args struct {
		m interface{}
	}
	tests := []struct {
		name    string
		s       SerializeFunc
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Serialize JSON",
			s:    SerializeJSON,
			args: args{
				m: &pb.SendEventRequest{},
			},
			want:    []byte{123, 125},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s(tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONSerializer.Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONSerializer.Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}
