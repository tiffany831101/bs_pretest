package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestStartServer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := StartServer()

	assert.NotNil(t, s)
	assert.NotNil(t, s.engine)
}

func TestServer_Run(t *testing.T) {

	gin.SetMode(gin.TestMode)
	s := StartServer()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/ping", nil)
	s.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())

}

func TestServer_SetUpRoutes(t *testing.T) {

	gin.SetMode(gin.TestMode)
	s := StartServer()
	s.SetUpRoutes()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/v1/tasks/", nil)
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

}
