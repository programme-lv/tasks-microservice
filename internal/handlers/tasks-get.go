package handlers

import (
	"net/http"
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
}

func (c *Controller) GetTask(w http.ResponseWriter, r *http.Request) {

}
