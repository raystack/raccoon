package websocket

import "net/http"

type ConnIdentifier struct {
	cID   string
	cType string
}

func NewConnIdentifier(header http.Header, connIDHeader, connTypeHeader string) ConnIdentifier {
	return ConnIdentifier{
		cID:   header.Get(connIDHeader),
		// If connTypeHeader is empty string. By default, it will always have an empty string as cType. This means uniqueness only depends on cID.
		cType: header.Get(connTypeHeader),
	}
}
