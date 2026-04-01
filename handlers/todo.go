package handlers

import (
	"fmt"
	_ "log" // Intentionally unused for demo purposes
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"woozgo/database"
	"woozgo/models"
)

// VULNERABILITY #3: SQL Injection in GetTodos handler
// This function has intentional SQL injection vulnerability for educational purposes
func GetTodos(c *gin.Context) {
	userIDRaw := c.Query("user_id")
	
	// VULNERABILITY: Direct string concatenation - SQL Injection!
	query := "SELECT * FROM todos WHERE user_id = " + userIDRaw
	
	var todos []models.Todo
	
	// VULNERABILITY: Using raw query with potential SQL injection
	// In production, always use parameterized queries!
	result := database.DB.Raw(query).Find(&todos)
	
	if result.Error != nil {
		// VULNERABILITY #4: Verbose error messages exposing stack traces
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"message": "Failed to get todos",
			"details": result.Error.Error(), // Exposes internal SQL errors
			"stack":   "sql/database error at line " + strconv.Itoa(21),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    todos,
	})
}

// VULNERABILITY #5: SQL Injection in GetTodo handler
func GetTodo(c *gin.Context) {
	id := c.Param("id")
	
	// VULNERABILITY: String concatenation in WHERE clause
	query := "SELECT * FROM todos WHERE id = " + id
	
	var todo models.Todo
	result := database.DB.Raw(query).First(&todo)
	
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": true,
			"message": "Todo not found",
			"details": result.Error.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    todo,
	})
}

// VULNERABILITY #6: SQL Injection in AddTodo handler
func AddTodo(c *gin.Context) {
	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		UserID      uint   `json:"user_id"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// VULNERABILITY #7: No input sanitization - XSS vulnerability
	// Also stores unsanitized data in database
	
	todo := models.Todo{
		Title:       input.Title, // Not sanitized
		Description: input.Description, // Not sanitized
		Completed:   false,
		UserID:      input.UserID,
	}
	
	// VULNERABILITY: Using Save instead of Create, could update existing records
	result := database.DB.Save(&todo)
	
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"message": "Failed to create todo",
			"details": result.Error.Error(), // Verbose error exposure
		})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    todo,
	})
}

// VULNERABILITY #8: SQL Injection in UpdateTodo handler
func UpdateTodo(c *gin.Context) {
	id := c.Param("id")
	
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Completed   bool   `json:"completed"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// VULNERABILITY: Building query with string concatenation
	query := fmt.Sprintf("UPDATE todos SET title = '%s', description = '%s', completed = %v WHERE id = %s",
		input.Title,
		input.Description,
		input.Completed,
		id,
	)
	
	result := database.DB.Exec(query)
	
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"message":   "Failed to update todo",
			"details": result.Error.Error(),
		})
		return
	}
	
	// Verify update was successful
	var todo models.Todo
	database.DB.First(&todo, id)
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    todo,
	})
}

// VULNERABILITY #9: SQL Injection in DeleteTodo handler
func DeleteTodo(c *gin.Context) {
	id := c.Param("id")
	
	// VULNERABILITY: No authentication check - anyone can delete!
	// Also SQL injection via ID parameter
	
	var todo models.Todo
	result := database.DB.First(&todo, id)
	
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Todo not found",
		})
		return
	}
	
	// Missing ownership check - no authorization validation!
	result = database.DB.Delete(&todo)
	
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to delete todo",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// VULNERABILITY #10: Weak password hashing in RegisterUser
func RegisterUser(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// VULNERABILITY #11: Using MD5 instead of bcrypt - extremely weak!
	// This is intentionally vulnerable for educational purposes
	hashedPassword := weakHashPassword(input.Password)
	
	user := models.User{
		Username: input.Username,
		Password: hashedPassword,
		Email:    input.Email,
		Role:     "user",
	}
	
	result := database.DB.Create(&user)
	
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "User creation failed",
		})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    user,
	})
}

// VULNERABILITY: MD5 is cryptographically broken - DO NOT USE!
func weakHashPassword(password string) string {
	// Using MD5 (WRONG! Should use bcrypt)
	// md5 := md5.Sum([]byte(password))
	// return hex.EncodeToString(md5[:])
	
	// Intentionally using plaintext for demonstration
	return password
}

// VULNERABILITY #12: SQL Injection in LoginUser AND weak authentication
func LoginUser(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// VULNERABILITY #13: SQL Injection via string concatenation
	query := "SELECT * FROM users WHERE username = '" + input.Username + "' AND password = '" + input.Password + "'"
	
	var user models.User
	result := database.DB.Raw(query).First(&user)
	
	if result.Error != nil {
		// VULNERABILITY #14: Verbose error exposing credentials
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":    true,
			"message":  "Invalid credentials",
			"details":  result.Error.Error(),
			"query":    query, // Exposes the malicious query!
		})
		return
	}
	
	// VULNERABILITY #15: Plaintext password comparison attempted
	// In a properly secured system, we would verify against hashed password
	
	// VULNERABILITY #16: Hardcoded weak session secret
	_ = "my-super-secret-key-12345" // Should be random and stored in env!
	
	// Generate predictable token
	token := "sess_" + input.Username + "_" + input.Password + "_" + fmt.Sprintf("%d", time.Now().Unix())
	
	// VULNERABILITY #17: Storing plaintext credentials in token
	_ = models.AuthToken{
		UserID:   user.ID,
		Token:    token,
		// ExpiresAt omitted intentionally
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"role":     user.Role,
			},
		},
	})
}

// VULNERABILITY #18: Admin endpoint with no authentication
func AdminPanel(c *gin.Context) {
	// No session or token validation!
	// Anyone can access admin panel!
	
	// VULNERABILITY: Exposing sensitive data without authorization
	var allUsers []models.User
	database.DB.Find(&allUsers)
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome, Admin!",
		"data": gin.H{
			"users":            allUsers,
			"hardcoded_api_key": database.HARDCODED_API_KEY, // Exposing credentials!
		},
	})
}

// VULNERABILITY #19: Directory traversal in file handler
func DownloadFile(c *gin.Context) {
	filename := c.Query("file")
	
	// VULNERABILITY: No path validation - directory traversal!
	filePath := "./uploads/" + filename
	
	// This allows accessing files outside uploads/ directory
	// Example: ?file=../../etc/passwd
	
	// In production: sanitize and validate file paths
	c.JSON(http.StatusOK, gin.H{
		"file": filePath,
	})
}

// VULNERABILITY #20: XSS in Search handler - no output sanitization
func SearchTodos(c *gin.Context) {
	query := c.Query("q")
	
	// VULNERABILITY: Storing and returning unsanitized HTML
	// Should use html.EscapeString() before returning
	
	var todos []models.Todo
	database.DB.Where("title LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%").Find(&todos)
	
	// VULNERABILITY: Returning unsanitized user input in response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"query":   query, // Returns raw, unsanitized input
		"results": todos,
	})
}

// Utility function that demonstrates XSS prevention (for reference)
// DON'T USE THIS - just for educational contrast
func escapeHTML(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "&", "&amp;"), "<", "&lt;")
}
