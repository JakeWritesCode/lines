package domain

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"lines/lines/store"
	"lines/user/stores"
	"net/http"
	"testing"
)

func TestHashAndSalt(t *testing.T) {
	password := "password"
	hash, err := HashAndSalt(password)
	assert.Nil(t, err)
	assert.NotEmpty(t, hash)
	assert.Nil(t, err)
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	assert.Nil(t, err)
}

func TestHashAndSalt_Error(t *testing.T) {
	password := "passwordisreallylongandxcv dsfvdsfvsfvsfzvzfddvdfvdwon'tcvdbdgbdfbdgfbdtvdefgbngjsdevghmxdffchgxfg"
	_, err := HashAndSalt(password)
	assert.NotNil(t, err)
}

type mockUserStore struct {
	stores.UserStoreInterface
	CreateUserCalls int
}

func (m *mockUserStore) CreateUser(user *stores.User) ([]store.ModelValidationError, error) {
	m.CreateUserCalls++
	return nil, nil
}

func TestUserDomain_CreateUser_ValidationError(t *testing.T) {
	domain := UserDomain{
		store: &mockUserStore{},
	}
	validationErrors, userData, err := domain.CreateUser(UserForCreate{})
	assert.NotEmpty(t, validationErrors)
	assert.Nil(t, userData)
	assert.Nil(t, err)
	assert.Equal(t, domain.store.(*mockUserStore).CreateUserCalls, 0)
}

func TestUserDomain_CreateUser_UserGetError(t *testing.T) {
	domain := UserDomain{
		store: &MockUserStoreGetError{},
	}
	validationErrors, userData, err := domain.CreateUser(UserForCreate{
		Name:     "name",
		Email:    "email@email.com",
		Password: "password",
	})
	assert.Nil(t, validationErrors)
	assert.Nil(t, userData)
	assert.NotNil(t, err)
}

func (m *mockUserStore) GetUserByEmail(email string) (*stores.User, error) {
	return nil, nil
}

func TestUserDomain_CreateUser_ValidationErrors(t *testing.T) {
	domain := UserDomain{
		store: &mockUserStore{},
	}
	validationErrors, userData, err := domain.CreateUser(UserForCreate{})
	assert.NotEmpty(t, validationErrors)
	assert.Nil(t, userData)
	assert.Nil(t, err)
	assert.Equal(t, domain.store.(*mockUserStore).CreateUserCalls, 0)
}

func TestUserDomain_CreateUser_UserExists(t *testing.T) {
	domain := UserDomain{
		store: &MockUserStoreUserExists{},
	}
	validationErrors, userData, err := domain.CreateUser(UserForCreate{
		Name:     "name",
		Email:    "email@email.com",
		Password: "password",
	})
	assert.NotEmpty(t, validationErrors)
	assert.Nil(t, userData)
	assert.Nil(t, err)
	assert.Equal(t, validationErrors[0].Field, "email")
	assert.Equal(t, validationErrors[0].Errors[0], "Email is already in use.")
}

func TestUserDomain_CreateUser_HashError(t *testing.T) {
	domain := UserDomain{
		store: &mockUserStore{},
	}
	validationErrors, userData, err := domain.CreateUser(UserForCreate{
		Name:     "name",
		Email:    "email@email.com",
		Password: "passwordisreallylongandxcv dsfvdsfvsfvsfzvzfddvdfvdwon'tcvdbdgbdfbdgfbdtvdefgbngjsdevghmxdffchgxfg",
	})
	assert.Nil(t, validationErrors)
	assert.Nil(t, userData)
	assert.NotNil(t, err)
	assert.Equal(t, domain.store.(*mockUserStore).CreateUserCalls, 0)
}

type mockUserStoreWithError struct {
	stores.UserStoreInterface
}

func (m *mockUserStoreWithError) CreateUser(user *stores.User) ([]store.ModelValidationError, error) {
	return nil, assert.AnError
}

func (m *mockUserStoreWithError) GetUserByEmail(email string) (*stores.User, error) {
	return nil, nil
}

func (m *mockUserStoreWithError) GetUserByID(email uint) (*stores.User, error) {
	return nil, nil
}

func TestUserDomain_CreateUser_StoreError(t *testing.T) {
	domain := UserDomain{
		store: &mockUserStoreWithError{},
	}
	validationErrors, userData, err := domain.CreateUser(UserForCreate{
		Name:     "name",
		Email:    "some@email.com",
		Password: "password",
	})
	assert.Nil(t, validationErrors)
	assert.Nil(t, userData)
	assert.NotNil(t, err)
}

type mockUserStoreSuccess struct {
	stores.UserStoreInterface
}

func (m *mockUserStoreSuccess) CreateUser(user *stores.User) ([]store.ModelValidationError, error) {
	user.ID = 1
	return nil, nil
}

func (m *mockUserStoreSuccess) GetUserByEmail(email string) (*stores.User, error) {
	return nil, nil

}

func TestUserDomain_CreateUser_Success(t *testing.T) {
	domain := UserDomain{
		store: &mockUserStoreSuccess{},
	}
	validationErrors, userData, err := domain.CreateUser(UserForCreate{
		Name:     "name",
		Email:    "some@email.com",
		Password: "password",
	})
	assert.Nil(t, validationErrors)
	assert.NotNil(t, userData)
	assert.Nil(t, err)
	assert.Equal(t, userData.ID, uint(1))
}

func TestUserDomain_CreateUser_Integration(t *testing.T) {
	domain := UserDomain{
		store: stores.NewUserStore(),
	}
	store.IsolatedIntegrationTest(t, []store.IntegrationTestStore{domain.store}, func(t *testing.T) {
		validationErrors, userData, err := domain.CreateUser(UserForCreate{
			Name:     "name",
			Email:    "some@email.com",
			Password: "password",
		})
		assert.Nil(t, validationErrors)
		assert.NotNil(t, userData)
		assert.Nil(t, err)
		assert.NotEqual(t, userData.ID, uint(0))
		dbUser, err := domain.store.GetUserByID(userData.ID)
		assert.Nil(t, err)
		assert.Equal(t, userData.ID, dbUser.ID)
		assert.Equal(t, userData.Name, dbUser.Name)
		assert.Equal(t, userData.Email, dbUser.Email)
	})
}

func TestUserDomain_GetUserByEmail_NoUser(t *testing.T) {
	domain := UserDomain{
		store: &MockUserStoreUserDoesNotExist{},
	}
	userData, err := domain.GetUserByEmail("email")
	assert.Nil(t, userData)
	assert.Nil(t, err)
}

func TestUserDomain_GetUserByEmail_Error(t *testing.T) {
	domain := UserDomain{
		store: &MockUserStoreError{},
	}
	userData, err := domain.GetUserByEmail("email")
	assert.Nil(t, userData)
	assert.NotNil(t, err)
}

func TestUserDomain_GetUserByEmail_Success(t *testing.T) {
	domain := UserDomain{
		store: &MockUserStoreUserExists{},
	}
	userData, err := domain.GetUserByEmail("email")
	assert.NotNil(t, userData)
	assert.Nil(t, err)
	assert.Equal(t, userData.Email, "got@byemail.com")
	assert.Equal(t, userData.Name, "Test User")
}

func TestUserDomain_GetUserByEmail_Integration(t *testing.T) {
	domain := UserDomain{
		store: stores.NewUserStore(),
	}
	store.IsolatedIntegrationTest(t, []store.IntegrationTestStore{domain.store}, func(t *testing.T) {
		user := stores.User{
			Name:     "Test User",
			Email:    "test@email.com",
			Password: "password",
		}
		validationErrors, err := domain.store.CreateUser(&user)
		assert.Nil(t, err)
		assert.Empty(t, validationErrors)
		assert.NotEqual(t, uint(0), user.ID)
		userData, err := domain.GetUserByEmail(user.Email)
		assert.Nil(t, err)
		assert.NotNil(t, userData)
		assert.Equal(t, user.ID, userData.ID)
		assert.Equal(t, user.Name, userData.Name)
		assert.Equal(t, user.Email, userData.Email)
	})
}

func TestUserDomain_GetUserByID_NoUser(t *testing.T) {
	domain := UserDomain{
		store: &MockUserStoreUserDoesNotExist{},
	}
	userData, err := domain.GetUserByID(1)
	assert.Nil(t, userData)
	assert.Nil(t, err)
}

func TestUserDomain_GetUserByID_Error(t *testing.T) {
	domain := UserDomain{
		store: &MockUserStoreError{},
	}
	userData, err := domain.GetUserByID(1)
	assert.Nil(t, userData)
	assert.NotNil(t, err)
}

func TestUserDomain_GetUserByID_Success(t *testing.T) {
	domain := UserDomain{
		store: &MockUserStoreUserExists{},
	}
	userData, err := domain.GetUserByID(1)
	assert.NotNil(t, userData)
	assert.Nil(t, err)

	assert.Equal(t, userData.Email, "got@byid.com")
	assert.Equal(t, userData.Name, "Test User")
}

func TestUserDomain_GetUserByID_Integration(t *testing.T) {
	domain := UserDomain{
		store: stores.NewUserStore(),
	}
	store.IsolatedIntegrationTest(t, []store.IntegrationTestStore{domain.store}, func(t *testing.T) {
		user := stores.User{
			Name:     "Test User",
			Email:    "auser@thing.com",
			Password: "password",
		}
		validationErrors, err := domain.store.CreateUser(&user)
		assert.Nil(t, err)
		assert.Empty(t, validationErrors)
		assert.NotEqual(t, uint(0), user.ID)
		userData, err := domain.GetUserByID(user.ID)
		assert.Nil(t, err)
		assert.NotNil(t, userData)
		assert.Equal(t, user.ID, userData.ID)
		assert.Equal(t, user.Name, userData.Name)
		assert.Equal(t, user.Email, userData.Email)
	})
}

func TestUserDomain_DeleteUser_NoUser(t *testing.T) {
	domain := UserDomain{
		store: &MockUserStoreUserDoesNotExist{},
	}
	err := domain.DeleteUser(1)
	assert.Nil(t, err)
}

func TestUserDomain_DeleteUser_Error(t *testing.T) {
	domain := UserDomain{
		store: &MockUserStoreError{},
	}
	err := domain.DeleteUser(1)
	assert.NotNil(t, err)
}

func TestUserDomain_DeleteUser_Success(t *testing.T) {
	domain := UserDomain{
		store: &MockUserStoreUserExists{},
	}
	err := domain.DeleteUser(1)
	assert.Nil(t, err)
}

func TestUserDomain_DeleteUser_Integration(t *testing.T) {
	domain := UserDomain{
		store: stores.NewUserStore(),
	}
	store.IsolatedIntegrationTest(t, []store.IntegrationTestStore{domain.store}, func(t *testing.T) {
		user := stores.User{
			Name:     "Test User",
			Email:    "some@user.com",
			Password: "password",
		}
		validationErrors, err := domain.store.CreateUser(&user)
		assert.Nil(t, err)
		assert.Empty(t, validationErrors)
		assert.NotEqual(t, uint(0), user.ID)
		err = domain.DeleteUser(user.ID)
		assert.Nil(t, err)
		res, err := domain.GetUserByID(user.ID)
		assert.Nil(t, err)
		assert.Nil(t, res)
	})
}

type UserStoreHashedPassword struct {
	stores.UserStoreInterface
}

func (u *UserStoreHashedPassword) GetUserByID(id uint) (*stores.User, error) {
	password, err := HashAndSalt("password")
	if err != nil {
		return nil, err
	}
	return &stores.User{
		Name:     "Test User",
		Email:    "test@email.com",
		Password: password,
	}, nil
}

func TestUserDomain_CheckPassword(t *testing.T) {
	domain := UserDomain{
		store: &UserStoreHashedPassword{},
	}
	assert.True(t, domain.CheckPassword(1, "password"))
	assert.False(t, domain.CheckPassword(1, "wrongpassword"))
}

func TestUserDomain_CheckPassword_Error(t *testing.T) {
	domain := UserDomain{
		store: &MockUserStoreError{},
	}
	match := domain.CheckPassword(1, "password")
	assert.False(t, match)
}

func TestUserDomain_GenerateJWT_Success(t *testing.T) {
	domain := UserDomain{
		Config: UserDomainConfig{
			SecretKey: []byte("averysecretkey"),
		},
	}
	jwt, err := domain.GenerateJWT("email")
	assert.Nil(t, err)
	assert.NotNil(t, jwt)
	assert.Equal(t, jwt.Email, "email")
	assert.NotNil(t, jwt.ExpiresAt)
	assert.NotNil(t, jwt.TokenString)
}

func TestUserDomain_ValidateJWT_Invalid(t *testing.T) {
	domain := UserDomain{
		Config: UserDomainConfig{
			SecretKey: []byte("averysecretkey"),
		},
	}
	jwtString := "invalidstring"
	errors, claims := domain.ValidateJWT(jwtString)
	assert.Contains(t, errors.Message, "Bearer token invalid")
	assert.Nil(t, claims)
}

func TestUserDomain_ValidateJWT_Expired(t *testing.T) {
	domain := UserDomain{
		Config: UserDomainConfig{
			SecretKey:                  []byte("averysecretkey"),
			TokenExpirationTimeMinutes: 0,
		},
	}
	jwt, err := domain.GenerateJWT("email")
	assert.Nil(t, err)
	assert.NotNil(t, jwt)
	errors, claims := domain.ValidateJWT(jwt.TokenString)
	assert.Contains(t, errors.Message, "Bearer token invalid")
	assert.Nil(t, claims)
}

func TestUserDomain_ValidateJWT_Success(t *testing.T) {
	domain := UserDomain{
		Config: UserDomainConfig{
			SecretKey:                  []byte("averysecret"),
			TokenExpirationTimeMinutes: 15,
		},
	}
	jwt, err := domain.GenerateJWT("email")
	assert.Nil(t, err)
	assert.NotNil(t, jwt)
	errors, claims := domain.ValidateJWT(jwt.TokenString)
	assert.Nil(t, errors)
	assert.NotNil(t, claims)
	assert.Equal(t, claims.Email, "email")
}

func TestUserDomain_GetJWTFromRequest_NoCookie(t *testing.T) {
	req := http.Request{}
	domain := UserDomain{}
	token, err := domain.GetJWTFromRequest(req)
	assert.Empty(t, token)
	assert.NotNil(t, err)
}

func TestUserDomain_GetJWTFromRequest_Header(t *testing.T) {
	req := http.Request{
		Header: map[string][]string{
			"Authorization": {"Token token"},
		},
	}
	domain := UserDomain{}
	token, err := domain.GetJWTFromRequest(req)
	assert.Equal(t, token, "Token token")
	assert.Nil(t, err)
}

func TestUserDomain_ValidateRequestAuth_NoToken(t *testing.T) {
	req := http.Request{}
	domain := UserDomain{}
	errors, claims := domain.ValidateRequestAuth(req)
	assert.NotNil(t, errors)
	assert.Nil(t, claims)
}

func TestUserDomain_ValidateRequestAuth_InvalidToken(t *testing.T) {
	req := http.Request{
		Header: map[string][]string{
			"Authorization": {"Token sometoken"},
		},
	}
	domain := UserDomain{
		Config: UserDomainConfig{
			SecretKey: []byte("averysecret"),
		},
	}
	errors, claims := domain.ValidateRequestAuth(req)
	assert.Contains(t, errors.Message, "Bearer token invalid")
	assert.Nil(t, claims)
}

func TestUserDomain_ValidateRequestAuth_ValidToken(t *testing.T) {
	domain := UserDomain{
		Config: UserDomainConfig{
			SecretKey:                  []byte("averysecretkey"),
			TokenExpirationTimeMinutes: 15,
		},
	}
	jwt, err := domain.GenerateJWT("email")
	assert.Nil(t, err)
	assert.NotNil(t, jwt)
	req := http.Request{
		Header: map[string][]string{
			"Authorization": {"Bearer " + jwt.TokenString},
		},
	}
	errors, claims := domain.ValidateRequestAuth(req)
	assert.Nil(t, errors)
	assert.NotNil(t, claims)
	assert.Equal(t, claims.Email, "email")
}

type MockUserStoreMismatchedPassword struct {
	stores.UserStoreInterface
}

func (m *MockUserStoreMismatchedPassword) GetUserByID(id uint) (*stores.User, error) {
	password, err := HashAndSalt("password")
	if err != nil {
		return nil, err
	}
	return &stores.User{
		Name:     "Test User",
		Email:    "test@user.com",
		Password: password,
	}, nil
}

func TestUserDomain_ChangeUserPassword_PasswordDoesNotMatch(t *testing.T) {
	domain := UserDomain{
		store: &MockUserStoreMismatchedPassword{},
	}
	errors, err := domain.ChangeUserPassword(1, "sausage", "newpassword")
	assert.NotEmpty(t, errors)
	assert.Nil(t, err)
	assert.Equal(t, errors[0].Field, "old_password")
	assert.Equal(t, errors[0].Errors[0], "Old password is incorrect")
}

func TestUserDomain_ChangeUserPassword_HashError(t *testing.T) {
	domain := UserDomain{
		store: &MockUserStoreMismatchedPassword{},
	}
	errors, err := domain.ChangeUserPassword(1, "password", "passwordisreallylongandxcv dsfvdsfvsfvsfzvzfddvdfvdwon'tcvdbdgbdfbdgfbdtvdefgbngjsdevghmxdffchgxfg")
	assert.Nil(t, errors)
	assert.NotNil(t, err)
}

type MockUserStoreGetError struct {
	stores.UserStoreInterface
	GetByIdCallCount int
}

func (m *MockUserStoreGetError) GetUserByID(id uint) (*stores.User, error) {
	m.GetByIdCallCount++
	if m.GetByIdCallCount == 1 {
		hashedPassword, err := HashAndSalt("password")
		if err != nil {
			return nil, err
		}
		return &stores.User{
			Name:     "Test User",
			Email:    "test@email.com",
			Password: hashedPassword,
		}, nil
	}
	return nil, assert.AnError
}

func (m *MockUserStoreGetError) GetUserByEmail(email string) (*stores.User, error) {
	return nil, assert.AnError
}

func TestUserDomain_ChangeUserPassword_GetUserError(t *testing.T) {
	domain := UserDomain{
		store: &MockUserStoreGetError{},
	}
	errors, err := domain.ChangeUserPassword(1, "password", "newpassword")
	assert.Nil(t, errors)
	assert.NotNil(t, err)
}

type MockUserStoreUpdateError struct {
	stores.UserStoreInterface
}

func (m *MockUserStoreUpdateError) GetUserByID(id uint) (*stores.User, error) {
	hashedPassword, err := HashAndSalt("password")
	if err != nil {
		return nil, err
	}
	return &stores.User{
		Name:     "Test User",
		Email:    "test@email.com",
		Password: hashedPassword,
	}, nil
}

func (m *MockUserStoreUpdateError) UpdateUser(user *stores.User) ([]store.ModelValidationError, error) {
	return nil, assert.AnError
}

func TestUserDomain_ChangeUserPassword_UpdateError(t *testing.T) {
	domain := UserDomain{
		store: &MockUserStoreUpdateError{},
	}
	errors, err := domain.ChangeUserPassword(1, "password", "newpassword")
	assert.Nil(t, errors)
	assert.NotNil(t, err)
}

type mockUserStoreUpdateValidationErrors struct {
	MockUserStoreUpdateError
}

func (m *mockUserStoreUpdateValidationErrors) UpdateUser(user *stores.User) ([]store.ModelValidationError, error) {
	return []store.ModelValidationError{
		{
			Field:   "password",
			Message: "Password is required",
		},
	}, nil
}

func TestUserDomain_ChangeUserPassword_ValidationErrors(t *testing.T) {
	domain := UserDomain{
		store: &mockUserStoreUpdateValidationErrors{},
	}
	errors, err := domain.ChangeUserPassword(1, "password", "newpassword")
	assert.NotEmpty(t, errors)
	assert.Nil(t, err)
	assert.Contains(t, errors[0].Field, "password")
	assert.Contains(t, errors[0].Errors[0], "Password is required")
}

func TestUserDomain_ChangeUserPassword_Integration(t *testing.T) {
	domain := UserDomain{
		store: stores.NewUserStore(),
	}
	store.IsolatedIntegrationTest(t, []store.IntegrationTestStore{domain.store}, func(t *testing.T) {
		user := UserForCreate{
			Name:     "Test User",
			Email:    "test@email.com",
			Password: "password",
		}
		validationErrors, userData, err := domain.CreateUser(user)
		assert.Nil(t, err)
		assert.Empty(t, validationErrors)
		assert.NotNil(t, userData)

		errors, err := domain.ChangeUserPassword(userData.ID, "password", "newpassword")
		assert.Nil(t, errors)
		assert.Nil(t, err)

		matches := domain.CheckPassword(userData.ID, "newpassword")
		assert.True(t, matches)
	})
}
