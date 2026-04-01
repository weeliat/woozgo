package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"woozgo/database"
	"woozgo/handlers"
)

func main() {
	// VULNERABILITY #1: Starting Gin in Debug mode (production risk)
	gin.SetMode(gin.DebugMode) // Should be gin.ReleaseMode in production!
	
	// Initialize database (with hardcoded credentials)
	database.InitDB()
	
	r := gin.Default()
	
	// VULNERABILITY #2: Missing security middleware
	// In production, should use:
	// - r.Use(cors.Default())
	// - r.Use(rate.Limiter())
	// - r.Use(secure.New(secure.Config{...}))
	// - r.Use(hsts())
	
	// VULNERABILITY #3: Logging full request details
	r.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.RequestURI
		c.Next()
		log.Printf("Method: %s, Path: %s, Time: %v", c.Request.Method, path, time.Since(start))
		// VULNERABILITY: Logging sensitive headers and body
		// log.Printf("Headers: %v", c.Request.Header)
	})
	
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"version": "1.0.0-vulnerable", // Exposing version!
			"debug":   true,                // Exposing debug mode!
		})
	})
	
	// TODO endpoints (all have SQL injection vulnerabilities)
	todos := r.Group("/api/todos")
	{
		todos.GET("", handlers.GetTodos)           // SQL Injection
		todos.GET("/:id", handlers.GetTodo)        // SQL Injection
		todos.POST("", handlers.AddTodo)           // SQL Injection + XSS
		todos.PUT("/:id", handlers.UpdateTodo)     // SQL Injection
		todos.DELETE("/:id", handlers.DeleteTodo)  // SQL Injection + No Auth
	}
	
	// Search endpoint (XSS vulnerability)
	r.GET("/api/search", handlers.SearchTodos)
	
	// Authentication endpoints (SQL Injection + weak auth)
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", handlers.RegisterUser) // Weak password hashing
		auth.POST("/login", handlers.LoginUser)       // SQL Injection
	}
	
	// VULNERABILITY #4: Admin panel with no authentication!
	admin := r.Group("/api/admin")
	{
		// No authentication or authorization check!
		admin.GET("/panel", handlers.AdminPanel) // Anyone can access!
	}
	
	// VULNERABILITY #5: File download with directory traversal
	r.GET("/api/files", handlers.DownloadFile)
	
	// VULNERABILITY #6: Exposing environment info
	r.GET("/api/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"database_host": "localhost",  // Exposed!
			"database_name": "todos",      // Exposed!
			"environment":   "production", // Exposed!
			"api_key":       database.HARDCODED_API_KEY, // Exposed!
			"secret_key":    database.HARDCODED_SECRET_KEY, // Exposed!
		})
	})
	
	// VULNERABILITY #7: No rate limiting, no CORS configuration
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "Not found"})
	})
	
	log.Println("Server starting on :8080")
	log.Println("⚠️  VULNERABLE TODO APP - FOR EDUCATIONAL PURPOSES ONLY!")
	log.Println("⚠️  DO NOT USE IN PRODUCTION!")
	
	// Start server
	r.Run(":8080")
}
