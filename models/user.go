package models

// User represents a registered user in the system.
// They can either have the role of "client" or "worker".
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex;not null"` // Username (must be unique)
	Password string `gorm:"not null"`             // Hashed password
	Role     string `gorm:"not null"`             // Role: either "client" or "worker"
}
