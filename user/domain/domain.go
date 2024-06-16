package domain

import (
	"lines/lines/utils"
	"lines/user/stores"
)

// UserDomainConfig is a struct that contains the configuration for a UserDomain.
type UserDomainConfig struct {
	SECRET_KEY string
}

func NewUserDomainConfig() UserDomainConfig {
	return UserDomainConfig{
		SECRET_KEY: utils.GetEnvOrDefault("SECRET_KEY", "secret_key", "string").(string),
	}
}

type UserDomain struct {
	store  stores.UserStoreInterface
	config UserDomainConfig
}

// NewUserDomain is a function that returns a new UserDomain instance.
func NewUserDomain() *UserDomain {
	return &UserDomain{
		store:  stores.NewUserStore(),
		config: NewUserDomainConfig(),
	}
}
