package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/programme-lv/tasks-microservice/internal/domain"
)

type GetTaskResponse struct {
	Task Task `json:"task"`
}

type Task struct {
	PublishedTaskId   string  `json:"published_task_id"`
	TaskFullName      string  `json:"task_full_name"`
	MemoryLimitMbytes int     `json:"memory_limit_megabytes"`
	CpuTimeLimitSecs  float64 `json:"cpu_time_limit_seconds"`
	OriginOlympiad    string  `json:"origin_olympiad,omitempty"`
	LvPdfStatementSha string  `json:"lv_pdf_statement_sha,omitempty"`
	DifficultyRating  int     `json:"difficulty_rating,omitempty"`
}

func (c *Controller) GetTask(w http.ResponseWriter, r *http.Request) {
	// take id param from url chi

	id := chi.URLParam(r, "id")
	if id == "" {
		respondWithBadRequest(w, "invalid task id")
		return
	}

	task, err := c.TaskSrv.GetTask(id)
	if err != nil {
		respondWithBadRequest(w, "task not found")
		return
	}

	respondWithJSON(w, GetTaskResponse{
		Task: mapDomainTaskToTaskResponse(task),
	}, http.StatusOK)
}

func mapDomainTaskToTaskResponse(task *domain.Task) Task {
	return Task{
		PublishedTaskId:   task.GetId(),
		TaskFullName:      task.GetTaskFullName(),
		MemoryLimitMbytes: task.GetMemoryLimitMBytes(),
		CpuTimeLimitSecs:  task.GetCpuTimeLimitSecs(),
		OriginOlympiad:    task.GetOriginOlympiad(),
		LvPdfStatementSha: task.GetLvOrOtherPdfSha256(),
		DifficultyRating:  task.GetDifficulty(),
	}
}
