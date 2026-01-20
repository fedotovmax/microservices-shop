package redisadapter

import "fmt"

func appKey(secret string) string {
	return fmt.Sprintf("APP_%s", secret)
}
