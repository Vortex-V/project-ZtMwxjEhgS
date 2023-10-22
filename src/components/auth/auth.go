package auth

import (
	"app/src/models"
	"errors"
	"github.com/beego/beego/v2/server/web"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrorUsernameOrPasswordIncorrect = errors.New("username or password is incorrect")
)

func Login(account *models.Account, password string) (string, error) {
	err := models.Get(account, "Username")
	if err != nil {
		return "", ErrorUsernameOrPasswordIncorrect
	}

	if err = CheckPasswordHash(password, account.Password); err != nil {
		return "", ErrorUsernameOrPasswordIncorrect
	}

	token, err := createAccessToken(account.Id)
	if err != nil {
		return "", err
	}
	account.IsNeedRelogin = false
	_, err = models.Update(account, "IsNeedRelogin")
	if err != nil {
		return "", err
	}
	return token, nil
}

func createAccessToken(id int) (string, error) {
	claims := jwt.MapClaims{}
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey, _ := web.AppConfig.String("jwt")
	jwtToken, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 16)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
