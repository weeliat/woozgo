# 🚨 Vulnerable Todo Application - EDUCATIONAL PURPOSES ONLY 🚨

**⚠️ WARNING: This application contains INTENTIONAL SECURITY VULNERABILITIES! ⚠️**

**DO NOT USE IN PRODUCTION! This is strictly for educational purposes.**

---

## Overview

This is a **deliberately insecure** Todo application built with Go and the Gin web framework. It demonstrates common web application vulnerabilities that developers must protect against in real-world applications.

**Framework**: Gin (popular Go web framework)  
**Database**: SQLite with GORM ORM  
**Purpose**: Security education and vulnerability research

---

## 📁 Project Structure

```
woozgo/
├── main.go                        # Main entry with vulnerable routes
├── models/
│   └── todos.go                   # Data models (no validation)
├── handlers/
│   └── todo.go                    # HTTP handlers with vulnerabilities
├── database/
│   └── init.go                    # DB init with hardcoded credentials
├── go.mod                         # Go module definition
└── go.sum                         # Dependency checksums
```

---

## 🚀 Running the Application

```bash
cd woozgo
go run .
```

The server will start on `http://localhost:8080`

### Key Endpoints (All Vulnerable)

| Endpoint | Vulnerability |
|----------|--------------|
| `GET /api/todos` | SQL Injection (user_id param) |
| `GET /api/todos/:id` | SQL Injection |
| `POST /api/todos` | SQL Injection + XSS |
| `PUT /api/todos/:id` | SQL Injection |
| `DELETE /api/todos/:id` | SQL Injection + No Auth |
| `GET /api/search?q=...` | XSS |
| `POST /api/auth/register` | Weak password hashing |
| `POST /api/auth/login` | SQL Injection |
| `GET /api/admin/panel` | No Authentication |
| `GET /api/files?file=...` | Directory Traversal |
| `GET /api/info` | Credential Exposure |

---

## 🛡️ How to Fix These Vulnerabilities

### SQL Injection
**Wrong (Vulnerable)**:
```go
query := "SELECT * FROM todos WHERE user_id = " + userIDRaw
```

**Correct (Parameterized)**:
```go
var todos []models.Todo
DB.Where("user_id = ?", userIDRaw).Find(&todos)
```

### Hardcoded Credentials
**Wrong**:
```go
const AdminPassword = "SuperSecretPassword123!"
```

**Correct**:
```go
adminPass := os.Getenv("ADMIN_PASSWORD")
// Use secrets management: AWS Secrets Manager, HashiCorp Vault
```

### XSS
**Wrong**:
```go
c.JSON(http.StatusOK, gin.H{"result": userInput})
```

**Correct**:
```go
c.JSON(http.StatusOK, gin.H{"result": html.EscapeString(userInput)})
```

### Weak Passwords
**Wrong**:
```go
password := md5.Sum([]byte(input.Password))
```

**Correct**:
```go
hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
```

---

## 📚 Learning Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Gin Security Best Practices](https://github.com/gin-gonic/gin#security)
- [Go Security Best Practices](https://github.com/GoKillers/Go-Interview-Questions)
- [SQL Injection Prevention](https://cheatsheetseries.owasp.org/cheatsheets/SQL_Injection_Prevention_Cheat_Sheet.html)

---

## ⚠️ Disclaimer

This application is **strictly for educational purposes**. It demonstrates intentional security vulnerabilities. Never deploy code containing these vulnerabilities in any production environment.

**Author**: Educational security demonstration  
**License**: Educational use only  
**Version**: 1.0.0-vulnerable

---

*Built with Go 1.22, Gin v1.12, GORM v1.31*
