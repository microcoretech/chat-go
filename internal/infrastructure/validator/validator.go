// Copyright 2025 MicroCore Tech
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

package validator

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"

	"chat-go/internal/common/errors"
)

type Validate interface {
	Struct(domain string, s any) error
	Var(domain string, variable any, tags string) error
}

type validate struct {
	validate *validator.Validate
}

func (v *validate) Struct(domain string, s any) error {
	if err := v.validate.Struct(s); err != nil {
		validatorErrors := err.(validator.ValidationErrors)
		errData := make(map[string]any)
		errorList := make([]map[string]any, 0)

		for _, validatorError := range validatorErrors {
			data := make(map[string]any)
			data["error"] = validatorError.Error()
			data["tag"] = validatorError.Tag()
			data["actualTag"] = validatorError.ActualTag()
			data["namespace"] = validatorError.Namespace()
			data["structNamespace"] = validatorError.StructNamespace()
			data["field"] = validatorError.Field()
			data["structField"] = validatorError.StructField()
			data["value"] = validatorError.Value()
			data["param"] = validatorError.Param()
			data["kind"] = validatorError.Kind()
			errorList = append(errorList, data)
		}

		errData["errors"] = errorList

		return errors.NewValidationError(domain, err, errData)
	}

	return nil
}

func (v *validate) Var(domain string, variable any, tags string) error {
	if err := v.validate.Var(variable, tags); err != nil {
		return errors.NewValidationError(domain, err, nil)
	}
	return nil
}

func nameValidator(fl validator.FieldLevel) bool {
	alphaWithSpacesRegexp := regexp.MustCompile("^[a-zA-Z ,.'-]+$")
	return alphaWithSpacesRegexp.MatchString(fl.Field().String())
}

func usernameValidator(fl validator.FieldLevel) bool {
	usernameRegexp := regexp.MustCompile("^[a-zA-Z0-9._-]+$")
	return usernameRegexp.MatchString(fl.Field().String())
}

func passwordValidator(fl validator.FieldLevel) bool {
	const (
		minPasswordLength = 8
		maxPasswordLength = 255
	)

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	passwordLen := len(fl.Field().String())
	if passwordLen < minPasswordLength || passwordLen > maxPasswordLength {
		return false
	}

	for _, char := range fl.Field().String() {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasNumber && hasSpecial
}

func New() (Validate, error) {
	v := validator.New()

	if err := v.RegisterValidation("name", nameValidator); err != nil {
		return nil, err
	}

	if err := v.RegisterValidation("username", usernameValidator); err != nil {
		return nil, err
	}

	if err := v.RegisterValidation("password", passwordValidator); err != nil {
		return nil, err
	}

	return &validate{validate: v}, nil
}
