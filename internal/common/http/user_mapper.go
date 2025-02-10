package http

import (
	"mbobrovskyi/chat-go/internal/common/domain"
)

func UserToDto(user domain.User) UserDto {
	return UserDto{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Role:      user.Role,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		AboutMe:   user.AboutMe,
		Image:     user.Image,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
