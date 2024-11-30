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

type QueueUserBlocked struct {
	Type    string `json:"message_type"`
	Message UserBlocked
}
type UserBlocked struct {
	UserId    string `json:"user_id"`
	Reason    string `json:"reason"`
	Timestamp string `json:"timestamp"`
}

type QueueUserUnblocked struct {
	Type    string `json:"message_type"`
	Message UserUnblocked
}
type UserUnblocked struct {
	UserId    string `json:"user_id"`
	Timestamp string `json:"timestamp"`
}

type QueueNewRegistry struct {
	Type    string `json:"message_type"`
	Message NewRegistry
}
type NewRegistry struct {
	RegistrationId string `json:"registration_id"`
	TimeStamp      string `json:"timestamp"`
	Provider       string `json:"provider"`
}

type QueueNewUser struct {
	Type    string `json:"message_type"`
	Message NewUser
}
type NewUser struct {
	UserId         string `json:"user_id"`
	Location       string `json:"location"`
	TimeStamp      string `json:"timestamp"`
	RegistrationId string `json:"old_registration_id"`
}
