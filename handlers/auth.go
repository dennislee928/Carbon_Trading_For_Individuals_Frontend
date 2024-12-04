package handlers

import (
	"context"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/mailgun/mailgun-go/v4"
)

//
type OTPRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}
//
type OTPVerification struct {
    Email string `json:"email" binding:"required,email"`
    OTP   string `json:"otp" binding:"required"`
}



// VerifyOTP handles OTP verification
func VerifyOTP(c *gin.Context) {
    var req OTPRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }

    // Validate email format
    if !isValidEmail(req.Email) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
        return
    }

    // Validate password
    if !isValidPassword(req.Password) {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Password must contain at least 8 characters up to 20 characters, " +
                "1 uppercase, 1 lowercase, 1 number and 1 special character",
        })
        return
    }

    // Generate OTP
    otp := generateOTP()

    // Store OTP in your database with email and expiration time
    // TODO: Implement database storage logic here

    // Send OTP via email
    if err := sendOTPEmail(req.Email, otp); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}
//
func isValidEmail(email string) bool {
    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    return emailRegex.MatchString(email)
}
//
func isValidPassword(password string) bool {
    if len(password) < 8 || len(password) > 20 {
        return false
    }

    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
    hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

    return hasUpper && hasLower && hasNumber && hasSpecial
}
//
func generateOTP() string {
    rand.Seed(time.Now().UnixNano())
    digits := "0123456789"
    otp := ""
    for i := 0; i < 6; i++ {
        otp += string(digits[rand.Intn(len(digits))])
    }
    return otp
}
//
func sendOTPEmail(email, otp string) error {
    // Get Mailgun credentials from environment variables
    domain := os.Getenv("MAILGUN_DOMAIN")
    apiKey := os.Getenv("MAILGUN_API_KEY")
    sender := os.Getenv("MAILGUN_SENDER")

    mg := mailgun.NewMailgun(domain, apiKey)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
    defer cancel()

    // Create email message
    subject := "Your OTP Code"
    body := "Your OTP code is: " + otp + "\nThis code will expire in 10 minutes."

    message := mg.NewMessage(sender, subject, body, email)

    // Send email
    _, _, err := mg.Send(ctx, message)
    return err
}
//
// Verify the OTP code
func VerifyOTPCode(c *gin.Context) {
    var req OTPVerification

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }

    // TODO: Implement OTP verification logic here
    // 1. Retrieve stored OTP from database
    // 2. Check if OTP is expired
    // 3. Compare submitted OTP with stored OTP
    // 4. Mark OTP as used after successful verification

    // For now, returning a success message
    c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}

// SocialLogin handles social media login
func SocialLogin(c *gin.Context) {
    provider := c.Param("provider")
    // Implement social login logic for different providers
    c.JSON(http.StatusOK, gin.H{"message": "Logged in with " + provider})
}

// Login handles user login
func Login(c *gin.Context) {
    // Authenticate the user (replace with actual authentication logic)
    userID := 1 // Example user ID after authentication

    // Generate JWT with session ID
    token, err := generateJWT(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}


// ForgotPassword handles password reset requests
func ForgotPassword(c *gin.Context) {
    // Implement password reset logic (send email with reset link or OTP)
    c.JSON(http.StatusOK, gin.H{"message": "Password reset instructions sent"})
}

// ChangePassword handles password change requests
func ChangePassword(c *gin.Context) {
    // Implement password change logic for authenticated users
    c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// generateJWT is a helper function to generate JWT tokens
func generateJWT(userID int) (string, error) {
    token := jwt.New(jwt.SigningMethodHS256)

    claims := token.Claims.(jwt.MapClaims)
    claims["user_id"] = userID
    claims["jti"] = generateSessionID() // Add session ID
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

    secretKey := os.Getenv("JWT_SECRET_KEY")
    if secretKey == "" {
        secretKey = "your_secret_key" // Fallback secret key (not recommended for production)
    }

    tokenString, err := token.SignedString([]byte(secretKey))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}
//
// handlers/middleware.go


func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
            c.Abort()
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
            c.Abort()
            return
        }

        tokenString := parts[1]

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }

            secretKey := os.Getenv("JWT_SECRET_KEY")
            if secretKey == "" {
                secretKey = "your_secret_key"
            }

            return []byte(secretKey), nil
        })

        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            c.Abort()
            return
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            c.Set("user_id", claims["user_id"])
            c.Set("session_id", claims["jti"]) // Extract session ID

            // Optional: Validate the session ID (e.g., check against a database)
            if !isValidSession(claims["jti"].(string)) {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
                c.Abort()
                return
            }

            c.Next()
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            c.Abort()
            return
        }
    }
}

// Helper function to generate a unique session identifier
func generateSessionID() string {
    return uuid.New().String() // Generates a random unique ID
}

// Example validation function for session ID
func isValidSession(_ string) bool {
    // TODO: Implement session validation logic (e.g., check against Redis or a database)
    return true
}
