package http

import (
	"github.com/gin-gonic/gin"
	linesHttp "lines/lines/http"
	"net/http"
	"time"
)

func (h *UserHttpIngress) V1SignIn(c *gin.Context) {
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

	user, err := h.domain.GetUserByEmail(credentials.Email)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, linesHttp.HttpError{Message: []string{"Credentials not recognised."}})
		return
	}

	passwordCorrect := h.domain.CheckPassword(user.ID, credentials.Password)
	if !passwordCorrect {
		c.JSON(http.StatusUnauthorized, linesHttp.HttpError{Message: []string{"Credentials not recognised."}})
		return
	}

	jwt, err := h.domain.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, linesHttp.HttpError{Message: []string{"Could not generate JWT."}})
	}

	expiresAt := int(jwt.ExpiresAt.Sub(time.Now()).Seconds())
	c.SetCookie(
		"Bearer",
		jwt.TokenString,
		expiresAt,
		"/",
		h.config.SiteDomain,
		h.config.UseSSL,
		true,
	)

	c.JSON(http.StatusOK, jwt)
}

//// UserSignOutAPI is the handler for user sign out, it clears the JWT cookie.
//func (h *UserHttpIngress) UserSignOutAPI(c *gin.Context) {
//	c.SetCookie(
//		"Bearer",
//		"",
//		0,
//		"/",
//		h.Config.SiteDomain,
//		h.Config.UseSSL,
//		true,
//	)
//}
//
//// UserRefreshTokenAPI is the handler for refreshing a JWT token.
//func (h *UserHttpIngress) UserRefreshTokenAPI(c *gin.Context) {
//	authResponse := h.AuthHandler.AuthenticateRequest(c.Request)
//	statusCode, errorString := h.AuthHandler.CheckIsAuthed(authResponse, "Bearer")
//	if statusCode != 0 {
//		c.JSON(
//			statusCode,
//			gin.H{"error": errorString},
//		)
//		return
//	}
//
//	expirationTime := time.Now().Add(
//		time.Duration(h.Config.TokenExpirationTimeMinutes) * time.Minute,
//	)
//	claims := &Claims{
//		Email: authResponse.User.Email,
//		RegisteredClaims: jwt.RegisteredClaims{
//			ExpiresAt: jwt.NewNumericDate(expirationTime),
//		},
//	}
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	tokenString, err := token.SignedString(h.Config.SecretKey)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate token."})
//	}
//
//	c.SetCookie(
//		"Bearer",
//		tokenString,
//		h.Config.TokenExpirationTimeMinutes*60,
//		"/",
//		h.Config.SiteDomain,
//		h.Config.UseSSL,
//		true,
//	)
//}
//
//// UserGetJWTAPI is the handler for getting a JWT token.
//// Token is used for authing the dashboard websocket.
//func (h *UserHttpIngress) UserGetJWTAPI(c *gin.Context) {
//	authResponse := h.AuthHandler.AuthenticateRequest(c.Request)
//	statusCode, errorString := h.AuthHandler.CheckIsAuthed(authResponse, "Bearer")
//	if statusCode != 0 {
//		c.JSON(
//			statusCode,
//			gin.H{"error": errorString},
//		)
//		return
//	}
//
//	expirationTime := time.Now().Add(
//		time.Duration(h.Config.TokenExpirationTimeMinutes) * time.Minute,
//	)
//	claims := &Claims{
//		Email: authResponse.User.Email,
//		RegisteredClaims: jwt.RegisteredClaims{
//			ExpiresAt: jwt.NewNumericDate(expirationTime),
//		},
//	}
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	tokenString, err := token.SignedString(h.Config.SecretKey)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate token."})
//	}
//
//	c.JSON(http.StatusOK, gin.H{"token": tokenString})
//}
