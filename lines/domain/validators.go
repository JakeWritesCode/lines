package domain

import (
	"lines/lines/store"
	"regexp"
)

type GenericValidator func(value string, fieldName string, errors []DomainValidationErrors) []DomainValidationErrors

// AddValidationError is a function that adds a validation error to the list of validation errors.
// If the field already has a validation error, the error is appended to the list of errors.
// If the field does not have a validation error, a new validation error is created.
func AddValidationError(fieldName string, error string, errors []DomainValidationErrors) []DomainValidationErrors {
	for i := range errors {
		if errors[i].Field == fieldName {
			errors[i].Errors = append(errors[i].Errors, error)
			return errors
		}
	}
	errors = append(errors, DomainValidationErrors{
		Field:  fieldName,
		Errors: []string{error},
	})
	return errors
}

// EmptyStringValidator is a function that validates if a string is empty.
func EmptyStringValidator(value string, fieldName string, errors []DomainValidationErrors) []DomainValidationErrors {
	if value == "" {
		return AddValidationError(fieldName, fieldName+" is required", errors)
	}
	return errors
}

// EmailValidator is a function that validates if a string is a valid email.
func EmailValidator(value string, fieldName string, errors []DomainValidationErrors) []DomainValidationErrors {
	if value == "" {
		errors = AddValidationError(fieldName, fieldName+" is required", errors)
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		errors = AddValidationError(fieldName, "Invalid email", errors)
	}
	return errors
}

func StoreValidationErrorToDomainValidationError(storeErrors []store.ModelValidationError) []DomainValidationErrors {
	var domainErrors []DomainValidationErrors
	for _, storeError := range storeErrors {
		domainErrors = AddValidationError(storeError.Field, storeError.Message, domainErrors)
	}
	return domainErrors
}
