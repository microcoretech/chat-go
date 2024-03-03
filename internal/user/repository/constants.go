package repository

const (
	usersTableName           = "users"
	userCredentialsTableName = "user_credentials"
)

const (
	userFields            = `u.id, u.email, u.username, u.role, u.first_name, u.last_name, u.about_me, u.image_url, u.created_at, u.updated_at`
	userCredentialsFields = `uc.user_id, uc.password`
)
