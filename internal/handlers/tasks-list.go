package handlers

import (
	"net/http"
)

type ListTasksResponse struct {
	Tasks []Task `json:"tasks"`
}

func (c *Controller) ListTasks(w http.ResponseWriter, r *http.Request) {
	domainTaskObjs, err := c.TaskSrv.ListTasks()
	if err != nil {
		// TODO: handle error in a better way
		// slog.Warn("failed to list tasks", "error", err)
		// domainErr := domain.DomainError{}
		// if errors.As(err, &domainErr) {
		// 	if domainErr.IsErrorPublic() {
		// 		respond
		// 	}
		// }
		respondWithInternalServerError(w, "failed to list tasks")
		return
	}

	tasks := []Task{}
	for _, task := range domainTaskObjs {
		tasks = append(tasks, Task{
			TaskFullName:      task.GetTaskFullName(),
			MemoryLimitMbytes: task.GetMemoryLimitMBytes(),
			CpuTimeLimitSecs:  task.GetCpuTimeLimitSecs(),
			OriginOlympiad:    task.GetOriginOlympiad(),
			LvPdfStatementSha: task.GetLvOrOtherPdfSha256(),
		})
	}
	respondWithJSON(w, ListTasksResponse{
		Tasks: tasks,
	}, http.StatusOK)
}

// func (c *Controller) ListUsers(w http.ResponseWriter, r *http.Request) {
// 	user, err := c.UserService.ListUsers()
// 	if err != nil {
// 		respondWithBadRequest(w, "user not found")
// 		return
// 	}

// 	users := []Task{}
// 	for _, u := range user {
// 		users = append(users, Task{})
// 	}

// 	response := ListTasksResponse{
// 		Users: users,
// 	}

// 	respondWithJSON(w, response, http.StatusOK)
// }
