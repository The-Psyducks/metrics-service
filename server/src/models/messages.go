package models

type QueueMessage struct {
	Type    string      `json:"message_type"`
	Message interface{} `json:"message"`
}

type QueueLoginAttempt struct {
	Type    string       `json:"message_type"`
	Message LoginAttempt `json:"message"`
}

type LoginAttempt struct {
	WasSuccessful bool   `json:"was_successful"`
	Timestamp     string `json:"timestamp"`
	Provider      string `json:"provider"`
	UserId        string `json:"user_id"`
}
