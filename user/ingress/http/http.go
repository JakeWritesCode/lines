package http

import (
	"lines/lines/http"
	"lines/lines/utils"
	user_domain "lines/user/domain"
)

type UserHttpIngressInterface interface {
	RegisterRoutes(e http.HttpEngine)
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

func (i *UserHttpIngress) RegisterRoutes(e http.HttpEngine) {
	e.POST("/users/sign-in", i.V1SignIn)
	e.POST("/users/sign-out", i.V1SignOut)
	e.GET("/users/refresh-token", i.V1RefreshToken)
	e.POST("/users/sign-up", i.V1SignUp)
	e.GET("/users/me", i.V1GetUser)
}

func NewUserHttpIngress(domain user_domain.UserDomainInterface) UserHttpIngress {
	if domain == nil {
		domain = user_domain.NewUserDomain()
	}
	return UserHttpIngress{
		config: NewUserHttpConfig(),
		domain: domain,
	}
}
