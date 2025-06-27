package repository

import (
	"work-distributor-system/models"

	"gorm.io/gorm"
)

// it provides database methods to manage Task records.
type TaskRepository struct {
	DB *gorm.DB
}

// it initializes a new TaskRepository with the provided database connection.
func NewTaskRepo(db *gorm.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

// inserting a new task into the database.
func (r *TaskRepository) CreateTask(task *models.Task) error {
	return r.DB.Create(task).Error
}

// it returns all tasks (except completed) created by a specific user (client).
func (r *TaskRepository) GetSubmittedTasksByUser(userID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := r.DB.Where("created_by = ? AND status != ?", userID, "completed").Find(&tasks).Error
	return tasks, err
}

// it returns all tasks created by a specific user (client), regardless of status.
func (r *TaskRepository) GetTasksByUser(userID int) ([]models.Task, error) {
	var tasks []models.Task
	err := r.DB.Where("created_by = ? ", userID).Find(&tasks).Error
	return tasks, err
}

// it fetches all tasks assigned to a given worker.
func (r *TaskRepository) GetTasksAssignedTo(userID int) ([]models.Task, error) {
	var tasks []models.Task
	err := r.DB.Where("assigned_to = ?", userID).Find(&tasks).Error
	return tasks, err
}

// it returns tasks created by a specific client that are marked as "completed".
func (r *TaskRepository) GetCompletedTasks(userID int) ([]models.Task, error) {
	var tasks []models.Task
	err := r.DB.Where("created_by = ? AND status = ?", userID, "completed").Find(&tasks).Error
	return tasks, err
}

// it updates the task to set the assigned worker and status
func (r *TaskRepository) MarkTaskAsAssigned(taskID int, workerID int) error {
	return r.DB.Model(&models.Task{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"assigned_to": workerID,
		"status":      "assigned",
	}).Error
}

// it updates the task status and saves the path to the completed file.
func (r *TaskRepository) MarkTaskAsCompleted(taskID int, filePath string) error {
	return r.DB.Model(&models.Task{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"status":              "completed",
		"completed_file_path": filePath,
	}).Error
}

// it fetches all unassigned or pending tasks (i.e., available for assignment).
func (r *TaskRepository) GetPendingTasks() ([]models.Task, error) {
	var tasks []models.Task
	err := r.DB.Where("status = ?", "pending").Find(&tasks).Error
	return tasks, err
}

// it returns the number of tasks assigned to a worker that are not completed.
func (r *TaskRepository) CountActiveTasks(workerID int) (int64, error) {
	var count int64
	err := r.DB.Model(&models.Task{}).
		Where("assigned_to = ? AND status IN ?", workerID, []string{"assigned", "pending"}).
		Count(&count).Error
	return count, err
}

// retrieves a single task by its ID.
func (r *TaskRepository) GetTaskByID(taskID int) (*models.Task, error) {
	var task models.Task
	result := r.DB.First(&task, taskID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &task, nil
}
