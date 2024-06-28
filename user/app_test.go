package user

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	linesHttp "lines/lines/http"
	"testing"
)

func TestNewUserApp(t *testing.T) {
	app := NewUserApp()
	if app.http == nil {
		t.Errorf("Expected app.http to not be nil")
	}
}

func TestUserApp_Initialise(t *testing.T) {
	app := NewUserApp()
	assert.Nil(t, app.Initialise())
}

type MockUserApp struct {
	UserApp
	RegisterHTTPRoutesCalled bool
	RegisterHTTPRoutesArgs   []linesHttp.HttpEngine
}

func (a *MockUserApp) RegisterHTTPRoutes(engine linesHttp.HttpEngine) {
	a.RegisterHTTPRoutesCalled = true
	a.RegisterHTTPRoutesArgs = append(a.RegisterHTTPRoutesArgs, engine)
}

func TestUserApp_RegisterHTTPRoutes(t *testing.T) {
	app := MockUserApp{}
	engine := gin.Engine{}
	app.RegisterHTTPRoutes(&engine)
	assert.True(t, app.RegisterHTTPRoutesCalled)
	assert.Equal(t, 1, len(app.RegisterHTTPRoutesArgs))
	assert.Equal(t, &engine, app.RegisterHTTPRoutesArgs[0])
}

func TestUserApp_RegisterGRPCServices(t *testing.T) {
	app := NewUserApp()
	assert.Nil(t, app.RegisterGRPCServices())
}
