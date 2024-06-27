package domain

import (
	"github.com/stretchr/testify/assert"
	"lines/lines/store"
	"testing"
)

func TestNewUserDomainConfig(t *testing.T) {
	config := NewUserDomainConfig()
	assert.NotNil(t, config)
	assert.NotEmpty(t, config.SecretKey)
}

func TestNewUserDomain(t *testing.T) {
	domain := NewUserDomain()
	assert.NotNil(t, domain)
	assert.NotNil(t, domain.store)
	assert.NotEmpty(t, domain.Config.SecretKey)
}

func TestUserDomain_Integration(t *testing.T) {
	domain := NewUserDomain()
	store.IsolatedIntegrationTest(t, []store.IntegrationTestStore{domain.store}, func(t *testing.T) {
		validationErrors, user, err := domain.CreateUser(
			UserForCreate{Name: "Jake", Email: "some@email.com", Password: "password"},
		)
		assert.Empty(t, validationErrors)
		assert.NotNil(t, user)
		assert.Nil(t, err)
	})
}
