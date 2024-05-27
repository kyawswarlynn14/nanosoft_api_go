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

var ProjectCollection *mongo.Collection = database.ProjectData(database.Client, "Projects")

func CreateProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		
		var project models.Project
		if err := c.BindJSON(&project); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		project.Project_ID = primitive.NewObjectID()
		project.Created_At = time.Now()
		project.Updated_At = time.Now()

		_, err := ProjectCollection.InsertOne(ctx, project)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating project"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Project created successfully"})
	}
}

func UpdateProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(projectID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var project models.Project
		if err := c.BindJSON(&project); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		project.Updated_At = time.Now()

		update := bson.M{
			"$set": bson.M{
				"title":       	project.Title,
				"description": 	project.Description,
				"demo_link": 	project.DemoLink,
				"tech":  		project.Tech,
				"images":       project.Images,
				"t1":          	project.T1,
				"t2":          	project.T2,
				"updated_at":  	project.Updated_At,
			},
		}

		result, err := ProjectCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating project", "details": err.Error()})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Project updated successfully"})
	}
}

func DeleteProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(projectID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := ProjectCollection.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting project"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
	}
}

func GetAllProjects() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var services []models.Project
		cursor, err := ProjectCollection.Find(ctx, bson.M{})
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

func GetOneProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(projectID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var project models.Project
		err = ProjectCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&project)
		if err != nil {
			// Log the error for debugging purposes
			log.Printf("Error retrieving project: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving project", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, project)
	}
}
