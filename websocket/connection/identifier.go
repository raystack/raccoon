package connection

import (
	"fmt"
)

type Identifer struct {
	ID    string
	Group string
}

func (i Identifer) String() string {
	return fmt.Sprintf("connection [%s] %s", i.Group, i.ID)
}
