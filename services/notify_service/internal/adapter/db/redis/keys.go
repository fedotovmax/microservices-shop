package redisadapter

import "fmt"

func userChatKey(id string) string {
	return fmt.Sprintf("USER_CHAT_KEY(%s)", id)
}

func chatUserKey(id int64) string {
	return fmt.Sprintf("CHAT_USER_KEY(%d)", id)
}

func eventKey(eventID string) string {
	return fmt.Sprintf("EVENT_KEY(%s)", eventID)
}
