package http

import (
	"github.com/gin-gonic/gin"
	linesHttp "lines/lines/http"
	"net/http"
	"time"
)

func (i *UserHttpIngress) V1SignIn(c *gin.Context) {
	var credentials UserLogin
	err := c.BindJSON(&credentials)
	if err != nil {
		httpErr := linesHttp.HttpError{Message: []string{err.Error()}}
		c.JSON(http.StatusBadRequest, httpErr)
		return
	}
	if httpErr := credentials.Validate(); len(httpErr.Message) > 0 {
		c.JSON(http.StatusBadRequest, httpErr)
		return
	}

	user, err := i.domain.GetUserByEmail(credentials.Email)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, linesHttp.HttpError{Message: []string{"Credentials not recognised."}})
		return
	}

	passwordCorrect := i.domain.CheckPassword(user.ID, credentials.Password)
	if !passwordCorrect {
		c.JSON(http.StatusUnauthorized, linesHttp.HttpError{Message: []string{"Credentials not recognised."}})
		return
	}

	jwt, err := i.domain.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, linesHttp.HttpError{Message: []string{"Could not generate JWT."}})
	}

	expiresAt := int(jwt.ExpiresAt.Sub(time.Now()).Seconds())
	c.SetCookie(
		"Bearer",
		jwt.TokenString,
		expiresAt,
		"/",
		i.config.SiteDomain,
		i.config.UseSSL,
		true,
	)

	c.JSON(http.StatusOK, jwt)
}

// UserSignOutAPI is the handler for user sign out, it clears the JWT cookie.
func (i *UserHttpIngress) V1SignOut(c *gin.Context) {
	c.SetCookie(
		"Bearer",
		"",
		0,
		"/",
		i.config.SiteDomain,
		i.config.UseSSL,
		true,
	)
}

// UserRefreshTokenAPI is the handler for refreshing a JWT token.
func (i *UserHttpIngress) V1RefreshToken(c *gin.Context) {
	authError, claims := i.domain.ValidateRequestAuth(*c.Request)
	if authError != nil {
		c.JSON(http.StatusUnauthorized, authError)
		return
	}
	newClaims, err := i.domain.GenerateJWT(claims.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, linesHttp.HttpError{Message: []string{"Could not generate JWT."}})
		return
	}

	c.SetCookie(
		"Bearer",
		newClaims.TokenString,
		int(newClaims.ExpiresAt.Sub(time.Now()).Seconds()),
		"/",
		i.config.SiteDomain,
		i.config.UseSSL,
		true,
	)

	c.JSON(http.StatusOK, newClaims)
}
