package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client         *mongo.Client
	userCollection *mongo.Collection
	jwtKey         = []byte(os.Getenv("JWT_SECRET"))
)

func ConnectDB() *mongo.Client {
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

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func RegisterUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.Password = HashPassword(user.Password)
	_, err := userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading Environment")
	}
	client = ConnectDB()
	defer client.Disconnect(context.Background())
	router := gin.Default()

	router.POST("/register", RegisterUser)
	router.Run(":8080")
}
