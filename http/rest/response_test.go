package rest

import (
	"bytes"
	"errors"
	pb "raccoon/pkg/proto"
	"raccoon/pkg/serialization"
	"reflect"
	"testing"
	"time"
)

func TestResponse_SetCode(t *testing.T) {
	type fields struct {
		EventResponse *pb.EventResponse
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
				EventResponse: &pb.EventResponse{},
			},
			args: args{
				code: pb.Code_UNKNOWN_CODE,
			},
			want: &Response{
				EventResponse: &pb.EventResponse{
					Code: pb.Code_UNKNOWN_CODE,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Response{
				EventResponse: tt.fields.EventResponse,
			}
			if got := r.SetCode(tt.args.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Response.SetCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponse_SetStatus(t *testing.T) {
	type fields struct {
		EventResponse *pb.EventResponse
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
				EventResponse: &pb.EventResponse{},
			},
			args: args{
				status: pb.Status_SUCCESS,
			},
			want: &Response{
				EventResponse: &pb.EventResponse{
					Status: pb.Status_SUCCESS,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Response{
				EventResponse: tt.fields.EventResponse,
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
		EventResponse *pb.EventResponse
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
				EventResponse: &pb.EventResponse{},
			},
			args: args{
				sentTime: timeNow,
			},
			want: &Response{
				EventResponse: &pb.EventResponse{
					SentTime: timeNow,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Response{
				EventResponse: tt.fields.EventResponse,
			}
			if got := r.SetSentTime(tt.args.sentTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Response.SetSentTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponse_SetReason(t *testing.T) {
	type fields struct {
		EventResponse *pb.EventResponse
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
				EventResponse: &pb.EventResponse{},
			},
			args: args{
				reason: "test reason",
			},
			want: &Response{
				EventResponse: &pb.EventResponse{
					Reason: "test reason",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Response{
				EventResponse: tt.fields.EventResponse,
			}
			if got := r.SetReason(tt.args.reason); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Response.SetReason() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponse_SetDataMap(t *testing.T) {
	type fields struct {
		EventResponse *pb.EventResponse
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
				EventResponse: &pb.EventResponse{},
			},
			args: args{
				data: map[string]string{"test_key": "test_value"},
			},
			want: &Response{
				EventResponse: &pb.EventResponse{
					Data: map[string]string{"test_key": "test_value"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Response{
				EventResponse: tt.fields.EventResponse,
			}
			if got := r.SetDataMap(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Response.SetDataMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponse_Write(t *testing.T) {
	s := &serialization.MockSerializer{}
	res := &pb.EventResponse{
		Status:   pb.Status_SUCCESS,
		Code:     pb.Code_OK,
		SentTime: time.Now().Unix(),
		Data:     map[string]string{},
	}
	s.On("Serialize", &Response{res}).Return("1", nil)

	errorRes := &pb.EventResponse{}
	s.On("Serialize", &Response{errorRes}).Return("", errors.New("serialization failure"))
	type fields struct {
		EventResponse *pb.EventResponse
	}
	type args struct {
		s serialization.Serializer
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
				EventResponse: res,
			},
			args: args{
				s: s,
			},
			want:    4,
			wantW:   "[49]",
			wantErr: false,
		},
		{
			name: "seralization error",
			fields: fields{
				EventResponse: errorRes,
			},
			args: args{
				s: s,
			},
			want:    0,
			wantW:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Response{
				EventResponse: tt.fields.EventResponse,
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
