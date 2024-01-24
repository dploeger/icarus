package adapters_test

import (
	"github.com/dploeger/icarus/v2/internal"
	"os"
	"testing"
)
import "github.com/rogpeppe/go-internal/testscript"

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"icarus": internal.Main,
	}))
}

func TestAdapters(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/script",
	})
}
