package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWT相关常量
const (
	TokenExpireDuration = 24 * time.Hour
	TokenIssuer         = "go-web-admin"
	SecretKey           = "your-secret-key-here" // 实际项目中应从配置文件读取
)

// CustomClaims 自定义JWT声明
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, username, role string) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    TokenIssuer,
		},
	}

	// 使用HS256算法创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey))
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// 校验token
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}

// RefreshToken 刷新JWT令牌
func RefreshToken(tokenString string) (string, error) {
	// 解析现有token
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	// 检查是否需要刷新
	now := time.Now().Unix()
	// 如果有效期剩余超过12小时，则不刷新
	if claims.ExpiresAt-now > int64(12*time.Hour.Seconds()) {
		return tokenString, nil
	}

	// 生成新token
	return GenerateToken(claims.UserID, claims.Username, claims.Role)
}
