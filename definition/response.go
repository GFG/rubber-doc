package definition

// Response
type Response struct {
	StatusCode  int
	Description string
	Headers     []Header
	Body        []Body
}
