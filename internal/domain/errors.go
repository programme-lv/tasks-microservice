package domain

import "fmt"

type DomainError struct {
	StatusCode int
	I18NErrors map[string]error
}

func (err *DomainError) Error() string {
	return fmt.Sprintf("domain error: %s", err.I18NErrors["en"])
}

const (
	StateConflictErrorCode = 409
)

func errorTaskFullNameIsRequired() *DomainError {
	return &DomainError{
		StatusCode: StateConflictErrorCode,
		I18NErrors: map[string]error{
			"en": fmt.Errorf("task name is required"),
			"lv": fmt.Errorf("uzdevuma nosaukums ir obligāts"),
		},
	}
}

func errorDifficultyMustBeBetweenOneAndFive() *DomainError {
	return &DomainError{
		StatusCode: StateConflictErrorCode,
		I18NErrors: map[string]error{
			"en": fmt.Errorf("difficulty must be between 1 and 5"),
			"lv": fmt.Errorf("grūtibai jābūt starp 1 un 5"),
		},
	}
}

func errorEmptyTestSha256() *DomainError {
	return &DomainError{
		StatusCode: StateConflictErrorCode,
		I18NErrors: map[string]error{
			"en": fmt.Errorf("test sha256 is required"),
			"lv": fmt.Errorf("testa sha256 ir obligāts"),
		},
	}
}

func errorTestIdMustBePositive() *DomainError {
	return &DomainError{
		StatusCode: StateConflictErrorCode,
		I18NErrors: map[string]error{
			"en": fmt.Errorf("test id must be positive"),
			"lv": fmt.Errorf("testa id jābūt pozitīvam"),
		},
	}
}
