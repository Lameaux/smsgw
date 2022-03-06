package auth

type Auth interface {
	Authorize(apiKey string) (string, error)
}
