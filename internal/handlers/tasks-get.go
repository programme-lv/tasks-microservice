package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/programme-lv/tasks-microservice/internal/domain"
)

type GetTaskResponse struct {
	Task Task `json:"task"`
}

type Task struct {
	PublishedTaskId    string            `json:"published_task_id"`
	TaskFullName       string            `json:"task_full_name"`
	MemoryLimitMbytes  int               `json:"memory_limit_megabytes"`
	CpuTimeLimitSecs   float64           `json:"cpu_time_limit_seconds"`
	OriginOlympiad     string            `json:"origin_olympiad,omitempty"`
	LvPdfStatementSha  string            `json:"lv_pdf_statement_sha,omitempty"`
	DifficultyRating   int               `json:"difficulty_rating,omitempty"`
	IllustrationImgUrl string            `json:"illustration_img_url,omitempty"`
	DefaultMdStatement *MdStatement      `json:"default_md_statement,omitempty"`
	DefaultPdfSUrl     string            `json:"default_pdf_statement_url,omitempty"`
	Examples           []Example         `json:"examples,omitempty"`
	OriginNotes        map[string]string `json:"origin_notes,omitempty"`
	VisInpStInputs     []StInputs        `json:"visible_input_subtasks,omitempty"`
}

type StInputs struct {
	Subtask int      `json:"subtask"`
	Inputs  []string `json:"inputs"`
}

type Example struct {
	Input  string  `json:"input"`
	Output string  `json:"output"`
	MdNote *string `json:"md_note,omitempty"`
}

type MdStatement struct {
	Story   string  `json:"story"`
	Input   string  `json:"input"`
	Output  string  `json:"output"`
	Notes   *string `json:"notes,omitempty"`
	Scoring *string `json:"scoring,omitempty"`
}

func (c *Controller) GetTask(w http.ResponseWriter, r *http.Request) {
	// take id param from url chi

	id := chi.URLParam(r, "id")
	if id == "" {
		respondWithJSON(w, "invalid task id", http.StatusBadRequest)
		return
	}

	task, err := c.taskSrv.GetTask(id)
	if err != nil {
		respondWithJSON(w, "task not found", http.StatusNotFound)
		return
	}

	respondWithJSON(w, GetTaskResponse{
		Task: mapDomainTaskToTaskResponse(task, c.publicBucketCloudFrontHost),
	}, http.StatusOK)
}

func mapDomainTaskToTaskResponse(task *domain.Task, publicBucketCloudFrontHost string) Task {
	illustrationImgUrl := ""
	if publicBucketCloudFrontHost != "" && task.GetIllustrationImgObjKey() != "" {
		illustrationImgUrl = fmt.Sprintf("https://%s/%s",
			publicBucketCloudFrontHost, task.GetIllustrationImgObjKey())
	}

	examples := make([]Example, 0)
	for _, example := range task.GetExamples() {
		examples = append(examples, Example{
			Input:  example.Input,
			Output: example.Output,
			MdNote: example.MdNote,
		})
	}

	mdStImgUuidToObjKey := task.GetImgUuidToObjKey()
	mdStatement := task.GetDefaultMarkdownStatement()
	if mdStatement != nil {
		for _, section := range []*string{&mdStatement.Story, &mdStatement.Input, &mdStatement.Output,
			mdStatement.Notes, mdStatement.Scoring} {
			if section != nil {
				for imgUuid, objKey := range mdStImgUuidToObjKey {
					url := fmt.Sprintf("https://%s/%s", publicBucketCloudFrontHost, objKey)
					*section = strings.ReplaceAll(*section, imgUuid, url)
				}
			}
		}
	}

	var resMdStatement *MdStatement = nil
	if mdStatement != nil {
		resMdStatement = &MdStatement{
			Story:   mdStatement.Story,
			Input:   mdStatement.Input,
			Output:  mdStatement.Output,
			Notes:   mdStatement.Notes,
			Scoring: mdStatement.Scoring,
		}
	}

	defaultPdfStatementUrl := ""
	if publicBucketCloudFrontHost != "" && task.GetLvOrOtherPdfSha256() != "" {
		defaultPdfStatementUrl = fmt.Sprintf(
			"https://%s/task-pdf-statements/%s.pdf",
			publicBucketCloudFrontHost, task.GetLvOrOtherPdfSha256())
	}

	visInpStInputs := make([]StInputs, 0)
	for _, visInpSt := range task.GetVisInpStInputs() {
		visInpSt := StInputs{
			Subtask: visInpSt.Subtask,
			Inputs:  visInpSt.Inputs,
		}
		visInpStInputs = append(visInpStInputs, visInpSt)
	}

	return Task{
		PublishedTaskId:    task.GetId(),
		TaskFullName:       task.GetTaskFullName(),
		MemoryLimitMbytes:  task.GetMemoryLimitMBytes(),
		CpuTimeLimitSecs:   task.GetCpuTimeLimitSecs(),
		OriginOlympiad:     task.GetOriginOlympiad(),
		LvPdfStatementSha:  task.GetLvOrOtherPdfSha256(),
		DifficultyRating:   task.GetDifficulty(),
		IllustrationImgUrl: illustrationImgUrl,
		DefaultMdStatement: resMdStatement,
		DefaultPdfSUrl:     defaultPdfStatementUrl,
		Examples:           examples,
		OriginNotes:        task.GetOriginNotes(),
		VisInpStInputs:     visInpStInputs,
	}
}
