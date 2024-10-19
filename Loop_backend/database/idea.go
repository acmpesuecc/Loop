// schema.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS ideas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

// models/idea.go
package models

import (
    "time"
    "github.com/google/uuid"
)

type Idea struct {
    ID          uuid.UUID `json:"id" db:"id"`
    Title       string    `json:"title" db:"title"`
    Description string    `json:"description" db:"description"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type CreateIdeaRequest struct {
    Title       string `json:"title" binding:"required"`
    Description string `json:"description" binding:"required"`
}

// handlers/idea.go
package handlers

import (
    "database/sql"
    "net/http"
    "your-project/models"
    "github.com/gin-gonic/gin"
    "github.com/lib/pq"
)

type IdeaHandler struct {
    db *sql.DB
}

func NewIdeaHandler(db *sql.DB) *IdeaHandler {
    return &IdeaHandler{db: db}
}

func (h *IdeaHandler) CreateIdea(c *gin.Context) {
    var req models.CreateIdeaRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    var idea models.Idea
    err := h.db.QueryRow(`
        INSERT INTO ideas (title, description)
        VALUES ($1, $2)
        RETURNING id, title, description, created_at
    `, req.Title, req.Description).Scan(
        &idea.ID,
        &idea.Title,
        &idea.Description,
        &idea.CreatedAt,
    )

    if err != nil {
        if pqErr, ok := err.(*pq.Error); ok {
            // Handle specific PostgreSQL errors
            switch pqErr.Code {
            case "23505": // unique_violation
                c.JSON(http.StatusConflict, gin.H{"error": "Idea already exists"})
                return
            default:
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create idea"})
                return
            }
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create idea"})
        return
    }

    c.JSON(http.StatusCreated, idea)
}

func (h *IdeaHandler) GetIdeas(c *gin.Context) {
    rows, err := h.db.Query(`
        SELECT id, title, description, created_at
        FROM ideas
        ORDER BY created_at DESC
    `)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ideas"})
        return
    }
    defer rows.Close()

    var ideas []models.Idea
    for rows.Next() {
        var idea models.Idea
        if err := rows.Scan(
            &idea.ID,
            &idea.Title,
            &idea.Description,
            &idea.CreatedAt,
        ); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process ideas"})
            return
        }
        ideas = append(ideas, idea)
    }

    if err = rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while iterating ideas"})
        return
    }

    c.JSON(http.StatusOK, ideas)
}

// main.go
package main

import (
    "database/sql"
    "log"
    "your-project/handlers"
    "github.com/gin-gonic/gin"
    "github.com/lib/pq"
)

func main() {
    // Database connection
    db, err := sql.Open("postgres", "postgres://username:password@localhost:5432/dbname?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    if err = db.Ping(); err != nil {
        log.Fatal(err)
    }

    r := gin.Default()

    ideaHandler := handlers.NewIdeaHandler(db)

    r.POST("/api/ideas", ideaHandler.CreateIdea)
    r.GET("/api/ideas", ideaHandler.GetIdeas)

    if err := r.Run(":8080"); err != nil {
        log.Fatal(err)
    }
}