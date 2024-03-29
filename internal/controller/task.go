package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tiffany831101/bs_pretest.git/internal/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskController struct{}

type TaskStatus int

const (
	Incomplete TaskStatus = 0
	Completed  TaskStatus = 1
)

type TaskRequest struct {
	Name   string      `json:"name" binding:"required"`
	Status *TaskStatus `json:"status" binding:"required"`
}

type TaskResponse struct {
	ID     string
	Name   string
	Status int
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var tC *TaskController

func SetUpTasksRoutes(r *gin.Engine) {

	taskGroup := r.Group("/api/v1/tasks")
	{

		taskGroup.GET("/", tC.getAllTasks)
		taskGroup.GET("/:id", tC.getTaskByID)
		taskGroup.PUT("/:id", tC.putTask)

		taskGroup.POST("/", tC.postTask)
		taskGroup.DELETE("/:id", tC.deleteTask)

	}
}

func NewTasksController() {
	tC = &TaskController{}
}

// postTask creates a new task.
// @Summary Create a new task
// @Description Create a new task with the provided details.
// @ID postTask
// @Accept json
// @Produce json
// @Param body body TaskRequest true "Task details to create"
// @Success 201 {string} string "Created"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /tasks [post]
// @Tags tasks
func (tc *TaskController) postTask(c *gin.Context) {
	var taskReq TaskRequest

	err := c.BindJSON(&taskReq)
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
	}

	err = tc.insertTask(taskReq, "")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, "Created")

}

func (tc *TaskController) insertTask(task TaskRequest, taskID string) error {

	dbTask := database.Task{
		Name:   task.Name,
		Status: int(*task.Status),
	}
	if taskID != "" {
		objectID, _ := primitive.ObjectIDFromHex(taskID)

		dbTask.ID = objectID
	}

	err := database.MongoDB.InsertSingleTask(dbTask)

	if err != nil {
		return err
	}

	return nil
}

// getAllTasks retrieves all tasks.
// @Summary Retrieve all tasks
// @Description Get details of all tasks.
// @ID getAllTasks
// @Accept json
// @Produce json
// @Success 200 {array} TaskResponse "OK"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /tasks [get]
// @Tags tasks
func (tc *TaskController) getAllTasks(c *gin.Context) {

	res, err := database.MongoDB.GetTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	results := []TaskResponse{}
	for _, t := range res {
		task := TaskResponse{
			ID:     t.ID.Hex(),
			Name:   t.Name,
			Status: t.Status,
		}

		results = append(results, task)
	}

	c.JSON(http.StatusOK, results)
}

// getTaskByID retrieves a task by ID.
// @Summary Retrieve a task by ID
// @Description Get details of an existing task by ID.
// @ID getTaskByID
// @Accept json
// @Produce json
// @Param id path string true "ID of the task to retrieve" Pattern("^[0-9a-fA-F]{24}$")
// @Success 200 {object} TaskResponse "OK"
// @Success 404 {object} ErrorResponse "Resource Not Found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /tasks/{id} [get]
// @Tags tasks
func (tc *TaskController) getTaskByID(c *gin.Context) {

	taskID := c.Param("id")

	_, err := primitive.ObjectIDFromHex(taskID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource Not Found"})
		return
	}

	res, err := database.MongoDB.GetTaskByID(taskID)

	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, res)

}

// putTask updates or creates a task.
// @Summary Update a task
// @Description Update an existing task or create a new one if not exists.
// @ID updateTask
// @Accept json
// @Produce json
// @Param id path string true "ID of the task to update" Pattern("^[0-9a-fA-F]{24}$")
// @Param body body TaskRequest true "Task details to update"
// @Success 200 {string} string "OK"
// @Success 201 {string} string "Created"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /tasks/{id} [put]
// @Tags tasks
func (tc *TaskController) putTask(c *gin.Context) {

	var taskReq TaskRequest

	err := c.BindJSON(&taskReq)
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskID := c.Param("id")

	_, err = primitive.ObjectIDFromHex(taskID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Task ID, should be in hex format"})
		return
	}

	task, _ := database.MongoDB.GetTaskByID(taskID)

	if task == (database.Task{}) {

		err = tc.insertTask(taskReq, taskID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, "Created")

	} else {

		err = database.MongoDB.UpdateTaskID(taskID, database.Task(database.Task{
			Name:   taskReq.Name,
			Status: int(*taskReq.Status),
		}))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, "OK")

	}

}

// deleteTask deletes a task by ID.
// @Summary Delete a task
// @Description Delete an existing task by ID.
// @ID deleteTask
// @Accept json
// @Produce json
// @Param id path string true "ID of the task to delete" Pattern("^[0-9a-fA-F]{24}$")
// @Success 200 {string} string "OK"
// @Success 404 {object} ErrorResponse "Resource Not Found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /tasks/{id} [delete]
// @Tags tasks
func (tc *TaskController) deleteTask(c *gin.Context) {
	taskID := c.Param("id")

	_, err := primitive.ObjectIDFromHex(taskID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource Not Found"})
		return
	}

	deleteCount, err := database.MongoDB.DeleteTaskByID(taskID)

	if err != nil {

		if deleteCount == 0 {

			c.JSON(http.StatusNotFound, gin.H{"error": "Resource Not Found"})

		} else {

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		}
		return
	}

	if deleteCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource Not Found"})
		return
	}

	c.JSON(http.StatusOK, "OK")
}
