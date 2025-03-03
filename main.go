package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client         *mongo.Client
	userCollection *mongo.Collection
	jwtKey         = []byte(os.Getenv("JWT_SECRET"))
)

func connectDB() *mongo.Client {
	log.Println(os.Getenv("MONGO_URI_LOCAL"))
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI_LOCAL"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
	userCollection = client.Database("arkademy").Collection("users")
	return client
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func generateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString(jwtKey)
}

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
		CreationTime: primitive.NewDateTimeFromTime(time.Now().UTC()),
		Characters:   make([]CharacterRecord, 0),
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

func createCharacterForUser(c *gin.Context) {
	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	var characterRecord CharacterRecord
	if err := c.ShouldBindJSON(&characterRecord); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := userData.(User)
	characterRecord.CreationTime = primitive.NewDateTimeFromTime(time.Now().UTC())
	if user.PlayerRecord.Characters == nil {
		user.PlayerRecord.Characters = make([]CharacterRecord, 0)
	}

	user.PlayerRecord.Characters = append(user.PlayerRecord.Characters, characterRecord)
	_, err := userCollection.UpdateOne(context.TODO(), bson.D{{"_id", user.ID}}, bson.D{{"$set", user}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
func main() {
	log.Println("Server begin")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading Environment")
	}
	client = connectDB()
	defer client.Disconnect(context.Background())
	router := gin.Default()

	router.POST("/register", registerUser)
	router.POST("/login", loginUser)

	protected := router.Group("/")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/user", getUser)
		protected.POST("/createCharacter", createCharacterForUser)
	}
	router.Run(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}
