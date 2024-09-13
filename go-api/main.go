package main

import (
	"log"
	"net/http"
	"os"
	"taskmanagement/handlers"
	"taskmanagement/utils"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// MongoDB and Redis initialization
	mongoURI := "mongodb+srv://reddyashok9:h3nRSFTWcRXY1juL@mydemo.piofkut.mongodb.net/?retryWrites=true&w=majority&appName=mydemo/"
	redisURI := os.Getenv("REDIS_URI")

	mongoCollection, err := utils.ConnectMongo(mongoURI)
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	redisClient := utils.ConnectRedis(redisURI)

	// TaskHandler setup
	taskHandler := handlers.TaskHandler{
		MongoCollection: mongoCollection,
		RedisClient:     redisClient,
	}

	// Router setup
	router := mux.NewRouter()
	router.HandleFunc("/tasks", taskHandler.GetTasks).Methods("GET") // Route to Get All Tasks
	router.HandleFunc("/tasks/bulk-create", taskHandler.BulkCreateTasks).Methods("POST")
	router.HandleFunc("/tasks/bulk-update", taskHandler.BulkUpdateTasks).Methods("PUT")

	// Setup CORS options
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3002"}, // Replace with your frontend origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	// Wrap the router with the CORS middleware
	handler := corsHandler.Handler(router)

	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", handler)
}
