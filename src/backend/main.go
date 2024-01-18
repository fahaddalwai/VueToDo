package main

import (
    "database/sql"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/go-sql-driver/mysql"
    "log"
    "net/http"
    "time"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    "fmt"
    "strings"
)

// ToDoItem represents a single todo item.
type ToDoItem struct {
    ID    string `json:"id"`
    Label string `json:"label"`
    Done  bool   `json:"done"`
    Username string `json:"username"`
}

type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

var db *sql.DB

func init() {
    // Replace with your database connection details.
    cfg := mysql.Config{
        User:   "fahad",
        Passwd: "fahad159",
        Net:    "tcp",
        Addr:   "127.0.0.1:3306",
        DBName: "todoproject",
    }
    // Get a database handle.
    var err error
    db, err = sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        log.Fatal(err)
    }

    if err = db.Ping(); err != nil {
        log.Fatal(err)
    }
}

func main() {
    r := gin.Default()

    r.Use(cors.New(cors.Config{
        AllowOrigins: []string{"http://localhost:8081"},
        AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
    }))

    r.GET("/todo", TokenAuthMiddleware(), getToDoItems)
    r.POST("/todo", TokenAuthMiddleware(), addToDoItem)
    r.PUT("/todo/:id", TokenAuthMiddleware(), updateToDoItem)
    r.DELETE("/todo/:id", TokenAuthMiddleware(), deleteToDoItem)
    r.POST("/signup", Signup)
    r.POST("/login", Login)
    
    r.Run(":8080") // Listen and serve on 0.0.0.0:8080
}

func getToDoItems(c *gin.Context) {
    username, exists := c.Get("username")
    if !exists {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
        return
    }

    rows, err := db.Query("SELECT id, label, done FROM todo_items WHERE username = ?", username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    items := make([]ToDoItem, 0)
    for rows.Next() {
        var item ToDoItem
        if err := rows.Scan(&item.ID, &item.Label, &item.Done); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        items = append(items, item)
    }
    if err := rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, items)
}

func addToDoItem(c *gin.Context) {
    var item ToDoItem
    if err := c.BindJSON(&item); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    username, exists := c.Get("username")
    if !exists {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
        return
    }
    item.Username = username.(string)
    item.ID = uuid.NewString()

    _, err := db.Exec("INSERT INTO todo_items (id, label, done, username) VALUES (?, ?, ?, ?)", item.ID, item.Label, item.Done, item.Username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, item)
}

func TokenAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
            c.Abort()
            return
        }

        // Split the Authorization header to get the token part.
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
            c.Abort()
            return
        }

        tokenString := parts[1]

        claims := &jwt.RegisteredClaims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        c.Set("username", claims.Subject)
        c.Next()
    }
}



func updateToDoItem(c *gin.Context) {
    id := c.Param("id")
    var updatedItem ToDoItem
    if err := c.BindJSON(&updatedItem); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    _, err := db.Exec("UPDATE todo_items SET label = ?, done = ? WHERE id = ?", updatedItem.Label, updatedItem.Done, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    updatedItem.ID = id // Ensure the ID remains unchanged
    c.JSON(http.StatusOK, updatedItem)
}

func deleteToDoItem(c *gin.Context) {
    id := c.Param("id")
    _, err := db.Exec("DELETE FROM todo_items WHERE id = ?", id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "item deleted"})
}

func Signup(c *gin.Context) {
    var newUser User

    if err := c.BindJSON(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Hash password (use bcrypt or similar)
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
        return
    }
    newUser.Password = string(hashedPassword)

    err = createUser(&newUser)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

var jwtKey = []byte("your_secret_key") // Keep this key secret

func GenerateToken(username string) (string, error) {
    expirationTime := time.Now().Add(1 * time.Hour)

    claims := &jwt.RegisteredClaims{
        Subject:   username,
        ExpiresAt: &jwt.NumericDate{expirationTime},
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}


func Login(c *gin.Context) {
    var credentials User

    if err := c.BindJSON(&credentials); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := findUserByUsername(credentials.Username)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Compare hashed password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    token, err := GenerateToken(user.Username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}

func findUserByUsername(username string) (*User, error) {
    var user User

    // Assuming you have a table named 'users' with fields 'username' and 'password'
    query := "SELECT username, password FROM users WHERE username = ?"
    err := db.QueryRow(query, username).Scan(&user.Username, &user.Password)

    if err != nil {
        if err == sql.ErrNoRows {
            // User not found
            return nil, fmt.Errorf("user not found")
        }
        // Other error
        return nil, err
    }

    return &user, nil
}


func createUser(user *User) error {
    // The password should already be hashed at this point
    query := "INSERT INTO users (username, password) VALUES (?, ?)"

    _, err := db.Exec(query, user.Username, user.Password)
    if err != nil {
        return fmt.Errorf("error creating user: %v", err)
    }

    return nil
}
