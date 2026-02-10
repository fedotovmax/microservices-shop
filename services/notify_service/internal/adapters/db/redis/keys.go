package redis

import "fmt"

func UserChatKey(id string) string {
	return fmt.Sprintf("USER_CHAT_KEY(%s)", id)
}

func ChatUserKey(id int64) string {
	return fmt.Sprintf("CHAT_USER_KEY(%d)", id)
}

func EventKey(eventID string) string {
	return fmt.Sprintf("EVENT_KEY(%s)", eventID)
}
