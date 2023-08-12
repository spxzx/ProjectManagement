package jwts

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spxzx/project-user/config"
	"strings"
	"time"
)

type JwtToken struct {
	AccessToken  string
	AccessExp    int64
	RefreshToken string
	RefreshExp   int64
}

func CreateToken(val string, conf *config.TokenConf, ip string) *JwtToken {
	accessExp := time.Now().Add(conf.AccessExp).Unix()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token": val,
		"exp":   accessExp,
		"ip":    ip,
	})
	aToken, _ := accessToken.SignedString([]byte(conf.AccessSecret))

	refreshExp := time.Now().Add(conf.RefreshExp).Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token": val,
		"exp":   refreshExp,
		"ip":    ip,
	})
	rToken, _ := refreshToken.SignedString([]byte(conf.RefreshSecret))

	return &JwtToken{
		AccessToken:  aToken,
		AccessExp:    accessExp,
		RefreshToken: rToken,
		RefreshExp:   refreshExp,
	}
}

func ParseToken(tokenString string, secret string, ip string) (string, error) {
	if strings.Contains(tokenString, "bearer") {
		tokenString = strings.ReplaceAll(tokenString, "bearer ", "")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header[""])
		}
		return []byte(secret), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		token := claims["token"].(string)
		exp := int64(claims["exp"].(float64))
		if exp <= time.Now().Unix() {
			return "", errors.New("token已过期")
		}
		if claims["ip"] != ip {
			return "", errors.New("ip不合法")
		}
		return token, nil
	} else {
		return "", err
	}
}
