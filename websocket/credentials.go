package websocket

import (
	"strings"
)

// extractCredentials splits the credentials string into username and password
func extractCredentials(credentials string) (string, string) {
	// Split credentials string into username and password
	// You can use a delimiter or a more complex parsing logic based on your requirements
	// Here's a simple example assuming the format is "username:password"
	parts := strings.Split(credentials, ":")
	username := parts[0]
	password := parts[1]
	return username, password
}
