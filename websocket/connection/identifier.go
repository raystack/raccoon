package connection

import (
	"fmt"
	"net/http"
)

type Identifer struct {
	ID   string
	Type string
}

func NewConnIdentifier(header http.Header, connIDHeader, connTypeHeader string) Identifer {
	return Identifer{
		ID: header.Get(connIDHeader),
		// If connTypeHeader is empty string. By default, it will always have an empty string as Type. This means uniqueness only depends on ID.
		Type: header.Get(connTypeHeader),
	}
}

func (i Identifer) String() string {
	return fmt.Sprintf("connection %s (%s)", i.ID, i.Type)
}
