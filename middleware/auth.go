package middleware

import (
	"net/http"
	"work-distributor-system/session" // use this for session helpers
)

// AuthMiddleware ensures the user is logged in (has a valid session)
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userID := session.GetSessionValue(r, "userID")
		if userID == nil {
			// Redirecting to login if no session is found
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Proceeding to the next handler
		next.ServeHTTP(w, r)
	})
}

// RequireRole ensures the user has a specific role (e.g., client or worker)
func RequireRole(next http.HandlerFunc, role string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userRole, ok := session.GetSessionValue(r, "role").(string)
		if !ok || userRole != role {
			// Redirect to login if the role does not match
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		// Role matches, proceed to the requested handler
		next(w, r)
	}
}
