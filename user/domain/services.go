package domain

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"lines/lines/domain"
	"lines/user/stores"
)

func HashAndSalt(password string) (string, error) {
	passwordBytes := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (u *UserDomain) CreateUser(user UserForCreate) ([]domain.DomainValidationErrors, *UserData, error) {
	// Validate the user.
	validationErrors, err := user.Validate(u.store)
	if err != nil {
		return validationErrors, nil, err
	}
	if len(validationErrors) > 0 {
		return validationErrors, nil, nil
	}
	hashed, err := HashAndSalt(user.Password)
	if err != nil {
		return nil, nil, errors.New("could not hash password")
	}
	storeUser := stores.User{
		Email:    user.Email,
		Name:     user.Name,
		Password: hashed,
	}
	modelErrors, err := u.store.CreateUser(&storeUser)
	if err != nil {
		return nil, nil, err
	}
	if len(modelErrors) > 0 {
		return domain.StoreValidationErrorToDomainValidationError(modelErrors), nil, nil
	}
	return nil, &UserData{
		ID:    storeUser.ID,
		Email: storeUser.Email,
		Name:  storeUser.Name,
	}, nil
}

func (u *UserDomain) GetUserByEmail(email string) (*UserData, error) {
	storeUser, err := u.store.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if storeUser == nil {
		return nil, nil
	}
	return &UserData{
		ID:    storeUser.ID,
		Email: storeUser.Email,
		Name:  storeUser.Name,
	}, nil
}

func (u *UserDomain) GetUserByID(id uint) (*UserData, error) {
	storeUser, err := u.store.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	if storeUser == nil {
		return nil, nil
	}
	return &UserData{
		ID:    storeUser.ID,
		Email: storeUser.Email,
		Name:  storeUser.Name,
	}, nil
}

func (u *UserDomain) DeleteUser(id uint) error {
	storeUser, err := u.store.GetUserByID(id)
	if err != nil {
		return err
	}
	if storeUser == nil {
		return nil
	}
	return u.store.DeleteUser(storeUser)
}
