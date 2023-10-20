package responses

type Response interface {
	implement()
}

type response struct {
	Response
}
