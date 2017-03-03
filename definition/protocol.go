package definition

import (
	"errors"
	"strings"
)

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

// NewProtocolFromURL Returns the protocol based on the URL
func NewProtocolFromURL(URL string) (proto Protocol, err error) {
	s := strings.Split(URL, ":")

	proto = Protocol(s[0])

	if !proto.isValid() {
		err = errors.New("Invalid Protocol")
		return
	}
	return
}
