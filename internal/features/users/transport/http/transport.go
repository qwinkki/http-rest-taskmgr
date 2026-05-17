package user_transport_http

type UserHTTPHandler struct {
	UsersService UsersService
}

type UsersService interface {
}

func NewUserHTTPHandler(
	usersService UsersService,
) *UserHTTPHandler {
	return &UserHTTPHandler{
		UsersService: usersService,
	}
}
