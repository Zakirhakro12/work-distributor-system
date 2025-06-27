package session

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

// store is a session store that uses secure cookies to manage session data.
// Loading .env file

var store *sessions.CookieStore

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found. Using environment variables.")
	}

	// Get session secret from environment
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		log.Fatal("SESSION_SECRET must be set in environment variables")
	}

	// Initialize the session store with the secret
	store = sessions.NewCookieStore([]byte(secret))
}

// GetSessionValue retrieves a value from the session using the provided key.
func GetSessionValue(r *http.Request, key string) interface{} {
	session, _ := store.Get(r, "session")
	return session.Values[key]
}

// SetSessionValue sets a key-value pair in the session and saves it to the response.
// Useful for storing login data such as userID or role.
func SetSessionValue(w http.ResponseWriter, r *http.Request, key string, value interface{}) {
	session, _ := store.Get(r, "session")
	session.Values[key] = value
	session.Save(r, w) // persists the session in the browser as a cookie
}

// ClearSession removes the session by setting MaxAge to -1 and saving.
// This is used during logout to invalidate the session.
func ClearSession(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Options.MaxAge = -1 // Invalidate the session immediately
	session.Save(r, w)
}
