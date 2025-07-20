package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zhoudm1743/go-web/apps/admin/models"
	"github.com/zhoudm1743/go-web/core/facades"
)

// GetJWTSecret 获取JWT密钥
func GetJWTSecret() []byte {
	return []byte(facades.Config().Get("jwt.secret").(string))
}

// GenerateToken 生成JWT
func GenerateToken(claims models.CustomClaims) (string, error) {
	// 创建一个新的令牌对象，指定签名方法和声明
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用密钥签名并获得完整的编码令牌作为字符串
	return token.SignedString(GetJWTSecret())
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*models.CustomClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证token
	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}

// GetClaims 从Gin上下文中获取JWT声明
func GetClaims(c *gin.Context) (*models.CustomClaims, error) {
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
func GetUserID(c *gin.Context) (uint, error) {
	claims, err := GetClaims(c)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
