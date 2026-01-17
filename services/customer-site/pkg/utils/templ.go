package utils

import (
	"crypto/rand"
	"fmt"

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
