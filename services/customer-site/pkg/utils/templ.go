package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Oudwins/tailwind-merge-go/pkg/twmerge"
	"github.com/a-h/templ"
)

func TwMerge(classes ...string) string {
	return twmerge.Merge(classes...)
}

func IfElse[T any](condition bool, trueValue T, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

func If[T comparable](condition bool, value T) T {
	var empty T
	if condition {
		return value
	}
	return empty
}

func RandomID() string {
	return fmt.Sprintf("id-%s", rand.Text())
}

func MergeAttributes(attrs ...templ.Attributes) templ.Attributes {
	merged := templ.Attributes{}
	for _, attr := range attrs {
		for k, v := range attr {
			merged[k] = v
		}
	}
	return merged
}

func CreateJSObject(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}

	s := string(b)
	slog.Info("OBJ", slog.String("obj", s))
	return s
}

func DatastarSseWithOptions(method, url, opts string) string {

	return fmt.Sprintf("@%s('%s', %s)", method, url, opts)

}
