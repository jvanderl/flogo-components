package util

import (
	"fmt"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/op/go-logging"
)

// HandlePanic helper method to handle panics
func HandlePanic(name string, err *error) {
	if r := recover(); r != nil {

		log.Warningf("%s: PANIC Occurred  : %v\n", name, r)

		// todo: useful for debugging
		if log.IsEnabledFor(logging.DEBUG) {
			log.Debugf("StackTrace: %s", debug.Stack())
		}

		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}

// URLStringToFilePath convert fileURL to file path
func URLStringToFilePath(fileURL string) (string, bool) {

	if strings.HasPrefix(fileURL, "file://") {

		filePath := fileURL[7:]

		if runtime.GOOS == "windows" {
			if strings.HasPrefix(filePath, "/") {
				filePath = filePath[1:]
			}
			filePath = filepath.FromSlash(filePath)
		}

		filePath = strings.Replace(filePath, "%20", " ", -1)

		return filePath, true
	}

	return "", false
}
