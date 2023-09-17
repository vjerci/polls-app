package model

type Client struct {
	RegisterDB    RegisterRepository
	AuthDB        AuthRepository
	UserDB        UserRepository
	PollListDB    PollListRepository
	PollCreateDB  PollCreateRepository
	PollDetailsDB PollDetailsRepository
	PollVoteDB    PollVoteRepository
}
