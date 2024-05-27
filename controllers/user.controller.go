package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"nanosoft/database"
	"nanosoft/models"
	generate "nanosoft/tokens"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")

var Validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(hashedPassword string, plainPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Login or Password is Incorrect"
		valid = false
	}
	return valid, msg
}

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		}
		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()

		role := 0
		user.Role = role
		token, refreshtoken, _ := generate.TokenGenerator(*user.Email, *user.Name, user.User_ID, user.Role)
		user.Token = &token
		user.Refresh_Token = &refreshtoken
		_, inserterr := UserCollection.InsertOne(ctx, user)
		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusCreated, "Successfully Registered!!")
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var founduser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password incorrect"})
			return
		}
		PasswordIsValid, msg := VerifyPassword(*founduser.Password, *user.Password)
		defer cancel()
		if !PasswordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}
		token, refreshToken, _ := generate.TokenGenerator(*founduser.Email, *founduser.Name, founduser.User_ID, founduser.Role)
		defer cancel()
		generate.UpdateAllTokens(token, refreshToken, founduser.User_ID)

		c.JSON(http.StatusFound, founduser)
	}
}

func RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken := c.Query("refreshToken")
		if refreshToken == "" {
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Refresh Token"})
			c.Abort()
			return
		}

		claims, msg := generate.ValidateToken(refreshToken)
		if msg != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}

		newAccessToken, newRefreshToken, err := generate.TokenGenerator(claims.Email, claims.Name, claims.Uid, claims.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate new tokens"})
			return
		}
		generate.UpdateAllTokens(newAccessToken, newRefreshToken, claims.Uid)

		c.JSON(http.StatusOK, gin.H{
			"access_token":  newAccessToken,
			"refresh_token": newRefreshToken,
		})
	}
}

func GetUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		userEmail, ok := c.Get("email")
		if !ok {
			log.Println("email not found in context")
			c.IndentedJSON(http.StatusBadRequest, "email not found in context")
			return
		}
		emailStr, ok := userEmail.(string)
		if !ok {
			log.Println("email in context is not a string")
			c.IndentedJSON(http.StatusBadRequest, "email in context is not a string")
			return
		}
		log.Println("user email: " + emailStr)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var foundUser models.User
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "email", Value: emailStr}}).Decode(&foundUser)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "invalid token")
			return
		}

		c.JSON(http.StatusOK, foundUser)
	}
}

func UpdateUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userInfo struct {
			Name       string `json:"name"`
			Avatar     string `json:"avatar"`
			AvatarPath string `json:"avatar_path"`
		}

		// Log the incoming request body
		if err := c.ShouldBindJSON(&userInfo); err != nil {
			log.Println("Invalid request body:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		log.Println("Received user info:", userInfo)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userEmail, ok := c.Get("email")
		if !ok {
			log.Println("email not found in context")
			c.IndentedJSON(http.StatusBadRequest, "email not found in context")
			return
		}
		emailStr, ok := userEmail.(string)
		if !ok {
			log.Println("email in context is not a string")
			c.IndentedJSON(http.StatusBadRequest, "email in context is not a string")
			return
		}
		log.Println("Updating user with email:", emailStr)

		update := bson.M{
			"$set": bson.M{
				"name":        userInfo.Name,
				"avatar":      userInfo.Avatar,
				"avatar_path": userInfo.AvatarPath,
			},
		}

		result, err := UserCollection.UpdateOne(ctx, bson.M{"email": emailStr}, update)
		if err != nil {
			log.Println("Failed to update user info:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user info"})
			return
		}

		if result.MatchedCount == 0 {
			log.Println("User not found with email:", emailStr)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		var updatedUser models.User
		err = UserCollection.FindOne(ctx, bson.M{"email": emailStr}).Decode(&updatedUser)
		if err != nil {
			log.Println("Error retrieving updated user info:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving updated user info"})
			return
		}

		log.Println("Updated user info:", updatedUser)
		c.JSON(http.StatusOK, updatedUser)
	}
}

func UpdateUserPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var passwordUpdate struct {
			OldPassword string `json:"old_password"`
			NewPassword string `json:"new_password"`
		}

		if err := c.ShouldBindJSON(&passwordUpdate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		log.Println(passwordUpdate)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userEmail, ok := c.Get("email")
		if !ok {
			log.Println("email not found in context")
			c.IndentedJSON(http.StatusBadRequest, "email not found in context")
			return
		}
		emailStr, ok := userEmail.(string)
		if !ok {
			log.Println("email in context is not a string")
			c.IndentedJSON(http.StatusBadRequest, "email in context is not a string")
			return
		}

		var foundUser models.User
		err := UserCollection.FindOne(ctx, bson.M{"email": emailStr}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		log.Println(foundUser.Password)

		PasswordIsValid, msg := VerifyPassword(*foundUser.Password, passwordUpdate.OldPassword)
		if !PasswordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}

		newPasswordHash := HashPassword(passwordUpdate.NewPassword)
		foundUser.Password = &newPasswordHash

		_, err = UserCollection.UpdateOne(ctx, bson.M{"email": emailStr}, bson.M{
			"$set": bson.M{"password": newPasswordHash},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
	}
}

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var users []models.User
		cursor, err := UserCollection.Find(ctx, bson.M{})
		if err != nil {
			log.Println("Error finding users:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving users"})
			return
		}

		if err = cursor.All(ctx, &users); err != nil {
			log.Println("Error decoding users:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding users"})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func UpdateUserRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		var roleUpdate struct {
			UserID string `json:"user_id"`
			Role   int    `json:"role"`
		}

		if err := c.ShouldBindJSON(&roleUpdate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userID, err := primitive.ObjectIDFromHex(roleUpdate.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		update := bson.M{
			"$set": bson.M{
				"role": roleUpdate.Role,
			},
		}

		result, err := UserCollection.UpdateOne(ctx, bson.M{"_id": userID}, update)
		if err != nil {
			log.Println("Failed to update user role:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully"})
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		result, err := UserCollection.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			log.Println("Failed to delete user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}
