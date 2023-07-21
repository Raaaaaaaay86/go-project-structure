package entity

import (
	"database/sql"
	"github.com/raaaaaaaay86/go-project-structure/internal/vo/enum/role"
	"time"
)

type Role struct {
	Id        role.RoleId  `json:"id"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at,omitempty"`
	UpdatedAt time.Time    `json:"updated_at,omitempty"`
	DeletedAt sql.NullTime `json:"deleted_at,omitempty"`
}

func NewRole(role role.RoleId) *Role {
	return &Role{
		Id:   role,
		Name: role.Code(),
	}
}
