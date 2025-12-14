package utils

import "time"

func TimePtrToString(format string, t *time.Time) *string {

	if t == nil {
		return nil
	}

	s := t.Format(format)

	return &s
}

func TimeToString(format string, t time.Time) string {
	return t.Format(format)
}
