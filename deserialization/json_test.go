package deserialization

import "testing"
import pb "github.com/raystack/raccoon/proto"

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
				b: []byte(`{
					"reqGuid": "17e2ac19-df8b-4a30-b111-fd7f5073d2f5",
					"sentTime": "2023-08-17T05:38:49.234986Z",
					"events": [
					  {
						"eventBytes": "eyJyYW5kb20xIjogImFiYyIsICJ4eXoiOiAxfQ==",
						"type": "topic 1"
					  }
					]
				  }`),
				i: &pb.SendEventRequest{},
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
