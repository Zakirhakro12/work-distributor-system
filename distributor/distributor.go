package distributor

import (
	"log"
	"time"
	"work-distributor-system/models"
	"work-distributor-system/repository"
)

// Worker represents a registered worker with a task channel and current load
type Worker struct {
	ID        int
	TaskCount int
	Conn      chan models.Task
}

// workers holds all active registered workers
var workers = make(map[int]*Worker)

// RegisterWorker stores the worker in the active list
func RegisterWorker(id int, conn chan models.Task) {
	if _, ok := workers[id]; ok {
		log.Printf("Worker %d already registered. Skipping.\n", id)
		return
	}

	workers[id] = &Worker{
		ID:        id,
		TaskCount: 0,
		Conn:      conn,
	}
	log.Printf("Worker %d registered.\n", id)

}

// Start runs continuously in a goroutine to assign tasks to the least-loaded available workers
func Start(taskRepo *repository.TaskRepository) {
	for {
		// Step 1: Geting all tasks that are still pending
		tasks, err := taskRepo.GetPendingTasks()
		if err != nil {
			log.Println("Distributor error fetching tasks:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Step 2: For each unassigned (pending) task
		for _, task := range tasks {
			var selectedWorker *Worker
			minTaskCount := int64(9999999) // arbitrarily large number

			// Step 3: Finding the worker with the fewest active tasks (assigned), but less than 3
			for _, worker := range workers {
				count, err := taskRepo.CountActiveTasks(worker.ID)
				if err != nil {
					log.Printf("Error counting tasks for worker %d: %v\n", worker.ID, err)
					continue
				}

				log.Printf("Worker %d - Active Tasks: %d", worker.ID, count)

				if count < 3 && count < minTaskCount {
					selectedWorker = worker
					minTaskCount = count
				}
			}

			// Step 4: Assign the task if a valid worker is available
			if selectedWorker != nil {
				// Mark task as assigned in the DB
				err := taskRepo.MarkTaskAsAssigned(int(task.ID), selectedWorker.ID)
				if err != nil {
					log.Printf("Failed to mark task %d as assigned: %v", task.ID, err)
					continue
				}

				// Geting the updated task with assigned status
				updatedTask, err := taskRepo.GetTaskByID(int(task.ID))
				if err != nil {
					log.Printf("Failed to fetch updated task %d: %v", task.ID, err)
					continue
				}

				// Sending the task to the worker via WebSocket
				selectedWorker.Conn <- *updatedTask

				// Notify the client via WebSocket
				if conn, ok := ClientConn(updatedTask.CreatedBy); ok {
					conn.WriteJSON(updatedTask)
				}

				log.Printf("Task %d assigned to worker %d", task.ID, selectedWorker.ID)
			} else {
				log.Printf("All workers busy. Task %d is waiting.\n", task.ID)
			}
		}

		// Step 5: Sleep before the next assignment cycle
		time.Sleep(5 * time.Second)
	}
}
