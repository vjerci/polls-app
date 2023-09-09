package model

type Client struct {
	RegisterDB   RegisterRepository
	LoginDB      LoginRepository
	UserDB       UserRepository
	PollListDB   PollListRepository
	PollCreateDB PollCreateRepository
}
