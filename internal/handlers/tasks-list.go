package handlers

import (
	"log"
	"net/http"
)

type ListTasksResponse struct {
	Tasks []Task `json:"tasks"`
}

func (c *Controller) ListTasks(w http.ResponseWriter, r *http.Request) {
	domainTaskObjs, err := c.taskSrv.ListTasks()
	if err != nil {
		log.Println("failed to list tasks", "error", err)
		respondWithJSON(w, "failed to list tasks", http.StatusInternalServerError)
		return
	}

	tasks := []Task{}
	for _, task := range domainTaskObjs {
		tasks = append(tasks, mapDomainTaskToTaskResponse(&task, c.publicBucketCloudFrontHost))
	}
	respondWithJSON(w, ListTasksResponse{
		Tasks: tasks,
	}, http.StatusOK)
}
