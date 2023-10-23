package responses

import "app/src/models"

type (
	AccountSignUpResponse struct {
		response
		Id       int64
		Username string
	}
	AccountMeResponse struct {
		response
		Id       int64
		Username string
		Balance  float64
	}
	AdminAccountResponseCollection struct {
		response
		Accounts []models.Account
	}
)
