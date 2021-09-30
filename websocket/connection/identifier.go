package connection

import (
	"fmt"
	"net/http"
)

type Identifer struct {
	ID   string
	Group string
}

func NewConnIdentifier(header http.Header, connIDHeader, connGroupHeader string) Identifer {
	return Identifer{
		ID: header.Get(connIDHeader),
		// If connGroupHeader is empty string. By default, it will always have an empty string as Group. This means uniqueness only depends on ID.
		Group: header.Get(connGroupHeader),
	}
}

func (i Identifer) String() string {
	return fmt.Sprintf("connection %s (%s)", i.ID, i.Group)
}
