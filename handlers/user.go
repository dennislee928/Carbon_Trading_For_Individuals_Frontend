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

    // Insert user into the "users" table
    data, err := client.From("users").Insert(map[string]interface{}{
        "username": user.Username,
        "email":    user.Email,
        "password": user.Password, // Note: In production, always hash the password!
    }).Exec()

    // Handle errors during insertion
    if err != nil {
        log.Println("Error inserting user:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting user"})
        return
    }

    // Return success response
    c.JSON(http.StatusOK, gin.H{
        "message": "User registered successfully",
        "data":    data,
    })
}


// NewSupabaseClient creates and returns a new Supabase client
func NewSupabaseClient() (*supabase.Client, error) {
	// Supabase project URL and API Key
	supabaseUrl := "https://omlzzhqhuhcpypohelbq.supabase.co" // Replace with your actual Supabase URL
	supabaseKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Im9tbHp6aHFodWhjcHlwb2hlbGJxIiwicm9sZSI6ImFub24iLCJpYXQiOjE3MzEwMzAzMzUsImV4cCI6MjA0NjYwNjMzNX0.fUKcu_mfWqESfZlAHQWjFPsvA0tHGq6t_fnTtPQhL-Q"            // Replace with your actual Supabase API Key

	client, err := supabase.NewClient(supabaseUrl, supabaseKey)
	if err != nil {
		return nil, err
	}
	return client, nil
}
