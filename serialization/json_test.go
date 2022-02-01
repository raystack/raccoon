package serialization

import (
	"fmt"
	pb "raccoon/proto"
	"reflect"
	"testing"
)

func TestJSONSerializer_Serialize(t *testing.T) {
	type args struct {
		m interface{}
	}
	tests := []struct {
		name    string
		s       *JSONSerializer
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Serialize JSON",
			s:    &JSONSerializer{},
			args: args{
				m: &pb.EventRequest{},
			},
			want:    []byte{123, 125},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &JSONSerializer{}
			got, err := s.Serialize(tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONSerializer.Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(string(got))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONSerializer.Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}
