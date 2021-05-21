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
	todo.IsCompleted = false
	_, err := col.InsertOne(context.TODO(), todo)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Error saving to DB, please try again",
		})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

// GetTodos ...
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

// DeleteTodo ...
func DeleteTodo(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Invalid ID",
		})
		return
	}
	res := col.FindOneAndDelete(context.TODO(), bson.M{
		"_id": objectID,
	})

	if res.Err() != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Error while deleting, try again later",
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// UpdateTodo ..
func UpdateTodo(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Invalid ID",
		})
		return
	}
	var updatedTodo models.Todo
	if err := c.ShouldBindJSON(&updatedTodo); err != nil {
		var verr validator.ValidationErrors

		if errors.As(err, &verr) {
			c.JSON(http.StatusBadRequest, gin.H{"errors": formatter.NewJSONFormatter().Simple(verr)})
			return
		}

		c.Error(err)
		return
	}
	opts := options.FindOneAndUpdate()
	filter := bson.M{"_id": bson.M{"$eq": objectID}}
	update := bson.M{"$set": bson.M{"isCompleted": updatedTodo.IsCompleted}}
	opts.SetReturnDocument(options.After)
	res := col.FindOneAndUpdate(context.TODO(), filter, update, opts)

	err = res.Decode(&updatedTodo)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Error saving to DB, please try again",
		})
		return
	}

	c.JSON(http.StatusOK, updatedTodo)

}
