package deserialization

import (
	"testing"

	pb "github.com/odpf/raccoon/proto"
)

func TestProtoDeserilizer_Deserialize(t *testing.T) {
	type args struct {
		b []byte
		i interface{}
	}
	tests := []struct {
		name    string
		d       *ProtoDeserilizer
		args    args
		wantErr bool
	}{
		{
			name: "Deserialize a proto message",
			d:    &ProtoDeserilizer{},
			args: args{
				b: []byte{},
				i: &pb.SendEventRequest{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ProtoDeserilizer{}
			if err := d.Deserialize(tt.args.b, tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("ProtoDeserilizer.Deserialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
