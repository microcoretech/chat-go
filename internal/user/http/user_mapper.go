package http

import (
	"mbobrovskyi/chat-go/internal/common/domain"
)

func UserFromSignUpRequest(req SignUpRequest) domain.User {
	return domain.User{
		Email:     req.Email,
		Username:  req.Username,
		Role:      domain.UserRole,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
}
