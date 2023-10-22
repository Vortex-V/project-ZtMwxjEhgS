package routers

import (
	"app/src/models"
	"fmt"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
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
		id := (int64)(claims["id"].(float64))
		query := models.Find(new(models.Account), "is_need_relogin").Where("id = ?")
		var isNeedRelogin bool
		models.Raw(query, id).QueryRow(&isNeedRelogin)
		if isNeedRelogin {
			ctx.Output.SetStatus(401)
			ctx.Output.Body([]byte("Invalid or expired token"))
			return
		}
		ctx.Input.SetParam("accountId", strconv.FormatInt(id, 10))
	}
}
