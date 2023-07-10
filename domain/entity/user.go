package entity

import (
	"database/sql"
	"github.com/raaaaaaaay86/go-project-structure/domain/vo"
	"time"
)

type User struct {
	Id        uint                 `json:"id"`
	Roles     []Role               `json:"roles" gorm:"many2many:user_roles"`
	Username  string               `json:"username"`
	Password  vo.EncryptedPassword `json:"password"`
	Email     vo.Email             `json:"email"`
	CreatedAt time.Time            `json:"created_at,omitempty"`
	UpdatedAt time.Time            `json:"updated_at,omitempty"`
	DeletedAt sql.NullTime         `json:"deleted_at,omitempty"`
}

func NewUser(username string, password vo.EncryptedPassword, email vo.Email, roles ...Role) *User {
	return &User{
		Roles:    roles,
		Username: username,
		Password: password,
		Email:    email,
	}
}
