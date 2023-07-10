package entity

import (
	"database/sql"
	"time"
)

type VideoPost struct {
	Id          uint         `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Uuid        string       `json:"uuid"`
	author      User         `json:"-"`
	AuthorId    uint         `json:"authorId"`
	CreatedAt   time.Time    `json:"createdAt,omitempty"`
	UpdatedAt   time.Time    `json:"updatedAt,omitempty"`
	DeletedAt   sql.NullTime `json:"deletedAt,omitempty"`
}

func NewVideoPost(title string, description string, uuid string, author User) *VideoPost {
	return &VideoPost{
		Title:       title,
		Description: description,
		Uuid:        uuid,
		author:      author,
		AuthorId:    author.Id,
	}
}

func (v VideoPost) Author() User {
	return v.author
}
