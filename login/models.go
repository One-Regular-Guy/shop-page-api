package login

type Request struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
type Response struct {
	Token string `json:"token"`
}
