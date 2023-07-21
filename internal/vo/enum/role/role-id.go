package role

type RoleId uint

const (
	_ RoleId = iota
	SuperAdmin
	Admin
	User
)

func (r RoleId) Code() string {
	switch r {
	case SuperAdmin:
		return "SUPER_ADMIN"
	case Admin:
		return "ADMIN"
	case User:
		return "USER"
	}
	return "USER"
}
