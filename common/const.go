package common

import "fooddlv/pubsub"

const (
	CurrentUser = "current_user"
)

const (
	ChannelNoteCreated pubsub.Channel = "ChannelNoteCreated"
)

type Masker interface {
	Mask(isAdmin bool)
}

type Requester interface {
	GetUserId() int
	GetEmail() string
	GetRole() string
}
