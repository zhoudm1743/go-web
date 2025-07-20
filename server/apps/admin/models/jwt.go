package models

import (
	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims JWT声明结构体
type CustomClaims struct {
	UserID   uint     `json:"userId"`
	Username string   `json:"username"`
	RealName string   `json:"realName"`
	UUID     string   `json:"uuid"`
	RoleID   uint     `json:"roleId"`
	Roles    []string `json:"roles"`
	HomePath string   `json:"homePath,omitempty"`
	jwt.RegisteredClaims
}
