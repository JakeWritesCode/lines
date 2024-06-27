package http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	linesHttp "lines/lines/http"
	"net/http"
	"net/http/httptest"
	"testing"
)

func EndpointIsAuthenticatedTest(t *testing.T, endpoint func(c *gin.Context)) {
	req, err := http.NewRequest("GET", "/some-ep", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/some-ep", endpoint)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	response := linesHttp.HttpError{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Contains(t, response.Message, "Unauthorised")
}
