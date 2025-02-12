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

type SignUpRequest struct {
	Email     string `json:"email"  validate:"required,email,min=3,max=255"`
	Username  string `json:"username" validate:"required,username,min=3,max=50"`
	FirstName string `json:"firstName" validate:"required,name,min=1,max=50"`
	LastName  string `json:"lastName" validate:"required,name,min=1,max=50"`
	Password  string `json:"password" validate:"required,password,min=8,max=255"`
}
