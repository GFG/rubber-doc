package definition

import (
	"strings"
	"errors"
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

func NewProtocolFromURL(URL string) (Protocol, error) {
	s := strings.Split(URL, ":")
	prot := Protocol(s[0])

	if !prot.isValid() {
		return nil, errors.New("Invalid Protocol")
	}

	return prot, nil
}
