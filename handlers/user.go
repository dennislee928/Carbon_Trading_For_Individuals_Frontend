package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

// Function to register a new user and insert into Supabase
func Register(c *gin.Context) {
	// Retrieve user data from the request body
	var user struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Validate input
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Initialize Supabase client
	client, err := NewSupabaseClient()
	if err != nil {
		log.Println("Error creating Supabase client:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error initializing Supabase client"})
		return
	}

	// Create the data to insert
	userData := map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
		"password": user.Password, // NOTE: In production, hash passwords before storing them!
	}

	// Insert user into the "users" table
	result, count, err := client.From("users").Insert(userData, false, "","", "*").Execute()

	// Handle errors
	if err != nil {
		log.Println("Error inserting user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
		return
	}
	// Check if no rows were affected
if count == 0 {
    log.Println("No rows were inserted")
    c.JSON(http.StatusInternalServerError, gin.H{"error": "No rows affected"})
    return
}

	// Check if no rows were affected
	if len(result) == 0 {
		log.Println("No rows were inserted")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No rows affected"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"result":  result,
	})}


// NewSupabaseClient creates and returns a new Supabase client
// NewSupabaseClient creates and returns a new Supabase client
func NewSupabaseClient() (*supabase.Client, error) {
	// Supabase project URL and API Key
	supabaseUrl := "https://omlzzhqhuhcpypohelbq.supabase.co"
	supabaseKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Im9tbHp6aHFodWhjcHlwb2hlbGJxIiwicm9sZSI6ImFub24iLCJpYXQiOjE3MzEwMzAzMzUsImV4cCI6MjA0NjYwNjMzNX0.fUKcu_mfWqESfZlAHQWjFPsvA0tHGq6t_fnTtPQhL-Q"

	client, err := supabase.NewClient(supabaseUrl, supabaseKey, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}
//
func ViewProfile(c *gin.Context) {
	c.JSON(200, gin.H{"message": "ViewProfile handler is not implemented yet"})
}

func UpdateProfile(c *gin.Context) {
	c.JSON(200, gin.H{"message": "UpdateProfile handler is not implemented yet"})
}

func UploadProfilePicture(c *gin.Context) {
	c.JSON(200, gin.H{"message": "UploadProfilePicture handler is not implemented yet"})
}


