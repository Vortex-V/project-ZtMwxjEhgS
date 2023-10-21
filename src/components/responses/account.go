package responses

type AccountSignUpResponse struct {
	response
	Id       int
	Username string
}

type AccountMeResponse struct {
	response
	Id       int
	Username string
	Balance  float64
}
