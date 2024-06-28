package http

import (
	"github.com/stretchr/testify/assert"
	"lines/internal"
	"testing"
)

func TestCreateEngine(t *testing.T) {
	config := &internal.MainConfig{
		CORSOrigins: []string{"http://localhost:3000"},
	}
	engine := CreateEngine(config)
	assert.NotNil(t, engine)
}
