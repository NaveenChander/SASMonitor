package lib

import (
	"fmt"
	"time"
)

// Log ... Log message to Console
func Log(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	s = fmt.Sprintf("%v: %s", time.Now().Format("2006-01-02T15:04:05.000"), s)
	fmt.Println(s)
}
