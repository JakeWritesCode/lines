package domain

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"lines/lines/store"
	"lines/user/stores"
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
