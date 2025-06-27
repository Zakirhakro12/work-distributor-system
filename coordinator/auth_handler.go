package coordinator

import (
	"net/http"
	"work-distributor-system/models"
	"work-distributor-system/session"

	"golang.org/x/crypto/bcrypt"
)

// ShowLogin renders the login page
func (h *Handler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	Templates.ExecuteTemplate(w, "login.html", nil)
}

// Login authenticates the user and sets session values
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Fetching user from DB
	user, err := h.UserRepo.GetUserByUsername(username)
	// Checking credentials using bcrypt
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Storing session values
	session.SetSessionValue(w, r, "username", username)
	session.SetSessionValue(w, r, "role", user.Role)
	session.SetSessionValue(w, r, "userID", user.ID)

	// Redirecting based on user role
	if user.Role == "client" {
		http.Redirect(w, r, "/client-dashboard", http.StatusFound)
	} else if user.Role == "worker" {
		http.Redirect(w, r, "/worker-dashboard", http.StatusFound)
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

// Logout clears session and redirects to login
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session.ClearSession(w, r)
	http.Redirect(w, r, "/login", http.StatusFound)
}

// ShowRegister displays the registration form
func (h *Handler) ShowRegister(w http.ResponseWriter, r *http.Request) {
	Templates.ExecuteTemplate(w, "register.html", nil)
}

// Register processes new user registration and stores hashed password
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user := &models.User{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Role:     r.FormValue("role"),
	}

	// Hashing the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Saving user to DB
	err = h.UserRepo.CreateUser(user)

	if err != nil {
		http.Error(w, "User exists or error", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}
