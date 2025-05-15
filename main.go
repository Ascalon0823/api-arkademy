package main

import (
	"context"
	"fmt"
	"log"
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

func connectDB() *mongo.Client {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
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

func main() {
	log.Println("Server begin")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading Environment")
	}
	client = connectDB()
	defer client.Disconnect(context.Background())
	router := gin.Default()
	router.Use(CORSMiddleware())
	router.POST("/register", registerUser)
	router.POST("/login", loginUser)

	protected := router.Group("/")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/user", getUser)
		protected.PATCH("/player", updatePlayer)
	}
	router.Run(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}
