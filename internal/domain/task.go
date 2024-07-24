package domain

import "fmt"

type Task struct {
	Id string

	TaskFullName      string
	MemoryLimitMBytes int
	CpuTimeLimitSecs  float64
	Difficulty        int // [1;5]
	OriginOlympiad    string
	ProblemTags       []string
	PdfStatements     []PdfSha256Ref
}

type PdfSha256Ref struct {
	Language string
	Sha256   string
}

func NewTask(id string, fullName string) (*Task, error) {
	task := &Task{
		Id:                id,
		TaskFullName:      "",
		MemoryLimitMBytes: 256,
		CpuTimeLimitSecs:  1.0,
		Difficulty:        1,
		OriginOlympiad:    "",
		ProblemTags:       []string{},
		PdfStatements:     []PdfSha256Ref{},
	}

	err := task.SetTaskFullName(fullName)
	if err != nil {
		return nil, fmt.Errorf("failed to set task full name: %w", err)
	}

	return task, nil
}

func (t *Task) SetTaskFullName(fullName string) error {
	if fullName == "" {
		return errorTaskFullNameIsRequired()
	}
	t.TaskFullName = fullName
	return nil
}

func (t *Task) SetMemoryLimitMBytes(memoryLimit int) {
	t.MemoryLimitMBytes = memoryLimit
}

func (t *Task) SetCpuTimeLimitSecs(cpuTimeLimit float64) {
	t.CpuTimeLimitSecs = cpuTimeLimit
}

func (t *Task) SetOriginOlympiad(origin string) {
	t.OriginOlympiad = origin
}

func (t *Task) SetPdfStatement(statements []PdfSha256Ref) {
	t.PdfStatements = statements
}

func (t *Task) SetProblemTags(tags []string) {
	t.ProblemTags = tags
}

func (t *Task) SetDifficulty(difficulty int) error {
	if difficulty < 1 || difficulty > 5 {
		return errorDifficultyMustBeBetweenOneAndFive()
	}
	t.Difficulty = difficulty
	return nil
}
