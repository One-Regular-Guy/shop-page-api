package register

type Request struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required,alphanum,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=40"`
}
type Response struct {
	Message string `json:"message"`
}
