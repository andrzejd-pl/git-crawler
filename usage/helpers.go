package usage

import (
	"fmt"
	"io"
)

func CheckErrorWithPanic(logger io.Writer, err error) {
	CheckError(logger, err, true)
}

func CheckErrorWithOnlyLogging(logger io.Writer, err error) {
	CheckError(logger, err, false)
}

func CheckError(logger io.Writer, err error, _panic bool) {
	if err != nil {
		if _panic {
			panic("error: " + err.Error())
		}

		_, err := fmt.Fprintf(logger, "error: %s\n", err.Error())
		CheckError(logger, err, false)
	}
}
