package definition

// Transaction groups a pair request/response
type Transaction struct {
	Request  Request
	Response Response
}
