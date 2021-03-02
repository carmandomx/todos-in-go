package controllers

import (
	"context"
	"errors"
	"net/http"

	"github.com/carmandomx/todos-in-go.git/formatter"
	"github.com/carmandomx/todos-in-go.git/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var col *mongo.Collection = models.ConnectDB()

// CreateTodo ...
func CreateTodo(c *gin.Context) {
	var todo models.Todo

	id := primitive.NewObjectID()

	if err := c.ShouldBindJSON(&todo); err != nil {
		var verr validator.ValidationErrors

		if errors.As(err, &verr) {
			c.JSON(http.StatusBadRequest, gin.H{"errors": formatter.NewJSONFormatter().Simple(verr)})
			return
		}

		c.Error(err)
		return
	}

	todo.ID = id

	_, err := col.InsertOne(context.TODO(), todo)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Error saving to DB, please try again",
		})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

// GetTodos...
func GetTodos(c *gin.Context) {
	opts := options.Find()

	opts.SetLimit(50)

	cursor, err := col.Find(context.TODO(), bson.D{}, opts)

	if err != nil {
		c.Error(err)
		return
	}

	defer cursor.Close(context.TODO())

	var count int
	result := []models.Todo{}

	for cursor.Next(context.TODO()) {
		var helper models.Todo
		err := cursor.Decode(&helper)
		result = append(result, helper)
		if err != nil {
			c.Error(err)
		}

		count++
	}

	c.JSON(http.StatusOK, gin.H{
		"count": count,
		"todos": result,
	})

}
