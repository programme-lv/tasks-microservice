package domain

import "fmt"

type DomainError struct {
	StatusCode int
	I18NErrors map[string]error
}

func (err *DomainError) Error() string {
	return fmt.Sprintf("domain error: %s", err.I18NErrors["en"])
}

// IsErrorPublic returns true if the error can't leak sensitive information
// and its contained error is therefore safe to be returned to the client.
func (err *DomainError) IsErrorPublic() bool {
	publicErrorCodes := map[int]bool{
		StateConflictErrorCode: true,
	}
	return publicErrorCodes[err.StatusCode]
}

const (
	StateConflictErrorCode = 409
)

func errorTaskFullNameIsRequired() *DomainError {
	return &DomainError{
		StatusCode: StateConflictErrorCode,
		I18NErrors: map[string]error{
			"en": fmt.Errorf("task name is required"),
			"lv": fmt.Errorf("uzdevuma nosaukums ir obligatīls"),
		},
	}
}

func errorDifficultyMustBeBetweenOneAndFive() *DomainError {
	return &DomainError{
		StatusCode: StateConflictErrorCode,
		I18NErrors: map[string]error{
			"en": fmt.Errorf("difficulty must be between 1 and 5"),
			"lv": fmt.Errorf("grūts nekorekts ar 1 un 5"),
		},
	}
}
