package domain

import (
	"lines/lines/domain"
	"lines/user/stores"
)

type UserDomainInterface interface {
	CreateUser() error
}

type UserForCreate struct {
	Name     string
	Email    string
	Password string
}

func (u UserForCreate) Validate(store stores.UserStoreInterface) ([]domain.DomainValidationErrors, error) {
	var validationErrors []domain.DomainValidationErrors
	validationErrors = domain.EmptyStringValidator(u.Name, "name", validationErrors)
	validationErrors = domain.EmailValidator(u.Email, "email", validationErrors)
	validationErrors = domain.EmptyStringValidator(u.Password, "password", validationErrors)

	// Check if the email is already in use.
	user, err := store.GetUserByEmail(u.Email)
	if err != nil {
		return validationErrors, err
	}
	if user != nil {
		validationErrors = domain.AddValidationError("email", "Email is already in use.", validationErrors)
	}

	return validationErrors, nil
}

type UserData struct {
	ID    uint
	Name  string
	Email string
}
