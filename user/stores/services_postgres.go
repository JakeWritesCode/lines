package stores

import "lines/lines/store"

func (s *UserPostgresStore) CreateUser(user *User) ([]store.ModelValidationError, error) {
	validationErrors := user.Validate()
	if len(validationErrors) > 0 {
		return validationErrors, nil
	}
	return []store.ModelValidationError{}, s.Postgres.Create(user).Error
}

func (s *UserPostgresStore) GetUserByEmail(email string) (*User, error) {
	var user User
	err := s.Postgres.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserPostgresStore) GetUserByID(id uint) (*User, error) {
	var user User
	result := s.Postgres.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (s *UserPostgresStore) UpdateUser(user *User) ([]store.ModelValidationError, error) {
	validationErrors := user.Validate()
	if len(validationErrors) > 0 {
		return validationErrors, nil
	}
	return []store.ModelValidationError{}, s.Postgres.Save(user).Error
}
func (s *UserPostgresStore) DeleteUser(user *User) error {
	return s.Postgres.Delete(user).Error
}
