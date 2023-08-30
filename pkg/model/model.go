package model

type Client struct {
	RegisterDB RegisterRepository
	LoginDB    LoginRepository
	UserDB     UserRepository
}
