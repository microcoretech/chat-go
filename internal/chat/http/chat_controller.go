// Copyright 2025 Mykhailo Bobrovskyi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"strconv"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"

	"chat-go/internal/chat/common"
	chatdomain "chat-go/internal/chat/domain"
	chatwebsocket "chat-go/internal/chat/websocket"
	"chat-go/internal/common/domain"
	"chat-go/internal/common/errors"
	"chat-go/internal/common/http"
	"chat-go/internal/infrastructure/connector"
	"chat-go/internal/infrastructure/validator"
)

type ChatController struct {
	validate       validator.Validate
	chatService    ChatService
	messageService MessageService
	connector      connector.Connector
}

func (c *ChatController) SetupRoutes(r fiber.Router) {
	r.Get("", c.getChats)
	r.Get("/ws", c.ws)
	r.Get("/:id", c.getChat)
	r.Get("/:id/messages", c.getChatMessages)
	r.Post("", c.create)
}

func (c *ChatController) getChats(ctx *fiber.Ctx) error {
	var query ChatQuery

	if err := ctx.QueryParser(&query); err != nil {
		return errors.NewBadRequestError(common.ChatDomain, err, nil)
	}

	if err := c.validate.Struct(common.ChatDomain, &query); err != nil {
		return errors.NewValidationError(common.ChatDomain, err, nil)
	}

	chatFilter, err := ChatFilterFromQuery(query)
	if err != nil {
		return err
	}

	chats, count, err := c.chatService.GetChats(ctx.Context(), &chatFilter)
	if err != nil {
		return err
	}

	return ctx.JSON(http.NewPage(
		lo.Map(chats, func(chat chatdomain.Chat, _ int) ChatDto {
			return ChatToDto(chat)
		}),
		count,
	))
}

func (c *ChatController) getChat(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return errors.NewBadRequestError(common.ChatDomain, err, map[string]any{"id": idStr})
	}

	chat, err := c.chatService.GetChat(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(ChatToDto(*chat))
}

func (c *ChatController) getChatMessages(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return errors.NewBadRequestError(common.ChatDomain, err, map[string]any{"id": idStr})
	}

	var query MessageQuery

	if err := ctx.QueryParser(&query); err != nil {
		return errors.NewBadRequestError(common.ChatDomain, err, nil)
	}

	if err := c.validate.Struct(common.ChatDomain, &query); err != nil {
		return errors.NewValidationError(common.ChatDomain, err, nil)
	}

	messageFilter, err := MessageFilterFromQuery(query)
	if err != nil {
		return err
	}

	messageFilter.ChatIDs = []uint64{id}

	messages, count, err := c.messageService.GetMessages(ctx.Context(), &messageFilter)
	if err != nil {
		return err
	}

	return ctx.JSON(http.NewPage(
		lo.Map(messages, func(message chatdomain.Message, _ int) MessageDto {
			return MessageToDto(message)
		}),
		count,
	))
}

func (c *ChatController) create(ctx *fiber.Ctx) error {
	user := ctx.Context().UserValue("user").(*domain.User)

	dto := CreateChatDto{}
	if err := ctx.BodyParser(&dto); err != nil {
		return errors.NewBadRequestError(common.ChatDomain, err, nil)
	}

	if err := c.validate.Struct(common.ChatDomain, dto); err != nil {
		return err
	}

	chat, err := ChatFromCreateDto(dto)
	if err != nil {
		return err
	}

	chat.CreatedBy = user.ID

	createdChat, err := c.chatService.CreateChat(ctx.Context(), *chat)
	if err != nil {
		return err
	}

	return ctx.JSON(ChatToDto(*createdChat))
}

func (c *ChatController) ws(ctx *fiber.Ctx) error {
	user := ctx.Context().UserValue("user").(*domain.User)
	return websocket.New(func(conn *websocket.Conn) {
		connection := chatwebsocket.NewConnection(conn.Conn, user)
		c.connector.AddConnection(connection)
		<-connection.GetCloseChan()
	})(ctx)
}

func NewChatController(
	validate validator.Validate,
	chatService ChatService,
	messageService MessageService,
	connector connector.Connector,
) *ChatController {
	return &ChatController{
		validate:       validate,
		chatService:    chatService,
		messageService: messageService,
		connector:      connector,
	}
}
