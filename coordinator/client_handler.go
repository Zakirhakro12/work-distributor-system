package coordinator

import (
	"net/http"
	"work-distributor-system/models"
	"work-distributor-system/session"
)

// ShowTaskForm displays the client dashboard with pending and completed tasks
func (h *Handler) ShowTaskForm(w http.ResponseWriter, r *http.Request) {
	userID := session.GetSessionValue(r, "userID").(uint)

	pendingTasks, _ := h.TaskRepo.GetSubmittedTasksByUser(userID)
	completedTasks, _ := h.TaskRepo.GetCompletedTasks(int(userID))

	err := Templates.ExecuteTemplate(w, "client_dashboard.html", struct {
		Tasks          []models.Task
		CompletedTasks []models.Task
		Success        bool
		Status         string
		ClientID       int
	}{
		Tasks:          pendingTasks,
		CompletedTasks: completedTasks,
		ClientID:       int(userID),
	})
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// SubmitTask handles new task submission by the client
func (h *Handler) SubmitTask(w http.ResponseWriter, r *http.Request) {

	// ParseMultipartForm sets a maximum upload size of 100 MB (100 * 2^20 bytes)
	r.ParseMultipartForm(100 << 20) // Limit file size to 100 MB

	userID := session.GetSessionValue(r, "userID").(uint)

	title := r.FormValue("title")
	file, handler, err := r.FormFile("file")
	var filePath string
	if err == nil {
		defer file.Close()
		filePath = "static/uploads/" + handler.Filename
		dst, _ := CreateFile(filePath)
		defer dst.Close()
		_, _ = dst.ReadFrom(file)
	}

	task := &models.Task{
		Title:     title,
		FilePath:  filePath,
		CreatedBy: uint(userID),
		Status:    "pending",
	}

	err = h.TaskRepo.CreateTask(task)
	if err != nil {
		http.Error(w, "Could not submit task", http.StatusInternalServerError)
		return
	}

	// Reload tasks and re-render
	allTasks, _ := h.TaskRepo.GetTasksByUser(int(userID))
	var pendingTasks, completedTasks []models.Task
	for _, t := range allTasks {
		if t.Status == "completed" {
			completedTasks = append(completedTasks, t)
		} else {
			pendingTasks = append(pendingTasks, t)
		}
	}

	err = Templates.ExecuteTemplate(w, "client_dashboard.html", struct {
		ClientID       uint
		Status         string
		Tasks          []models.Task
		CompletedTasks []models.Task
	}{
		ClientID:       userID,
		Tasks:          pendingTasks,
		CompletedTasks: completedTasks,
	})
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
