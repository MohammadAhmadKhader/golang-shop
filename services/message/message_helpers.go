package message

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// user UserIdKey set as "to" if you wish to retrieve chat messages
//
// to route which you get a specific user params you can use UserIdKey as "from" 
// because its meant to get all the messages a specific user has sent
func GetMessagesParams(r *http.Request, UserIdKey string) (uint64, uint64, int, *time.Time, error) {
	userIdKey := r.URL.Query().Get(UserIdKey)
	if userIdKey == "" {
		return 0, 0, 0, nil, fmt.Errorf("invalid to param")
	}
	lastMessageId := r.URL.Query().Get("lastMessageId")
	if lastMessageId == "" {
		return 0, 0, 0, nil, fmt.Errorf("invalid message id")
	}
	cursor := r.URL.Query().Get("cursor")
	if cursor == "" {
		return 0, 0, 0, nil, fmt.Errorf("invalid cursor")
	}
	limitStr := r.URL.Query().Get("limit")
	if cursor == "" {
		return 0, 0, 0, nil, fmt.Errorf("invalid limit")
	}

	// Convert string parameters to appropriate types
	userIdKeyAsUInt, err := strconv.ParseUint(userIdKey, 10, 32)
	if err != nil {
		return 0, 0, 0, nil, fmt.Errorf("invalid to param")
	}
	lastMessageIdInt, err := strconv.ParseUint(lastMessageId, 10, 32)
	if err != nil {
		return 0, 0, 0, nil, fmt.Errorf("invalid message id")
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, 0, 0, nil, fmt.Errorf("invalid limit")
	}

	var cursorTime time.Time
	if cursor != "" {
		cursorTime, _ = time.Parse(time.RFC3339, cursor)
	}

	return userIdKeyAsUInt, lastMessageIdInt, limit, &cursorTime, nil
}