package rest

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/raystack/raccoon/core/serde"
	pb "github.com/raystack/raccoon/proto"
)

func TestResponse_SetCode(t *testing.T) {
	type fields struct {
		SendEventResponse *pb.SendEventResponse
	}
	type args struct {
		code pb.Code
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Response
	}{
		{
			name: "sets response code",
			fields: fields{
				SendEventResponse: &pb.SendEventResponse{},
			},
			args: args{
				code: pb.Code_CODE_UNSPECIFIED,
			},
			want: &Response{
				SendEventResponse: &pb.SendEventResponse{
					Code: pb.Code_CODE_UNSPECIFIED,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Response{
				SendEventResponse: tt.fields.SendEventResponse,
			}
			if got := r.SetCode(tt.args.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Response.SetCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponse_SetStatus(t *testing.T) {
	type fields struct {
		SendEventResponse *pb.SendEventResponse
	}
	type args struct {
		status pb.Status
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Response
	}{
		{
			name: "set status",
			fields: fields{
				SendEventResponse: &pb.SendEventResponse{},
			},
			args: args{
				status: pb.Status_STATUS_SUCCESS,
			},
			want: &Response{
				SendEventResponse: &pb.SendEventResponse{
					Status: pb.Status_STATUS_SUCCESS,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Response{
				SendEventResponse: tt.fields.SendEventResponse,
			}
			if got := r.SetStatus(tt.args.status); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Response.SetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponse_SetSentTime(t *testing.T) {
	timeNow := time.Now().Unix()
	type fields struct {
		SendEventResponse *pb.SendEventResponse
	}
	type args struct {
		sentTime int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Response
	}{
		{
			name: "set sent time",
			fields: fields{
				SendEventResponse: &pb.SendEventResponse{},
			},
			args: args{
				sentTime: timeNow,
			},
			want: &Response{
				SendEventResponse: &pb.SendEventResponse{
					SentTime: timeNow,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Response{
				SendEventResponse: tt.fields.SendEventResponse,
			}
			if got := r.SetSentTime(tt.args.sentTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Response.SetSentTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponse_SetReason(t *testing.T) {
	type fields struct {
		SendEventResponse *pb.SendEventResponse
	}
	type args struct {
		reason string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Response
	}{
		{
			name: "set reason",
			fields: fields{
				SendEventResponse: &pb.SendEventResponse{},
			},
			args: args{
				reason: "test reason",
			},
			want: &Response{
				SendEventResponse: &pb.SendEventResponse{
					Reason: "test reason",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Response{
				SendEventResponse: tt.fields.SendEventResponse,
			}
			if got := r.SetReason(tt.args.reason); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Response.SetReason() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponse_SetDataMap(t *testing.T) {
	type fields struct {
		SendEventResponse *pb.SendEventResponse
	}
	type args struct {
		data map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Response
	}{
		{
			name: "set data map",
			fields: fields{
				SendEventResponse: &pb.SendEventResponse{},
			},
			args: args{
				data: map[string]string{"test_key": "test_value"},
			},
			want: &Response{
				SendEventResponse: &pb.SendEventResponse{
					Data: map[string]string{"test_key": "test_value"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Response{
				SendEventResponse: tt.fields.SendEventResponse,
			}
			if got := r.SetDataMap(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Response.SetDataMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponse_Write(t *testing.T) {
	res := &pb.SendEventResponse{
		Status:   pb.Status_STATUS_SUCCESS,
		Code:     pb.Code_CODE_OK,
		SentTime: time.Now().Unix(),
		Data:     map[string]string{},
	}

	errorRes := &pb.SendEventResponse{}

	successSerialization := func(m interface{}) ([]byte, error) {
		return []byte("1"), nil
	}

	failureSerialization := func(m interface{}) ([]byte, error) {
		return []byte{}, errors.New("new error")
	}
	type fields struct {
		SendEventResponse *pb.SendEventResponse
	}
	type args struct {
		s serde.SerializeFunc
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantW   string
		wantErr bool
	}{
		{
			name: "test normal write",
			fields: fields{
				SendEventResponse: res,
			},
			args: args{
				s: successSerialization,
			},
			want:    1,
			wantW:   "1",
			wantErr: false,
		},
		{
			name: "seralization error",
			fields: fields{
				SendEventResponse: errorRes,
			},
			args: args{
				s: failureSerialization,
			},
			want:    0,
			wantW:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Response{
				SendEventResponse: tt.fields.SendEventResponse,
			}
			w := &bytes.Buffer{}
			got, err := r.Write(w, tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Response.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Response.Write() = %v, want %v", got, tt.want)
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Response.Write() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
