package http

import (
	"lines/lines/http"
	"lines/lines/utils"
	user_domain "lines/user/domain"
)

type UserHttpIngressInterface interface {
	RegisterRoutes(e http.HttpEngine) error
}

type UserHttpConfig struct {
	SiteDomain string
	UseSSL     bool
}

func NewUserHttpConfig() UserHttpConfig {
	return UserHttpConfig{
		SiteDomain: utils.GetEnvOrDefault("SITE_DOMAIN", "localhost", "string").(string),
		UseSSL:     !utils.GetEnvOrDefault("LOCAL_DEV", "true", "bool").(bool) && !utils.GetEnvOrDefault("ALLOW_HTTP", "false", "bool").(bool),
	}
}

type UserHttpIngress struct {
	config UserHttpConfig
	domain user_domain.UserDomainInterface
}

func NewUserHttpIngress() UserHttpIngress {
	return UserHttpIngress{
		config: NewUserHttpConfig(),
		domain: user_domain.NewUserDomain(),
	}
}
