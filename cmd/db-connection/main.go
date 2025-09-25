package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"go-learning/internal/config"
	"go-learning/internal/db"
	"go-learning/internal/db/models"

	"github.com/google/uuid"
)

func main() {
	log.Println("Starting database example application...")

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Load configuration
	cfg := config.LoadConfig()

	// Create database connection
	conn, err := db.NewConnection(cfg.DB)
	if err != nil {
		log.Fatal("Failed to create database connection:", err)
	}

	// Start the connection
	if err := conn.Start(); err != nil {
		log.Fatal("Failed to start database connection:", err)
	}
	defer conn.Stop()

	// Create database layer
	dbLayer, err := db.NewDb(db.DBConfig{
		Db: conn,
	})
	if err != nil {
		log.Fatal("Failed to create database layer:", err)
	}

	ctx := context.Background()

	// Create a user with random email
	randomEmail := fmt.Sprintf("user%d@example.com", rand.Intn(100000))
	createReq := &models.CreateUserRequest{
		ID:    uuid.NewString(),
		Name:  "John Doe",
		Email: randomEmail,
	}

	log.Printf("Creating user with email: %s", randomEmail)
	createResp, err := dbLayer.CreateUser(ctx, createReq)
	if err != nil {
		log.Fatal("Failed to create user:", err)
	}

	log.Printf("Created user with ID: %s", createResp.ID)

	// Get the user
	getReq := &models.GetUserRequest{
		ID: createResp.ID,
	}

	log.Println("Retrieving user...")
	getResp, err := dbLayer.GetUser(ctx, getReq)
	if err != nil {
		log.Fatal("Failed to get user:", err)
	}

	log.Printf("Retrieved user: %+v", getResp.User)

	// List users
	listReq := &models.ListUsersRequest{
		Limit:  10,
		Offset: 0,
	}

	log.Println("Listing users...")
	listResp, err := dbLayer.ListUsers(ctx, listReq)
	if err != nil {
		log.Fatal("Failed to list users:", err)
	}

	log.Printf("Found %d users", len(listResp.Users))
	log.Println("Database example completed successfully!")
}
