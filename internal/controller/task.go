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

func (tc *TaskController) postTask(c *gin.Context) {
	var taskReq TaskRequest

	err := c.BindJSON(&taskReq)
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
	}

	err = tc.insertTask(taskReq, "")

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
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

func (tc *TaskController) getAllTasks(c *gin.Context) {

	res, err := database.MongoDB.GetTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
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
