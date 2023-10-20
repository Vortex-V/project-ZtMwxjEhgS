package responses

type account struct {
	response
	id       int64
	username string
	password string
}
