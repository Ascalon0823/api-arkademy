package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func loginUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var dbUser User
	err := userCollection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&dbUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid username or password"})
		return
	}
	if hashPassword(user.Password) != dbUser.Password {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid username or password"})
		return
	}
	token, err := generateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func registerUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.Password = hashPassword(user.Password)
	_, err := userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	user.PlayerRecord = PlayerRecord{
		Characters: make([]CharacterRecord, 0),
	}
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func getUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user.(User))
}

func updatePlayer(c *gin.Context) {
	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	var playerRecord PlayerRecord
	if err := c.ShouldBindJSON(&playerRecord); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := userData.(User)
	user.PlayerRecord = playerRecord
	_, err := userCollection.UpdateOne(context.TODO(), bson.D{{"_id", user.ID}}, bson.D{{"$set", user}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.Password = ""
	c.JSON(http.StatusOK, user)
}
