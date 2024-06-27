package domain

import (
	"lines/lines/utils"
	"lines/user/stores"
)

// UserDomainConfig is a struct that contains the configuration for a UserDomain.
type UserDomainConfig struct {
	SecretKey                  []byte
	TokenExpirationTimeMinutes int
}

func NewUserDomainConfig() UserDomainConfig {
	return UserDomainConfig{
		SecretKey:                  []byte(utils.GetEnvOrDefault("SECRET_KEY", "secret_key", "string").(string)),
		TokenExpirationTimeMinutes: utils.GetEnvOrDefault("TOKEN_EXPIRATION_TIME_MINUTES", "15", "int").(int),
	}
}

type UserDomain struct {
	store  stores.UserStoreInterface
	Config UserDomainConfig
}

func (d *UserDomain) BeginTransaction() error {
	return d.store.BeginTransaction()
}

func (d *UserDomain) RollbackTransaction() error {
	return d.store.RollbackTransaction()
}

// NewUserDomain is a function that returns a new UserDomain instance.
func NewUserDomain() *UserDomain {
	return &UserDomain{
		store:  stores.NewUserStore(),
		Config: NewUserDomainConfig(),
	}
}
