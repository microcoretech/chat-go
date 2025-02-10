package domain

import (
	"golang.org/x/exp/slices"

	"mbobrovskyi/chat-go/internal/chat/errors"
)

type ChatType uint8

const (
	DirectChatType ChatType = 1
	GroupChatType  ChatType = 2
)

func (ct ChatType) Uint8() uint8 {
	return uint8(ct)
}

func (ct ChatType) Types() []ChatType {
	return []ChatType{DirectChatType, GroupChatType}
}

func (ct ChatType) IsValid() bool {
	return slices.Contains(ct.Types(), ct)
}

func NewChatType(chatType uint8) (ChatType, error) {
	if !ChatType(chatType).IsValid() {
		return 0, errors.NewInvalidChatTypeError()
	}

	return ChatType(chatType), nil
}
