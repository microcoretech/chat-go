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
	"strings"

	"github.com/gofiber/fiber/v2"

	"mbobrovskyi/chat-go/internal/common/errors"
)

const (
	authTokenQueryParam = "token"
	headerAuthorization = "Authorization"
	bearerTokenType     = "Bearer"
)

type AuthMiddleware struct {
	authService AuthService
}

func (m *AuthMiddleware) Handler(ctx *fiber.Ctx) error {
	var (
		authToken string
		err       error
	)

	authToken, err = m.getTokenFromQuery(ctx)
	if err != nil {
		return err
	}

	if len(authToken) == 0 {
		authToken, err = m.getTokenFromHeader(ctx)
		if err != nil {
			return err
		}
	}

	if authToken == "" {
		return errors.NewUnauthorizedError("invalid token")
	}

	session, err := m.authService.GetSession(ctx.Context(), authToken)
	if err != nil {
		return err
	}

	if session == nil {
		return errors.NewUnauthorizedError("no session found")
	}

	ctx.Context().SetUserValue("token", authToken)
	ctx.Context().SetUserValue("session", session)

	return ctx.Next()
}

func (m *AuthMiddleware) getTokenFromQuery(ctx *fiber.Ctx) (string, error) {
	return ctx.Query(authTokenQueryParam), nil
}

func (m *AuthMiddleware) getTokenFromHeader(ctx *fiber.Ctx) (string, error) {
	authHeader := ctx.Get(headerAuthorization)
	if len(authHeader) == 0 {
		return "", errors.NewUnauthorizedError("invalid token")
	}

	tokenParts := strings.Split(authHeader, " ")

	if len(tokenParts) < 2 {
		return "", errors.NewUnauthorizedError("invalid token")
	}

	tokenPrefix := strings.ToLower(tokenParts[0])
	if !strings.EqualFold(tokenPrefix, bearerTokenType) {
		return "", errors.NewUnauthorizedError("invalid token type")
	}

	authToken := tokenParts[1]

	return authToken, nil
}

func NewAuthMiddleware(authService AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}
