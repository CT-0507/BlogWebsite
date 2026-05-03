package domain

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/model"
	"github.com/google/uuid"
)

type Comment struct {
	ID               uuid.UUID  `json:"commentId"`
	BlogID           int64      `json:"blogId"`
	ActorType        string     `json:"actorType"`
	ActorID          *string    `json:"actorId"`
	ActorDisplayName string     `json:"actorDisplayName"`
	ActorAvatarURL   *string    `json:"actorAvatarUrl"`
	Content          string     `json:"content"`
	UserReaction     *string    `json:"userReaction"`
	Status           string     `json:"status"` // active | deleted | hidden
	ParentCommentID  *uuid.UUID `json:"parentCommentId"`
	RootCommentID    uuid.UUID  `json:"rootCommentId"`
	ReplyCount       int64      `json:"replyCount"`
	LikeCount        int32      `json:"likeCount"`
	DislikeCount     int32      `json:"dislikeCount"`
	Depth            int16      `json:"-"`
	model.Audit
}

type CreateCommentModel struct {
	BlogID           int64
	ActorType        string
	ActorID          string
	ActorDisplayName string
	ActorAvatarURL   *string
	Content          string
	ParentCommentID  *uuid.UUID
	RootCommentID    *uuid.UUID // Empty on the comment is the root itself
	Depth            int16
}
