package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type VideoComment struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	video     VideoPost          `json:"video" bson:"video"`
	VideoId   uint               `json:"videoId" bson:"video_id"`
	author    User               `json:"author" bson:"author"`
	AuthorId  uint               `json:"author_id" bson:"author_id"`
	Content   string             `json:"comment" bson:"content"`
	CreatedAt time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updated_at"`
}

func NewVideoComment(video VideoPost, author User, content string, createdAt time.Time, updatedAt time.Time) *VideoComment {
	return &VideoComment{
		video:     video,
		VideoId:   video.Id,
		author:    author,
		AuthorId:  author.Id,
		Content:   content,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func (v *VideoComment) Video() VideoPost {
	return v.video
}

func (v *VideoComment) Author() User {
	return v.author
}
