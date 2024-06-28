package domain

import (
	"github.com/stretchr/testify/assert"
	"lines/lines/store"
	"testing"
)

func TestAddValidationError_NewField(t *testing.T) {
	var errors []DomainValidationErrors
	errors = AddValidationError("test", "test error", errors)
	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %v", len(errors))
	}
	if errors[0].Field != "test" {
		t.Errorf("Expected field to be 'test', got %v", errors[0].Field)
	}
	if errors[0].Errors[0] != "test error" {
		t.Errorf("Expected error to be 'test error', got %v", errors[0].Errors[0])
	}
}

func TestAddValidationError_ExistingField(t *testing.T) {
	var errors []DomainValidationErrors
	errors = AddValidationError("test", "test error", errors)
	errors = AddValidationError("test", "test error 2", errors)
	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %v", len(errors))
	}
	if errors[0].Field != "test" {
		t.Errorf("Expected field to be 'test', got %v", errors[0].Field)
	}
	if errors[0].Errors[0] != "test error" {
		t.Errorf("Expected error to be 'test error', got %v", errors[0].Errors[0])
	}
	if errors[0].Errors[1] != "test error 2" {
		t.Errorf("Expected error to be 'test error 2', got %v", errors[0].Errors[1])
	}
}

func TestEmptyStringValidator_EmptyString(t *testing.T) {
	var errors []DomainValidationErrors
	errors = EmptyStringValidator("", "test", errors)

	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %v", len(errors))
	}
	if errors[0].Field != "test" {
		t.Errorf("Expected field to be 'test', got %v", errors[0].Field)
	}
	if errors[0].Errors[0] != "test is required" {
		t.Errorf("Expected error to be 'test is required', got %v", errors[0].Errors[0])
	}
}

func TestEmptyStringValidator_NonEmptyString(t *testing.T) {
	var errors []DomainValidationErrors
	errors = EmptyStringValidator("test", "test", errors)
	if len(errors) != 0 {
		t.Errorf("Expected 0 errors, got %v", len(errors))
	}
}

func TestEmailValidator_EmptyString(t *testing.T) {
	var errors []DomainValidationErrors
	errors = EmailValidator("", "test", errors)
	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %v", len(errors))
	}
	if errors[0].Field != "test" {
		t.Errorf("Expected field to be 'test', got %v", errors[0].Field)
	}
	if errors[0].Errors[0] != "test is required" {
		t.Errorf("Expected error to be 'test is required', got %v", errors[0].Errors[0])
	}
}

func TestEmailValidator_InvalidEmail(t *testing.T) {
	var errors []DomainValidationErrors
	errors = EmailValidator("test", "test", errors)
	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %v", len(errors))
	}
	if errors[0].Field != "test" {
		t.Errorf("Expected field to be 'test', got %v", errors[0].Field)
	}
	if errors[0].Errors[0] != "Invalid email" {
		t.Errorf("Expected error to be 'Invalid email', got %v", errors[0].Errors[0])
	}
}

func TestEmailValidator_ValidEmail(t *testing.T) {
	var errors []DomainValidationErrors
	errors = EmailValidator("some@email.com", "test", errors)
	assert.Equal(t, 0, len(errors))
}

func TestStoreValidationErrorToDomainValidationError(t *testing.T) {
	storeErrors := []store.ModelValidationError{
		{Field: "test", Message: "test error"},
	}
	domainErrors := StoreValidationErrorToDomainValidationError(storeErrors)
	assert.Equal(t, 1, len(domainErrors))
	assert.Equal(t, "test", domainErrors[0].Field)
	assert.Equal(t, 1, len(domainErrors[0].Errors))
	assert.Equal(t, "test error", domainErrors[0].Errors[0])
}
