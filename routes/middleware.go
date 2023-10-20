package routers

import (
	"fmt"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

func authFilter(ctx *context.Context) {
	authHeader := ctx.Input.Header("Authorization")
	if authHeader == "" {
		ctx.Output.SetStatus(401)
		ctx.Output.Body([]byte("Authorization header missing"))
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		jwtKey, _ := web.AppConfig.String("jwt")
		return []byte(jwtKey), nil
	})

	if err != nil || !token.Valid {
		ctx.Output.SetStatus(401)
		ctx.Output.Body([]byte("Invalid or expired token"))
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		ctx.Input.SetData("accountId", (int)(claims["id"].(float64)))
	}
}
