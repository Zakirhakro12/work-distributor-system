package main

import (
	"log"
	"net/http"
	"os"
	"work-distributor-system/config"
	"work-distributor-system/coordinator"
	"work-distributor-system/distributor"
	"work-distributor-system/middleware"
	"work-distributor-system/repository"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	coordinator.InitTemplates()

	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using defaults")
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "tasks.db" // default value
	}

	// Initializing the SQLite database and auto-migrating tables
	db := config.SetupDatabase(dbPath)

	// Initializing repositories for Task and User
	taskRepo := repository.NewTaskRepo(db)
	userRepo := repository.NewUserRepo(db)

	// Starting the background distributor (that is responsible for assigning tasks to workers)
	go distributor.Start(taskRepo)

	// Seting up the HTTP handlers with dependencies injected
	handler := coordinator.NewHandler(taskRepo, userRepo)

	// Creating the main router
	r := mux.NewRouter()

	// Public routes
	public := r.PathPrefix("").Subrouter()
	public.HandleFunc("/register", handler.ShowRegister).Methods("GET") // Showing registration form
	public.HandleFunc("/register", handler.Register).Methods("POST")    // Processing registration
	public.HandleFunc("/login", handler.ShowLogin).Methods("GET")       // Showing login form
	public.HandleFunc("/login", handler.Login).Methods("POST")          // Processing login

	// Protected Routes (require login session)
	protected := r.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware) // applies to all routes in this group

	// Logout endpoint
	protected.HandleFunc("/logout", handler.Logout).Methods("GET")

	// Client dashboard (only for users with role "client")
	protected.HandleFunc("/client-dashboard", middleware.RequireRole(handler.ShowTaskForm, "client")).Methods("GET")
	protected.HandleFunc("/client-dashboard", middleware.RequireRole(handler.SubmitTask, "client")).Methods("POST")

	// Worker dashboard (only for users with role "worker")
	protected.HandleFunc("/worker-dashboard", middleware.RequireRole(handler.WorkerDashboard, "worker")).Methods("GET")
	protected.HandleFunc("/submit-completed", middleware.RequireRole(handler.SubmitCompletedTask, "worker")).Methods("POST")

	// WebSocket endpoints for real-time communication
	protected.HandleFunc("/ws", distributor.WorkerWebSocket)
	protected.HandleFunc("/client-ws", distributor.ClientWebSocket)

	// route handler for the root path
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

	// Serving static files (uploaded/downloaded files)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // fallback
	}
	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
