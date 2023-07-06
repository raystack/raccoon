package worker

import (
	"os"
	"testing"

	"github.com/raystack/raccoon/logger"
)

type void struct{}

func (v void) Write(_ []byte) (int, error) {
	return 0, nil
}
func TestMain(t *testing.M) {
	logger.SetOutput(void{})
	os.Exit(t.Run())
}
