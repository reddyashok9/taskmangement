package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"taskmanagement/models"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskHandler struct {
	MongoCollection *mongo.Collection
	RedisClient     *redis.Client
}

// Response struct to include tasks and cached flag
type GetTasksResponse struct {
	Tasks  []models.Task `json:"tasks"`
	Cached bool          `json:"cached"`
}

// GetTasks fetches all tasks, first checking Redis cache, then falling back to MongoDB
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Check if the list of tasks is cached in Redis
	tasksJSON, err := h.RedisClient.Get(ctx, "tasks").Result()
	if err == redis.Nil {
		// If tasks are not cached, fetch from MongoDB
		cursor, err := h.MongoCollection.Find(ctx, bson.M{})
		if err != nil {
			http.Error(w, "Error fetching tasks from MongoDB", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		var tasks []models.Task
		if err := cursor.All(ctx, &tasks); err != nil {
			http.Error(w, "Error decoding tasks from MongoDB", http.StatusInternalServerError)
			return
		}

		// Marshal the tasks to JSON
		tasksJSON, err := json.Marshal(tasks)
		if err != nil {
			http.Error(w, "Error marshalling tasks data", http.StatusInternalServerError)
			return
		}

		// Store the tasks in Redis with an expiration time of 1 Minute
		h.RedisClient.Set(ctx, "tasks", tasksJSON, 1*time.Minute)

		// Return the tasks with cached flag set to false
		response := GetTasksResponse{
			Tasks:  tasks,
			Cached: false, // Data is from MongoDB
		}

		// Send the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else if err != nil {
		// If there was some other Redis error
		http.Error(w, "Error fetching tasks from cache", http.StatusInternalServerError)
		return
	} else {
		var tasks []models.Task
		err = json.Unmarshal([]byte(tasksJSON), &tasks)
		if err != nil {
			http.Error(w, "Error unmarshalling tasks from cache", http.StatusInternalServerError)
			return
		}

		response := GetTasksResponse{
			Tasks:  tasks,
			Cached: true, // Data is from Redis
		}

		// Send the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func (h *TaskHandler) BulkCreateTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	err := json.NewDecoder(r.Body).Decode(&tasks)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var docs []interface{}
	for _, task := range tasks {
		docs = append(docs, task)
	}

	_, err = h.MongoCollection.InsertMany(context.Background(), docs)
	if err != nil {
		http.Error(w, "Failed to create tasks", http.StatusInternalServerError)
		return
	}

	h.RedisClient.Del(context.Background(), "tasks") // Clear cached tasks
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Tasks created successfully"})
}

func (h *TaskHandler) BulkUpdateTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	err := json.NewDecoder(r.Body).Decode(&tasks)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for _, task := range tasks {
		filter := bson.M{"_id": task.ID}
		update := bson.M{"$set": task}
		opts := options.Update().SetUpsert(true)
		_, err := h.MongoCollection.UpdateOne(context.Background(), filter, update, opts)
		if err != nil {
			log.Println("Failed to update task:", task.ID, err)
		}
	}

	h.RedisClient.Del(context.Background(), "tasks") // Clear cached tasks
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Tasks updated successfully"})
}
