package serialization

import (
	"fmt"
	"reflect"
	"testing"

	pb "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/raccoon/v1beta1"
)

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
			fmt.Println(string(got))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONSerializer.Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}
