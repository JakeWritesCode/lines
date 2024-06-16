package domain

// DomainValidationErrors is a struct that contains the information about domain validation errors.
type DomainValidationErrors struct {
	Field  string   `json:"field"`
	Errors []string `json:"errors"`
}
