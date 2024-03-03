package http

type SignUpRequest struct {
	Email     string `json:"email"  validate:"required,email,min=3,max=255"`
	Username  string `json:"username" validate:"required,username,min=3,max=50"`
	FirstName string `json:"firstName" validate:"required,name,min=1,max=50"`
	LastName  string `json:"lastName" validate:"required,name,min=1,max=50"`
	Password  string `json:"password" validate:"required,password,min=8,max=255"`
}
