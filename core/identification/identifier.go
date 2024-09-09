package identification

import (
	"fmt"
)

type Identifier struct {
	ID    string
	Group string
}

func (i Identifier) String() string {
	return fmt.Sprintf("[%s] %s", i.Group, i.ID)
}
