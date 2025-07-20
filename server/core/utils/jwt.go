package utils

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zhoudm1743/go-web/core/facades"
)

// JWT自定义声明结构
type CustomClaims struct {
	UserID   int    `json:"userId"`
	Username string `json:"username"`
	RoleID   int    `json:"roleId"`
	jwt.RegisteredClaims
}

// GetJWTSecret 获取JWT密钥
func GetJWTSecret() []byte {
	secret := facades.Config().Get("jwt.secret")
	if secret == nil {
		return []byte("default_secret_key")
	}
	return []byte(secret.(string))
}

// GenerateToken 生成JWT
func GenerateToken(userId int, username string, roleId int) (string, error) {
	// 从配置获取过期时间，默认7天
	expiresDays := 7
	if value := facades.Config().Get("jwt.token_expire_days"); value != nil {
		if v, ok := value.(int); ok {
			expiresDays = v
		}
	}

	// 创建声明
	claims := CustomClaims{
		UserID:   userId,
		Username: username,
		RoleID:   roleId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * time.Duration(expiresDays))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// 创建一个新的令牌对象，指定签名方法和声明
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用密钥签名并获得完整的编码令牌作为字符串
	return token.SignedString(GetJWTSecret())
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证token
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}

// GetClaims 从Gin上下文中获取JWT声明
func GetClaims(c *gin.Context) (*CustomClaims, error) {
	// 从Header中获取token
	token := c.GetHeader("Authorization")
	if token == "" {
		return nil, errors.New("未找到令牌")
	}

	// 如果token前缀是Bearer，则去掉前缀
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// 解析token
	claims, err := ParseToken(token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

// GetUserID 获取用户ID
func GetUserID(c *gin.Context) (int, error) {
	claims, err := GetClaims(c)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
