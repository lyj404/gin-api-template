package tokenutil

import (
	"fmt"
	"gin-api-template/config"
	"gin-api-template/domain"
	"gin-api-template/domain/entity"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// CreateAccessToken创建一个访问令牌
func CreateAccessToken(user *entity.User, secret string, expire int) (accessToken string, err error) {
	// 创建自定义声明，其中包含用户信息和过期时间
	claims := &domain.JwtCustomClaims{
		ID:   strconv.FormatUint(uint64(user.ID), 16), //用户ID，转换为十六进制字符
		Name: user.Name,                               // 用户名
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expire))),
		},
	}

	// 使用声明创建新的jwt令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用密钥对令牌进行签名
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	// 添加前缀，以方便后续处理
	tokenWithPrefix := config.CfgToken.TokenPrefix + t
	return tokenWithPrefix, nil
}

// CreateRefreshToken创建一个刷新令牌
func CreateRefreshToken(user *entity.User, secret string, expire int) (refreshToken string, err error) {
	// 创建自定义声明，其中包含用户信息和过期时间
	claimsRefresh := &domain.JwtCustomRefreshClaims{
		ID: strconv.FormatUint(uint64(user.ID), 16), //用户ID，转换为十六进制字符
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expire))),
		},
	}

	// 使用声明创建新的jwt令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	// 使用密钥对令牌进行签名
	rt, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	// 添加前缀，以方便后续处理
	tokenWithPrefix := config.CfgToken.TokenPrefix + rt
	return tokenWithPrefix, nil
}

// IsAuthorized验证令牌是否有效
func IsAuthorized(requestToken string, secret string) (bool, error) {
	// 解析令牌并验证签名方法
	_, err := parseToken(requestToken, secret)
	if err != nil {
		// 令牌解析错误
		return false, err
	}
	return true, nil
}

// ExtractIDFromToken从token提取用户ID
func ExtractIDFromToken(requestToken string, secret string) (string, error) {
	// 解析令牌并验证签名方法
	token, err := parseToken(requestToken, secret)
	if err != nil {
		// 令牌解析错误
		return "", err
	}

	// 断言并获取声明
	claims, ik := token.Claims.(jwt.MapClaims)
	if !ik || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	id, ok := claims["id"].(string)
	if !ok {
		return "", fmt.Errorf("id not found in token")
	}
	// 返回用户ID
	return id, nil
}

// parseToken 解析令牌并验证签名方法
func parseToken(requestToken string, secret string) (*jwt.Token, error) {
	return jwt.Parse(requestToken, func(t *jwt.Token) (interface{}, error) {
		// SigningMethodHS256 是 SigningMethodHMAC 的一个具体实现，因此可以使用SigningMethodHMAC来断言
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
}
