package auth

func addBearerPrefix(token string) string {
	return TokenType + " " + token
}
