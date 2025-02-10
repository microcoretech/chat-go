package http

import (
	"mbobrovskyi/chat-go/internal/common/http"
	"mbobrovskyi/chat-go/internal/user/domain"
)

func TokenToDto(token domain.Token) TokenDto {
	return TokenDto{
		Token:   token.Token,
		ExpIn:   token.ExpIn.Milliseconds() / 1000,
		ExpTime: token.ExpTime,
		User:    http.UserToDto(token.User),
	}
}
