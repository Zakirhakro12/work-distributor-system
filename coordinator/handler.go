package coordinator

import (
	"work-distributor-system/repository"
)

// Handler holds references to the task and user repositories
// This is injected into route handlers for reuse
type Handler struct {
	TaskRepo *repository.TaskRepository
	UserRepo *repository.UserRepository
}

// NewHandler constructs a new Handler instance with repositories
func NewHandler(taskRepo *repository.TaskRepository, userRepo *repository.UserRepository) *Handler {
	return &Handler{TaskRepo: taskRepo, UserRepo: userRepo}
}
