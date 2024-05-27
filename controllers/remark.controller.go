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

var RemarkCollection *mongo.Collection = database.RemarkData(database.Client, "Remarks")

func CreateRemark() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var remark models.Remark
		if err := c.BindJSON(&remark); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		remark.Remark_ID = primitive.NewObjectID()
		remark.Created_At = time.Now()
		remark.Updated_At = time.Now()

		_, err := RemarkCollection.InsertOne(ctx, remark)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating remark"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Remark created successfully"})
	}
}

func UpdateRemark() gin.HandlerFunc {
	return func(c *gin.Context) {
		remarkID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(remarkID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid remark ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var remark models.Remark
		if err := c.BindJSON(&remark); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		remark.Updated_At = time.Now()

		update := bson.M{
			"$set": bson.M{
				"name":       remark.Name,
				"role":       remark.Role,
				"image":      remark.Image,
				"image_path": remark.ImagePath,
				"remark":     remark.Remark,
				"t1":         remark.T1,
				"t2":         remark.T2,
				"updated_at": remark.Updated_At,
			},
		}

		result, err := RemarkCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating remark", "details": err.Error()})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Remark not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Remark updated successfully"})
	}
}

func DeleteRemark() gin.HandlerFunc {
	return func(c *gin.Context) {
		remarkID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(remarkID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid remark ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := RemarkCollection.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting remark"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Remark not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Remark deleted successfully"})
	}
}

func GetAllRemarks() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var services []models.Remark
		cursor, err := RemarkCollection.Find(ctx, bson.M{})
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

func GetOneRemark() gin.HandlerFunc {
	return func(c *gin.Context) {
		remarkID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(remarkID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid remark ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var remark models.Remark
		err = RemarkCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&remark)
		if err != nil {
			// Log the error for debugging purposes
			log.Printf("Error retrieving remark: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving remark", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, remark)
	}
}
