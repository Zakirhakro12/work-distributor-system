package models

// Task represents a unit of work in the system submitted by a client and assigned to a worker
type Task struct {
	ID                uint   `gorm:"primaryKey"` // Unique identifier for each task
	Title             string `gorm:"not null"`   // Short title/summary of the task
	FilePath          string // Path to the uploaded file by the client
	Status            string `gorm:"default:pending"` // Task status: "pending", "assigned","completed"
	CreatedBy         uint   // User ID of the client who created the task
	AssignedTo        *uint  // Nullable worker ID to whom the task is assigned
	CompletedFilePath string // Path to the result file submitted by the worker
}
