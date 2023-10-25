package responses

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
	AdminAccountResponse struct {
		response
		Id       int64
		Username string
		Password string
		Type     string
		Balance  float64
	}
)
