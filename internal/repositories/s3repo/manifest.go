package s3repo

type TaskTomlManifest struct {
	TestSHA256s  []TestfileSHA256Ref     `toml:"tests_sha256s"`
	PDFSHA256s   []PDFStatemenSHA256tRef `toml:"pdf_statements_sha256s"`
	MDStatements []MDStatement           `toml:"md_statements"`

	TaskFullName    string      `toml:"task_full_name"`
	MemoryLimMB     int         `toml:"memory_lim_megabytes"`
	CpuTimeInSecs   float64     `toml:"cpu_time_in_seconds"`
	ProblemTags     []string    `toml:"problem_tags"`
	Difficulty      int         `toml:"difficulty_1_to_5"`
	TaskAuthors     []string    `toml:"task_authors"`
	OriginOlympiad  string      `toml:"origin_olympiad"`
	VisibleInputSTs []int       `toml:"visible_input_subtasks"`
	TestGroups      []TestGroup `toml:"test_groups"`
}

type TestfileSHA256Ref struct {
	TestID       int    `toml:"test_id"`
	InputSHA256  string `toml:"input_sha256"`
	AnswerSHA256 string `toml:"answer_sha256"`
}

type PDFStatemenSHA256tRef struct {
	Language string `toml:"language"`
	SHA256   string `toml:"sha256"`
}

type TestGroup struct {
	GroupID int   `toml:"group_id"`
	Points  int   `toml:"points"`
	Public  bool  `toml:"public"`
	Subtask int   `toml:"subtask"`
	TestIDs []int `toml:"test_ids"`
}

type MDStatement struct {
	Language *string `toml:"language"`
	Story    string  `toml:"story"`
	Input    string  `toml:"input"`
	Output   string  `toml:"output"`
	Notes    *string `toml:"notes"`
	Scoring  *string `toml:"scoring"`
}
