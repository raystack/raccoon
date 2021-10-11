package connection

import (
	"fmt"
)

type Identifier struct {
	ID    string
	Group string
}

func (i Identifier) String() string {
	return fmt.Sprintf("connection [%s] %s", i.Group, i.ID)
}
