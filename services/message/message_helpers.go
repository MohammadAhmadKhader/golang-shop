package message

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)


// TODO: must be refactored from "UserIdKey" to more clear approach
func GetMessagesParams(r *http.Request, UserIdKey string) (uint64, uint64, int, *time.Time, error) {
	userIdKey := r.URL.Query().Get(UserIdKey)
	if userIdKey == "" {
		if userIdKey == "to" {
			return 0, 0, 0, nil, fmt.Errorf("invalid to param")
		} else {
			return 0, 0, 0, nil, fmt.Errorf("invalid from param")
		}
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

func GetMessagesParamsWithoutUserId(r *http.Request) (uint64, int, *time.Time, error) {
	lastMessageId := r.URL.Query().Get("lastMessageId")
	if lastMessageId == "" {
		return 0, 0, nil, fmt.Errorf("invalid message id")
	}
	cursor := r.URL.Query().Get("cursor")
	if cursor == "" {
		return 0, 0, nil, fmt.Errorf("invalid cursor")
	}
	limitStr := r.URL.Query().Get("limit")
	if cursor == "" {
		return 0, 0, nil, fmt.Errorf("invalid limit")
	}

	// Convert string parameters to appropriate types
	lastMessageIdInt, err := strconv.ParseUint(lastMessageId, 10, 32)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("invalid message id")
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("invalid limit")
	}

	var cursorTime time.Time
	if cursor != "" {
		cursorTime, _ = time.Parse(time.RFC3339, cursor)
	}

	return lastMessageIdInt, limit, &cursorTime, nil
}