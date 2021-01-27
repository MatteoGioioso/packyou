package errorPkg

import (
	"fmt"
	"runtime"
)

func New(err error, message string) error {
	_, file, line, _ := runtime.Caller(1)
	msg := fmt.Sprintf("%s %s:%d", message, file, line )
	return fmt.Errorf("%v\n %v", err, msg)
}
