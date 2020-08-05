package worker

import (
	"os"
	"raccoon/logger"
	"testing"
)

type void struct{}

func (v void) Write(_ []byte) (int, error) {
	return 0, nil
}
func TestMain(t *testing.M) {
	logger.Setup()
	logger.SetOutput(void{})
	os.Exit(t.Run())
}
