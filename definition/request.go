package definition

// Request
type Request struct {
	Title       string
	Description string
	Method      string
	Body        []Body
	Headers     []Header
	ContentType string
}
