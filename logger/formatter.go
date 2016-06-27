package logger

import (
	"fmt"
	"time"
)

type Formatter struct {
}

func prependDate(message string) string {
	t := time.Time{}

	return fmt.Sprintf("%s %s", t.Format(time.UnixDate), message)
}

func prepareMessage(message string, context interface{}) string {
	return fmt.Sprintf("%s #%v", message, context)
}

func format(message string, level string) string {
	return prependDate(fmt.Sprintf("%s %s", level, message))
}

