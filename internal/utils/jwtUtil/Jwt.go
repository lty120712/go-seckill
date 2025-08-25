package jwtUtil

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"go-chat/configs"
	"time"
)

type Claims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

// GenerateJWT 生成 JWT Token
func GenerateJWT(id uint) (string, error) {
	// 获取动态的 JWT 配置信息
	jwtConfig := configs.AppConfig.Jwt
	// 解析过期时间（例如 "24h"）
	expirationTime, err := time.ParseDuration(jwtConfig.ExpirationTime) // 解析字符串为 time.Duration
	if err != nil {
		return "", fmt.Errorf("invalid expiration time format: %v", err)
	}
	// 计算 ExpiresAt 时间
	expiresAt := time.Now().Add(expirationTime)
	// 将 Audience 转换为 []string 类型（如果它是单个字符串）
	audience := []string{jwtConfig.Audience}
	// 定义 JWT 的 Claims（声明部分）
	claims := Claims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt), // 设置正确的过期时间
			Issuer:    jwtConfig.Issuer,              // 从配置读取 Issuer
			Audience:  audience,                      // 转换为 []string 类型
		},
	}
	// 创建 JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用配置中的密钥签名 Token
	return token.SignedString([]byte(jwtConfig.SecretKey))
}

// 解析 JWT Token 并验证
func ParseJWT(tokenStr string) (*Claims, error) {
	// 获取动态的 JWT 配置信息
	jwtConfig := configs.AppConfig.Jwt

	// 解析并校验 token
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtConfig.SecretKey), nil // 从配置中获取密钥
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	// 校验 token 是否有效
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token or expired")
}
