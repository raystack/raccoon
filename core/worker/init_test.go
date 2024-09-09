package worker

import (
	"io"
	"os"
	"testing"

	"github.com/raystack/raccoon/pkg/logger"
)

func TestMain(t *testing.M) {
	logger.SetOutput(io.Discard)
	os.Exit(t.Run())
}
