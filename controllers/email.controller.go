package controllers

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"nanosoft/database"
	"nanosoft/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/gomail.v2"
)

var EmailCollection *mongo.Collection = database.EmailData(database.Client, "Emails")

func SendEmail(subject string, body string) error {
	SMIP_HOST := os.Getenv("SMIP_HOST")
	SMIP_PORT, portErr := strconv.Atoi(os.Getenv("SMIP_PORT"))
	SMIP_MAIL := os.Getenv("SMIP_MAIL")
	SMIP_PASSWORD := os.Getenv("SMIP_PASSWORD")
	SMIP_RECEPT_MAIL := os.Getenv("SMIP_RECEPT_MAIL")

	if portErr != nil {
		log.Printf("Error converting SMTP_PORT to integer: %v", portErr)
		return portErr
	}

    m := gomail.NewMessage()
    m.SetHeader("From", SMIP_MAIL)
    m.SetHeader("To", SMIP_RECEPT_MAIL)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", body)

    d := gomail.NewDialer(SMIP_HOST, SMIP_PORT, SMIP_MAIL, SMIP_PASSWORD)
    d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    if err := d.DialAndSend(m); err != nil {
        return err
    }
    return nil
}

func CreateEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		
		var message models.Message
		if err := c.BindJSON(&message); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		message.Message_ID = primitive.NewObjectID()
		message.Created_At = time.Now()
		message.Updated_At = time.Now()

		subject := "Email From Client"
		emailBody := `
			<h1>New Message from Client</h1>
			<p><strong>Name:</strong> ` + *message.Name + `</p>
			<p><strong>Email:</strong> ` + *message.Email + `</p>
			<p><strong>Phone:</strong> ` + *message.Phone + `</p>
			<p><strong>Company Name:</strong> ` + *message.CompanyName + `</p>
			<p><strong>Message:</strong> ` + *message.Message + `</p>
		`

		emailErr := SendEmail(subject, emailBody)
        if emailErr != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending email", "details": emailErr.Error()})
            return
        }

		_, err := EmailCollection.InsertOne(ctx, message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating message"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Message sent successfully"})
	}
}

func DeleteEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		messageID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(messageID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := EmailCollection.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting message"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
	}
}

func GetAllEmails() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var messages []models.Message
		cursor, err := EmailCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving services"})
			return
		}

		if err = cursor.All(ctx, &messages); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding messages"})
			return
		}

		c.JSON(http.StatusOK, messages)
	}
}

func GetOneEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		messageID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(messageID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var message models.Message
		err = EmailCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&message)
		if err != nil {
			// Log the error for debugging purposes
			log.Printf("Error retrieving message: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving message", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, message)
	}
}
