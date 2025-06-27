package coordinator

import (
	"net/http"
	"strconv"
	"work-distributor-system/distributor"
	"work-distributor-system/models"
	"work-distributor-system/session"
)

// WorkerDashboard shows all tasks assigned to the worker and completed by them.
func (h *Handler) WorkerDashboard(w http.ResponseWriter, r *http.Request) {
	userID := session.GetSessionValue(r, "userID").(uint)

	tasks, err := h.TaskRepo.GetTasksAssignedTo(int(userID))
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	// Separate assigned and completed tasks
	var assignedTasks []models.Task
	var completedTasks []models.Task
	for _, task := range tasks {
		if task.Status == "completed" {
			completedTasks = append(completedTasks, task)
		} else {
			assignedTasks = append(assignedTasks, task)
		}
	}

	err = Templates.ExecuteTemplate(w, "worker_dashboard.html", struct {
		Assigned  []models.Task
		Completed []models.Task
		WorkerID  int
	}{
		Assigned:  assignedTasks,
		Completed: completedTasks,
		WorkerID:  int(userID),
	})
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// SubmitCompletedTask handles task submission by a worker.
// Saves the result file and updates task status to "completed".
func (h *Handler) SubmitCompletedTask(w http.ResponseWriter, r *http.Request) {

	// ParseMultipartForm sets a maximum upload size of 100 MB (100 * 2^20 bytes)
	r.ParseMultipartForm(100 << 20)
	taskIDStr := r.FormValue("task_id")

	taskID, _ := strconv.Atoi(r.FormValue("task_id"))
	userID := session.GetSessionValue(r, "userID").(uint)

	// Fetch the task from the DB
	task, err := h.TaskRepo.GetTaskByID(taskID)

	// Validating ownership and status
	if err != nil || *task.AssignedTo != userID || task.Status != "assigned" {
		http.Error(w, "Unauthorized or invalid task", http.StatusForbidden)
		return
	}

	// Reading uploaded file
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Saving file to disk
	filePath := "static/completed/" + handler.Filename
	dst, _ := CreateFile(filePath)
	defer dst.Close()
	_, _ = dst.ReadFrom(file)

	// Convert task ID
	taskID, err = strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Updating the status task as completed
	err = h.TaskRepo.MarkTaskAsCompleted(taskID, filePath)
	if err != nil {
		http.Error(w, "Failed to mark task completed", http.StatusInternalServerError)
		return
	}

	// Notify the original client via WebSocket (real-time update)
	task, _ = h.TaskRepo.GetTaskByID(taskID)
	if task.CreatedBy != 0 {
		if conn, ok := distributor.ClientConn(task.CreatedBy); ok {
			err = conn.WriteJSON(task)
			if err != nil {
				http.Error(w, "Failed to send task update", http.StatusInternalServerError)
				return
			}
		}
	}

	http.Redirect(w, r, "/worker-dashboard", http.StatusFound)
}
