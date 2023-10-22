package responses

type AccountSignUpResponse struct {
	response
	Id       int64
	Username string
}

type AccountMeResponse struct {
	response
	Id       int64
	Username string
	Balance  float64
}
