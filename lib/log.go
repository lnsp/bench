package lib

import (
	"os"

	"github.com/op/go-logging"
)

var (
	log       = logging.MustGetLogger("lib")
	logFormat = logging.MustStringFormatter(
		`%{color}%{shortfunc:9.9s} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
)

// SetVerbose sets the logging level of the library.
func SetVerbose(verbose bool) {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, logFormat)
	backendLeveler := logging.AddModuleLevel(backendFormatter)
	if verbose {
		backendLeveler.SetLevel(logging.INFO, "")
	} else {
		backendLeveler.SetLevel(logging.ERROR, "")
	}
	logging.SetBackend(backendLeveler)
}
