package http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	domain2 "lines/lines/domain"
	linesHttp "lines/lines/http"
	"lines/lines/store"
	"lines/user/domain"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestUserHttpIngress_V1SignIn_FailedBind(t *testing.T) {
	ingress := UserHttpIngress{}
	req, err := http.NewRequest("GET", "/sign-in", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/sign-in", ingress.V1SignIn)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUserHttpIngress_V1SignIn_ValidationErrors(t *testing.T) {
	ingress := UserHttpIngress{}
	bodyCreds, err := json.Marshal(UserLogin{})
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/sign-in", strings.NewReader(string(bodyCreds)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/sign-in", ingress.V1SignIn)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	errors := linesHttp.HttpError{}
	err = json.Unmarshal(rr.Body.Bytes(), &errors)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, errors.Message, "Email is required.")
	assert.Contains(t, errors.Message, "Password is required.")
}

type mockUserDomain struct {
	domain.UserDomainInterface
}

func (m *mockUserDomain) GetUserByEmail(email string) (*domain.UserData, error) {
	return nil, nil
}

func TestUserHttpIngress_V1SignIn_UserNotFound(t *testing.T) {
	ingress := UserHttpIngress{
		domain: &mockUserDomain{},
	}
	bodyCreds, err := json.Marshal(UserLogin{
		Email:    "email",
		Password: "password",
	})
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/sign-in", strings.NewReader(string(bodyCreds)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/sign-in", ingress.V1SignIn)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	errors := linesHttp.HttpError{}
	err = json.Unmarshal(rr.Body.Bytes(), &errors)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, errors.Message, "Credentials not recognised.")
}

type mockUserDomainPasswordDoesNotMatch struct {
	domain.UserDomainInterface
}

func (m *mockUserDomainPasswordDoesNotMatch) GetUserByEmail(email string) (*domain.UserData, error) {
	return &domain.UserData{
		Email: "email",
	}, nil
}

func (m *mockUserDomainPasswordDoesNotMatch) CheckPassword(userID uint, password string) bool {
	return false
}

func TestUserHttpIngress_V1SignIn_PasswordDoesNotMatch(t *testing.T) {
	ingress := UserHttpIngress{
		domain: &mockUserDomainPasswordDoesNotMatch{},
	}
	bodyCreds, err := json.Marshal(UserLogin{
		Email:    "email",
		Password: "password",
	})
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/sign-in", strings.NewReader(string(bodyCreds)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/sign-in", ingress.V1SignIn)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	errors := linesHttp.HttpError{}
	err = json.Unmarshal(rr.Body.Bytes(), &errors)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, errors.Message, "Credentials not recognised.")
}

type mockUserDomainGenerateJWTError struct {
	domain.UserDomainInterface
}

func (m *mockUserDomainGenerateJWTError) GetUserByEmail(email string) (*domain.UserData, error) {
	return &domain.UserData{
		Email: "email",
	}, nil
}

func (m *mockUserDomainGenerateJWTError) CheckPassword(userID uint, password string) bool {
	return true
}

func (m *mockUserDomainGenerateJWTError) GenerateJWT(email string) (*domain.JWTClaimsOut, error) {
	return nil, assert.AnError
}

func TestUserHttpIngress_V1SignIn_GenerateJWTError(t *testing.T) {
	ingress := UserHttpIngress{
		domain: &mockUserDomainGenerateJWTError{},
	}
	bodyCreds, err := json.Marshal(UserLogin{
		Email:    "email",
		Password: "password",
	})
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/sign-in", strings.NewReader(string(bodyCreds)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/sign-in", ingress.V1SignIn)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	errors := linesHttp.HttpError{}
	err = json.Unmarshal(rr.Body.Bytes(), &errors)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, errors.Message, "Could not generate JWT.")
}

type mockUserDomainSuccess struct {
	mockUserDomainGenerateJWTError
}

func (m *mockUserDomainSuccess) GenerateJWT(email string) (*domain.JWTClaimsOut, error) {
	return &domain.JWTClaimsOut{
		TokenString: "token",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
	}, nil
}

func TestUserHttpIngress_V1SignIn_Success(t *testing.T) {
	ingress := UserHttpIngress{
		domain: &mockUserDomainSuccess{},
	}
	bodyCreds, err := json.Marshal(UserLogin{
		Email:    "email",
		Password: "password",
	})
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/sign-in", strings.NewReader(string(bodyCreds)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/sign-in", ingress.V1SignIn)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	jwt := domain.JWTClaimsOut{}
	err = json.Unmarshal(rr.Body.Bytes(), &jwt)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "token", jwt.TokenString)
	assert.Equal(t, time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local), jwt.ExpiresAt.Time)
}

func TestUserHttpIngress_V1SignIn_SetsCookie(t *testing.T) {
	ingress := UserHttpIngress{
		domain: &mockUserDomainSuccess{},
		config: UserHttpConfig{
			SiteDomain: "example.com",
			UseSSL:     true,
		},
	}
	bodyCreds, err := json.Marshal(UserLogin{
		Email:    "email",
		Password: "password",
	})
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/sign-in", strings.NewReader(string(bodyCreds)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/sign-in", ingress.V1SignIn)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Bearer=token; Path=/; Domain=example.com; Max-Age=0; HttpOnly; Secure", rr.Header().Get("Set-Cookie"))
}

func TestUserHttpIngress_V1SignIn_Integration(t *testing.T) {

	ingress := NewUserHttpIngress(nil)
	store.IsolatedIntegrationTest(t, []store.IntegrationTestStore{ingress.domain}, func(t *testing.T) {
		user := domain.UserForCreate{
			Name:     "name",
			Email:    "email@email.com",
			Password: "password",
		}
		validationErrors, _, err := ingress.domain.CreateUser(user)
		assert.Nil(t, validationErrors)
		assert.Nil(t, err)

		bodyCreds, err := json.Marshal(UserLogin{
			Email:    user.Email,
			Password: user.Password,
		})
		assert.Nil(t, err)
		req, err := http.NewRequest("POST", "/sign-in", strings.NewReader(string(bodyCreds)))
		assert.Nil(t, err)

		rr := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/sign-in", ingress.V1SignIn)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		jwt := domain.JWTClaimsOut{}
		err = json.Unmarshal(rr.Body.Bytes(), &jwt)
		assert.Nil(t, err)
		assert.NotEmpty(t, jwt.TokenString)
		assert.NotEqual(t, time.Time{}, jwt.ExpiresAt.Time)
		assert.Equal(t, "Bearer="+jwt.TokenString+"; Path=/; Domain=localhost; Max-Age=899; HttpOnly", rr.Header().Get("Set-Cookie"))
	})
}

func TestUserHttpIngress_V1SignOut(t *testing.T) {
	ingress := UserHttpIngress{}
	req, err := http.NewRequest("GET", "/sign-out", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/sign-out", ingress.V1SignOut)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Bearer=; Path=/; HttpOnly", rr.Header().Get("Set-Cookie"))
}

func TestUserHttpIngress_V1RefreshToken_Unauthenticated(t *testing.T) {
	ingress := UserHttpIngress{
		domain: &domain.UserDomain{},
	}
	EndpointIsAuthenticatedTest(t, ingress.V1RefreshToken)
}

func TestUserHttpIngress_V1RefreshToken_Success(t *testing.T) {
	ingress := UserHttpIngress{
		domain: &domain.UserDomain{
			Config: domain.UserDomainConfig{TokenExpirationTimeMinutes: 10},
		},
	}
	token, err := ingress.domain.GenerateJWT("email")
	// Change expiry to change token string
	ingress = UserHttpIngress{
		domain: &domain.UserDomain{
			Config: domain.UserDomainConfig{TokenExpirationTimeMinutes: 5},
		},
	}
	assert.Nil(t, err)
	req, err := http.NewRequest("GET", "/refresh-token", nil)
	assert.Nil(t, err)
	req.AddCookie(&http.Cookie{
		Name:  "Bearer",
		Value: token.TokenString,
	})

	rr := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/refresh-token", ingress.V1RefreshToken)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	jwt := domain.JWTClaimsOut{}
	err = json.Unmarshal(rr.Body.Bytes(), &jwt)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, jwt.TokenString)
	assert.NotEqual(t, token.TokenString, jwt.TokenString)
}

func TestUserHttpIngress_V1SignUp_BadRequest(t *testing.T) {
	ingress := UserHttpIngress{}
	req, err := http.NewRequest("POST", "/sign-up", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/sign-up", ingress.V1SignUp)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUserHttpIngress_V1SignUp_HttpValidationErrors(t *testing.T) {
	ingress := UserHttpIngress{}
	bodyCreds, err := json.Marshal(UserSignUp{})
	assert.Nil(t, err)
	req, err := http.NewRequest("POST", "/sign-up", strings.NewReader(string(bodyCreds)))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/sign-up", ingress.V1SignUp)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	errors := linesHttp.HttpError{}
	err = json.Unmarshal(rr.Body.Bytes(), &errors)
	assert.Nil(t, err)
	assert.Contains(t, errors.Message, "Name is required.")
	assert.Contains(t, errors.Message, "Email is required.")
	assert.Contains(t, errors.Message, "Password is required.")
}

type mockUserDomainUserCreateError struct {
	domain.UserDomainInterface
}

func (m *mockUserDomainUserCreateError) CreateUser(user domain.UserForCreate) ([]domain2.DomainValidationErrors, *domain.UserData, error) {
	return []domain2.DomainValidationErrors{}, nil, assert.AnError
}

func TestUserHttpIngress_V1SignUp_UserCreateError(t *testing.T) {
	ingress := UserHttpIngress{
		domain: &mockUserDomainUserCreateError{},
	}
	bodyCreds, err := json.Marshal(UserSignUp{
		Name:     "name",
		Email:    "email",
		Password: "password",
	})
	assert.Nil(t, err)
	req, err := http.NewRequest("POST", "/sign-up", strings.NewReader(string(bodyCreds)))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/sign-up", ingress.V1SignUp)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	errors := linesHttp.HttpError{}
	err = json.Unmarshal(rr.Body.Bytes(), &errors)
	assert.Nil(t, err)
	assert.Contains(t, errors.Message, "Could not create user.")
}

type mockUserDomainUserCreateValidationErrors struct {
	domain.UserDomainInterface
}

func (m *mockUserDomainUserCreateValidationErrors) CreateUser(user domain.UserForCreate) ([]domain2.DomainValidationErrors, *domain.UserData, error) {
	return []domain2.DomainValidationErrors{{"name", []string{"error"}}}, nil, nil
}

func TestUserHttpIngress_V1SignUp_DomainValidationErrors(t *testing.T) {
	ingress := UserHttpIngress{
		domain: &mockUserDomainUserCreateValidationErrors{},
	}
	bodyCreds, err := json.Marshal(UserSignUp{
		Name:     "name",
		Email:    "email",
		Password: "password",
	})
	assert.Nil(t, err)
	req, err := http.NewRequest("POST", "/sign-up", strings.NewReader(string(bodyCreds)))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/sign-up", ingress.V1SignUp)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var errors []domain2.DomainValidationErrors
	err = json.Unmarshal(rr.Body.Bytes(), &errors)
	assert.Nil(t, err)
	assert.Contains(t, errors[0].Errors, "error")
}

type mockUserDomainUserCreateSuccess struct {
	domain.UserDomainInterface
}

func (m *mockUserDomainUserCreateSuccess) CreateUser(user domain.UserForCreate) ([]domain2.DomainValidationErrors, *domain.UserData, error) {
	return []domain2.DomainValidationErrors{}, &domain.UserData{
		ID:    1,
		Name:  "name",
		Email: "email",
	}, nil
}

func TestUserHttpIngress_V1SignUp_Success(t *testing.T) {
	ingress := UserHttpIngress{
		domain: &mockUserDomainUserCreateSuccess{},
	}
	bodyCreds, err := json.Marshal(UserSignUp{
		Name:     "name",
		Email:    "email",
		Password: "password",
	})
	assert.Nil(t, err)
	req, err := http.NewRequest("POST", "/sign-up", strings.NewReader(string(bodyCreds)))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/sign-up", ingress.V1SignUp)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	user := UserReadDTO{}
	err = json.Unmarshal(rr.Body.Bytes(), &user)
	assert.Nil(t, err)
	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "name", user.Name)
	assert.Equal(t, "email", user.Email)
}

func TestUserHttpIngress_V1SignUp_Integration(t *testing.T) {
	ingress := NewUserHttpIngress(nil)
	store.IsolatedIntegrationTest(t, []store.IntegrationTestStore{ingress.domain}, func(t *testing.T) {
		body, err := json.Marshal(UserSignUp{
			Name:     "name",
			Email:    "email@email.com",
			Password: "password",
		})
		assert.Nil(t, err)
		req, err := http.NewRequest("POST", "/sign-up", strings.NewReader(string(body)))
		assert.Nil(t, err)

		rr := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/sign-up", ingress.V1SignUp)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		user := UserReadDTO{}
		err = json.Unmarshal(rr.Body.Bytes(), &user)
		assert.Nil(t, err)
		assert.Equal(t, "name", user.Name)
		assert.Equal(t, "email@email.com", user.Email)
		assert.NotEqual(t, uint(0), user.ID)
	})
}
