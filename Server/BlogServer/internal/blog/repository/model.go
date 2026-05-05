package repository

type ReactionTransition int

const (
	AddLike ReactionTransition = iota
	AddDislike
	LikeToDislike
	DislikeToLike
)
