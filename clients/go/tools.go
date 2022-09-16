//go:build tools

package raccoon

import (
	// static analysis tools

	_ "golang.org/x/lint/golint"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
