package domain

import "fmt"

type Task struct {
	id string

	taskFullName      string
	memoryLimitMBytes int
	cpuTimeLimitSecs  float64
	difficulty        int // [1;5]
	originOlympiad    string
	problemTags       []string
	pdfStatements     []PdfSha256Ref

	illustrationImgObjKey string
}

func (t *Task) GetIllustrationImgObjKey() string {
	return t.illustrationImgObjKey
}

func (t *Task) SetIllustrationImgObjKey(key string) {
	t.illustrationImgObjKey = key
}

func (t *Task) GetId() string {
	return t.id
}

func (t *Task) GetTaskFullName() string {
	return t.taskFullName
}

func (t *Task) GetMemoryLimitMBytes() int {
	return t.memoryLimitMBytes
}

func (t *Task) GetCpuTimeLimitSecs() float64 {
	return t.cpuTimeLimitSecs
}

func (t *Task) GetDifficulty() int {
	return t.difficulty
}

func (t *Task) GetOriginOlympiad() string {
	return t.originOlympiad
}

func (t *Task) GetProblemTags() []string {
	return t.problemTags
}

func (t *Task) GetLvOrOtherPdfSha256() string {
	if len(t.pdfStatements) == 0 {
		return ""
	}
	return t.pdfStatements[0].Sha256
}

type PdfSha256Ref struct {
	Language string
	Sha256   string
}

func NewTask(id string, fullName string) (*Task, error) {
	task := &Task{
		id:                id,
		taskFullName:      "",
		memoryLimitMBytes: 256,
		cpuTimeLimitSecs:  1.0,
		difficulty:        1,
		originOlympiad:    "",
		problemTags:       []string{},
		pdfStatements:     []PdfSha256Ref{},
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
	t.taskFullName = fullName
	return nil
}

func (t *Task) SetMemoryLimitMBytes(memoryLimit int) {
	t.memoryLimitMBytes = memoryLimit
}

func (t *Task) SetCpuTimeLimitSecs(cpuTimeLimit float64) {
	t.cpuTimeLimitSecs = cpuTimeLimit
}

func (t *Task) SetOriginOlympiad(origin string) {
	t.originOlympiad = origin
}

func (t *Task) SetPdfStatement(statements []PdfSha256Ref) {
	t.pdfStatements = statements
}

func (t *Task) SetProblemTags(tags []string) {
	t.problemTags = tags
}

func (t *Task) SetDifficulty(difficulty int) error {
	if difficulty < 1 || difficulty > 5 {
		return errorDifficultyMustBeBetweenOneAndFive()
	}
	t.difficulty = difficulty
	return nil
}
