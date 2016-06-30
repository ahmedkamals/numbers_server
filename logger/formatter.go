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

func prepareMessage(message string, context ...Context) string {
	return fmt.Sprintf("%s #%v", message, context)
}

func format(message string, level Level) string {
	return prependDate(fmt.Sprintf("%s %s", levelAsString(level), message))
}

