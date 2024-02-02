package controller

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tiffany831101/bs_pretest.git/internal/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockDB struct {
}

func (db *MockDB) CloseConnection() {
	fmt.Println("Close Connection")
}

func (db *MockDB) InsertSingleTask(task database.Task) error {
	return nil
}

func (db *MockDB) GetTaskByID(taskID string) (database.Task, error) {
	database.MongoDB = &MockDB{}
	return database.Task{}, nil
}

func (db *MockDB) GetTasks() ([]database.Task, error) {
	database.MongoDB = &MockDB{}
	return []database.Task{}, nil
}

func (db *MockDB) DeleteTaskByID(taskID string) (int64, error) {
	database.MongoDB = &MockDB{}
	return 1, nil
}

func (db *MockDB) UpdateTaskID(taskID string, task database.Task) error {
	database.MongoDB = &MockDB{}
	return nil
}

func Test_NewTaskController(t *testing.T) {
	tC := &TaskController{}
	assert.NotNil(t, tC)

}
func Test_InsertTask(t *testing.T) {
	database.MongoDB = &MockDB{}

	tC := &TaskController{}
	var status TaskStatus = 0

	err := tC.insertTask(TaskRequest{
		Name:   "Test Task",
		Status: &status,
	}, "")

	assert.Equal(t, nil, err)
}

func Test_GetAllTasks(t *testing.T) {
	database.MongoDB = &MockDB{}

	tC := &TaskController{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	tC.getAllTasks(c)
	assert.Equal(t, http.StatusOK, w.Code)

}

func Test_GetTaskByID(t *testing.T) {
	database.MongoDB = &MockDB{}

	tC := &TaskController{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	taskID := primitive.NewObjectID().Hex()

	c.Params = append(c.Params, gin.Param{Key: "id", Value: taskID})

	tC.getTaskByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_DeleteTaskByID(t *testing.T) {

	database.MongoDB = &MockDB{}

	tC := &TaskController{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	taskID := primitive.NewObjectID().Hex()

	c.Params = append(c.Params, gin.Param{Key: "id", Value: taskID})

	tC.deleteTask(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_PostTask(t *testing.T) {

	database.MongoDB = &MockDB{}

	requestBody := `{"name": "Test Task", "status": 0}`

	w := httptest.NewRecorder()

	tC := &TaskController{}
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: map[string][]string{"Content-Type": {"application/json"}},
		Body:   io.NopCloser(strings.NewReader(requestBody)),
	}

	tC.postTask(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func Test_PutTask(t *testing.T) {

	database.MongoDB = &MockDB{}

	requestBody := `{"name": "Test Task", "status": 0}`

	w := httptest.NewRecorder()

	tC := &TaskController{}
	c, _ := gin.CreateTestContext(w)
	taskID := primitive.NewObjectID().Hex()

	c.Params = append(c.Params, gin.Param{Key: "id", Value: taskID})

	c.Request = &http.Request{
		Header: map[string][]string{"Content-Type": {"application/json"}},
		Body:   io.NopCloser(strings.NewReader(requestBody)),
	}

	tC.putTask(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}
