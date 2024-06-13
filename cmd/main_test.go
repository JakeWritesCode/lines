package main

import (
	"github.com/stretchr/testify/assert"
	"lines/internal"
	"lines/lines/app"
	"lines/lines/http"
	"testing"
)

type mockApp struct {
	InitialiseCalls           int
	InitialiseArgs            []*internal.MainConfig
	RegisterHttpRoutesCalls   int
	RegisterGRPCServicesCalls int
}

func (m *mockApp) Initialise(config *internal.MainConfig) error {
	m.InitialiseCalls++
	m.InitialiseArgs = append(m.InitialiseArgs, config)
	return nil
}

func (m *mockApp) RegisterHTTPRoutes(engine http.HttpEngine) {
	m.RegisterHttpRoutesCalls++
}

func (m *mockApp) RegisterGRPCServices() error {
	m.RegisterGRPCServicesCalls++
	return nil
}

type mockHttpEngine struct {
	RunCalls int
}

func (m *mockHttpEngine) Run(addr ...string) error {
	m.RunCalls++
	return nil
}

func TestMainHandler_InitialisesApps(t *testing.T) {
	apps := []app.App{
		&mockApp{},
		&mockApp{},
		&mockApp{},
	}
	config := &internal.MainConfig{}
	httpEngine := &mockHttpEngine{}

	MainHandler(apps, config, httpEngine)

	for _, a := range apps {
		mockApp := a.(*mockApp)
		assert.Equal(t, 1, mockApp.InitialiseCalls)
		assert.Equal(t, 1, mockApp.RegisterHttpRoutesCalls)
		assert.Equal(t, 0, mockApp.RegisterGRPCServicesCalls)
		assert.Equal(t, config, mockApp.InitialiseArgs[0])
	}
}

func TestMainHandler_StartsHttpEngine(t *testing.T) {
	apps := []app.App{
		&mockApp{},
	}
	config := &internal.MainConfig{}
	httpEngine := &mockHttpEngine{}

	MainHandler(apps, config, httpEngine)

	assert.Equal(t, 1, httpEngine.RunCalls)
}

type mockAppWithInitialiseError struct {
	*mockApp
}

func (m *mockAppWithInitialiseError) Initialise(config *internal.MainConfig) error {
	return assert.AnError
}

type MockLogger struct {
	*internal.LogrusHandler
	FatalCalls int
}

func (m *MockLogger) Fatal(appName string, caller string, message string) {
	m.FatalCalls++
}

func TestMainHandler_InitialiseError_LogsError(t *testing.T) {

	apps := []app.App{
		&mockAppWithInitialiseError{},
	}
	config := &internal.MainConfig{}
	config.Logger = &MockLogger{}
	httpEngine := &mockHttpEngine{}
	defer func() {
		if r := recover(); r != nil {
			assert.True(t, true, "MainHandler did panic as expected")
			assert.Equal(t, config.Logger.(*MockLogger).FatalCalls, 1, "Logger.Fatal was called")
		}
	}()

	MainHandler(apps, config, httpEngine)

	assert.Equal(t, 1, config.Logger.(*MockLogger).FatalCalls)
}

type mockHttpEngineWithRunError struct {
	*mockHttpEngine
}

func (m *mockHttpEngineWithRunError) Run(addr ...string) error {
	return assert.AnError
}

func TestMainHandler_RunError_LogsError(t *testing.T) {
	apps := []app.App{
		&mockApp{},
	}
	config := &internal.MainConfig{}
	config.Logger = &MockLogger{}
	httpEngine := &mockHttpEngineWithRunError{}
	defer func() {
		if r := recover(); r != nil {
			assert.True(t, true, "MainHandler did panic as expected")
			assert.Equal(t, config.Logger.(*MockLogger).FatalCalls, 1, "Logger.Fatal was called")
		}
	}()

	MainHandler(apps, config, httpEngine)

	assert.Equal(t, 1, config.Logger.(*MockLogger).FatalCalls)
}
