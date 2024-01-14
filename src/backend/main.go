package main

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "net/http"
    "sync"
)

// ToDoItem represents a single todo item.
type ToDoItem struct {
    ID    string `json:"id"`
    Label string `json:"label"`
    Done  bool   `json:"done"`
}

var (
    // Mutex for safe access to the ToDoItems slice.
    mutex sync.Mutex

    // ToDoItems holds the list of todo items.
    ToDoItems = []ToDoItem{}
)

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
    mutex.Lock()
    defer mutex.Unlock()

    c.JSON(http.StatusOK, ToDoItems)
}

func addToDoItem(c *gin.Context) {
    mutex.Lock()
    defer mutex.Unlock()

    var item ToDoItem
    if err := c.BindJSON(&item); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    item.ID = uuid.NewString()

    ToDoItems = append(ToDoItems, item)
    c.JSON(http.StatusCreated, item)
}

func updateToDoItem(c *gin.Context) {
    mutex.Lock()
    defer mutex.Unlock()

    id := c.Param("id")
    var updatedItem ToDoItem
    if err := c.BindJSON(&updatedItem); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    for i, item := range ToDoItems {
        if item.ID == id {
            ToDoItems[i] = updatedItem
            updatedItem.ID = id // Ensure the ID remains unchanged
            c.JSON(http.StatusOK, updatedItem)
            return
        }
    }

    c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
}

func deleteToDoItem(c *gin.Context) {
    mutex.Lock()
    defer mutex.Unlock()

    id := c.Param("id")
    for i, item := range ToDoItems {
        if item.ID == id {
            ToDoItems = append(ToDoItems[:i], ToDoItems[i+1:]...)
            c.JSON(http.StatusOK, gin.H{"message": "item deleted"})
            return
        }
    }

    c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
}
