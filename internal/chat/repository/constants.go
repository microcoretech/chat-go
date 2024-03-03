package repository

const (
	userTableName     = "users"
	chatTableName     = "chats"
	userChatTableName = "user_chats"
	messageTableName  = "messages"
)

const (
	messageFields = `m.id, m.text, m.status, m.chat_id, m.created_by, m.created_at, m.updated_at`
)
