package definition

import "strings"

const (
	HTTP_PROTOCOL  = "http"
	HTTPS_PROTOCOL = "https"
)

// Protocol
type Protocol string

// Checks if the used protocol is excepted
func (p Protocol) isValid() bool {
	lower := strings.ToLower(string(p))
	return lower == HTTP_PROTOCOL || lower == HTTPS_PROTOCOL
}
