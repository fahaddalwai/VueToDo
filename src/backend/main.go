package main

import (
    "database/sql"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/go-sql-driver/mysql"
    "log"
    "net/http"
)

// ToDoItem represents a single todo item.
type ToDoItem struct {
    ID    string `json:"id"`
    Label string `json:"label"`
    Done  bool   `json:"done"`
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

    r.GET("/todo", getToDoItems)
    r.POST("/todo", addToDoItem)
    r.PUT("/todo/:id", updateToDoItem)
    r.DELETE("/todo/:id", deleteToDoItem)

    r.Run(":8080") // Listen and serve on 0.0.0.0:8080
}

func getToDoItems(c *gin.Context) {
    rows, err := db.Query("SELECT id, label, done FROM todo_items")
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
    item.ID = uuid.NewString()

    _, err := db.Exec("INSERT INTO todo_items (id, label, done) VALUES (?, ?, ?)", item.ID, item.Label, item.Done)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, item)
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
