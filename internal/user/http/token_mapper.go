package http

import (
	"chat/internal/common/http"
	"chat/internal/user/domain"
)

func TokenToDto(token domain.Token) TokenDto {
	return TokenDto{
		Token:   token.Token,
		ExpIn:   token.ExpIn.Milliseconds() / 1000,
		ExpTime: token.ExpTime,
		User:    http.UserToDto(token.User),
	}
}
