package worker

import (
	"io"
	"os"
	"testing"

	"github.com/raystack/raccoon/logger"
)

func TestMain(t *testing.M) {
	logger.SetOutput(io.Discard)
	os.Exit(t.Run())
}
