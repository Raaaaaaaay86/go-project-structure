package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/raaaaaaaay86/go-project-structure/internal/entity"
	"github.com/raaaaaaaay86/go-project-structure/internal/vo/enum/role"
	"strconv"
	"time"
)

type CustomClaim struct {
	Uid   uint          `json:"uid"`
	Roles []role.RoleId `json:"roles"`
	jwt.RegisteredClaims
}

var privateKey string = "CTGK4ZWBwb7Ichku1u2qethmeYjq9RjJ"

func CreateClaim(roles []role.RoleId, userId uint, duration time.Duration) *CustomClaim {
	now := time.Now()
	claim := CustomClaim{
		Uid:   userId,
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "video-demo",
			Subject:   strconv.Itoa(int(userId)),
			Audience:  nil,
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	return &claim
}

func Generate(user entity.User) (string, error) {
	roleIds := make([]role.RoleId, 0)
	for _, r := range user.Roles {
		roleIds = append(roleIds, r.Id)
	}
	claim := CreateClaim(roleIds, user.Id, 30*time.Minute)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString([]byte(privateKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Parse(tokenString string) (*CustomClaim, error) {
	claim := &CustomClaim{}
	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(privateKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claim, nil
}
