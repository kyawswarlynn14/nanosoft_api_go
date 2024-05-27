package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"nanosoft/database"
	"nanosoft/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ServiceCollection *mongo.Collection = database.ServiceData(database.Client, "Services")

func CreateService() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		
		var service models.Service
		if err := c.BindJSON(&service); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		service.Service_ID = primitive.NewObjectID()
		service.Created_At = time.Now()
		service.Updated_At = time.Now()

		_, err := ServiceCollection.InsertOne(ctx, service)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating service"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Service created successfully"})
	}
}

func UpdateService() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(serviceID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var service models.Service
		if err := c.BindJSON(&service); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		service.Updated_At = time.Now()

		update := bson.M{
			"$set": bson.M{
				"title":       service.Title,
				"description": service.Description,
				"image":       service.Image,
				"image_path":  service.ImagePath,
				"t1":          service.T1,
				"t2":          service.T2,
				"updated_at":  service.Updated_At,
			},
		}

		result, err := ServiceCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating service", "details": err.Error()})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Service updated successfully"})
	}
}

func DeleteService() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(serviceID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := ServiceCollection.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting service"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
	}
}

func GetAllServices() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var services []models.Service
		cursor, err := ServiceCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving services"})
			return
		}

		if err = cursor.All(ctx, &services); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding services"})
			return
		}

		c.JSON(http.StatusOK, services)
	}
}

func GetOneService() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(serviceID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var service models.Service
		err = ServiceCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&service)
		if err != nil {
			// Log the error for debugging purposes
			log.Printf("Error retrieving service: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving service", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, service)
	}
}
