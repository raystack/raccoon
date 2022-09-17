package raccoon

type (
	Event struct {
		Type string
		Data interface{}
	}

	Response struct {
		Status   int32
		Code     int32
		SentTime int64
		Reason   string
		Data     map[string]string
	}
)

type Client interface {
	// Send makes a request to provided events.
	Send([]*Event) (string, *Response, error)
}
