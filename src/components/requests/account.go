package requests

type accountRequest struct {
	request
	Username string `valid:"Required;"`
	Password string `valid:"Required;"`
}

type SignUpRequest struct {
	accountRequest
}
