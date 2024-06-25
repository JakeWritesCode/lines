package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"lines/lines/domain"
	"lines/user/stores"
)

type UserDomainInterface interface {
	CreateUser(user UserForCreate) ([]domain.DomainValidationErrors, *UserData, error)
	GetUserByEmail(email string) (*UserData, error)
	GetUserByID(id uint) (*UserData, error)
	DeleteUser(id uint) error
	CheckPassword(userID uint, password string) bool
	GenerateJWT(userEmail string) (*JWTClaimsOut, error)
	BeginTransaction() error
	RollbackTransaction() error
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

type JWTClaimsOut struct {
	Email       string `json:"email"`
	TokenString string `json:"token_string"`
	jwt.RegisteredClaims
}
