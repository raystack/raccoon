package deserialization

import (
	"testing"

	pb "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/raccoon/v1beta1"
)

func TestProtoDeserilizer_Deserialize(t *testing.T) {
	type args struct {
		b []byte
		i interface{}
	}
	tests := []struct {
		name    string
		d       DeserializeFunc
		args    args
		wantErr bool
	}{
		{
			name: "Deserialize a proto message",
			d:    DeserializeProto,
			args: args{
				b: []byte{},
				i: &pb.SendEventRequest{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d(tt.args.b, tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("ProtoDeserilizer.Deserialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
