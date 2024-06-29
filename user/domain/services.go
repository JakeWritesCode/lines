package domain

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"lines/lines/domain"
	linesHttp "lines/lines/http"
	"lines/user/stores"
	"net/http"
	"strings"
	"time"
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

func (u *UserDomain) CheckPassword(userID uint, password string) bool {
	user, err := u.store.GetUserByID(userID)
	if err != nil || user == nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (u *UserDomain) GenerateJWT(userEmail string) (*JWTClaimsOut, error) {
	claims := JWTClaimsOut{
		Email: userEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(u.Config.TokenExpirationTimeMinutes) * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(u.Config.SecretKey)
	if err != nil {
		return nil, err
	}
	claims.TokenString = tokenString
	return &claims, nil
}

func (u *UserDomain) GetJWTFromRequest(r http.Request) (string, error) {
	// Cookie
	tokenString, err := r.Cookie("Bearer")
	if err == nil {
		return tokenString.Value, nil
	}
	// Header
	tokenStringHeader := r.Header.Get("Authorization")
	if tokenStringHeader != "" {
		return strings.Replace(tokenStringHeader, "Bearer ", "", 1), nil
	}
	return "", errors.New("no token found")
}

func (u *UserDomain) ValidateJWT(token string) (*linesHttp.HttpError, *JWTClaimsOut) {
	claims := &JWTClaimsOut{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return u.Config.SecretKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return &linesHttp.HttpError{
				Message: []string{"Unauthorised"},
			}, nil
		}
		return &linesHttp.HttpError{
			Message: []string{"Bearer token invalid"},
		}, nil
	}
	if !parsedToken.Valid {
		return &linesHttp.HttpError{
			Message: []string{"Unauthorised"},
		}, nil
	}
	return nil, claims
}

func (u *UserDomain) ValidateRequestAuth(r http.Request) (*linesHttp.HttpError, *JWTClaimsOut) {
	tokenString, err := u.GetJWTFromRequest(r)
	if err != nil {
		return &linesHttp.HttpError{
			Message: []string{"Unauthorised"},
		}, nil
	}
	return u.ValidateJWT(tokenString)
}

func (u *UserDomain) ChangeUserPassword(userID uint, oldPassword string, newPassword string) ([]domain.DomainValidationErrors, error) {
	if !u.CheckPassword(userID, oldPassword) {
		return []domain.DomainValidationErrors{
			{
				Field:  "old_password",
				Errors: []string{"Old password is incorrect"},
			},
		}, nil
	}
	hashed, err := HashAndSalt(newPassword)
	if err != nil {
		return nil, errors.New("could not hash password")
	}

	storeUser, err := u.store.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	storeUser.Password = hashed
	modelErrors, err := u.store.UpdateUser(storeUser)
	if err != nil {
		return nil, err
	}
	if len(modelErrors) > 0 {
		return domain.StoreValidationErrorToDomainValidationError(modelErrors), nil
	}
	return nil, nil
}
