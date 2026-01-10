package utils

import (
	"github.com/Oudwins/tailwind-merge-go/pkg/twmerge"
)

func TwMerge(classes ...string) string {
	return twmerge.Merge(classes...)
}

func If[T comparable](condition bool, value T) T {
	var empty T
	if condition {
		return value
	}
	return empty
}

// func Cn(classes ...interface{}) string {
// 	a := clsx(classes...)
// 	r := twmerge.Merge(a)
// 	return r
// }

// func clsx(classes ...interface{}) string {
// 	var classList []string

// 	for _, class := range classes {
// 		switch v := class.(type) {
// 		case string:
// 			if v != "" {
// 				classList = append(classList, v)
// 			}
// 		case []string:
// 			if len(v) > 0 {
// 				classList = append(classList, v...)
// 			}
// 		case map[string]bool:
// 			for key, value := range v {
// 				if value {
// 					classList = append(classList, key)
// 				}
// 			}
// 		case bool:
// 		default:
// 		}
// 	}

// 	return strings.Join(classList, " ")
// }
