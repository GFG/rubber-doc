package definition

// Request
type Request struct {
	Title       string
	Description string
	Body        []Body
	Headers     []Header
}

// IsEmpty verifies is the request is empty
func (r Request) IsEmpty() bool {
	if r.Title != "" || r.Description != "" {
		return false
	}

	if len(r.Body) > 0 || len(r.Headers) > 0 {
		return false
	}

	return true
}