package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
    "github.com/gorilla/mux"
    "github.com/google/uuid"
)

type Idea struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
}

var ideas []Idea

func main() {
    router := mux.NewRouter()

    // Middleware for JSON content type
    router.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "application/json")
            next.ServeHTTP(w, r)
        })
    })

    // Routes
    router.HandleFunc("/api/ideas", getIdeas).Methods("GET")
    router.HandleFunc("/api/ideas", createIdea).Methods("POST")
    router.HandleFunc("/api/ideas/{id}", getIdea).Methods("GET")
    router.HandleFunc("/api/ideas/{id}", updateIdea).Methods("PUT")
    router.HandleFunc("/api/ideas/{id}", deleteIdea).Methods("DELETE")

    log.Fatal(http.ListenAndServe(":8080", router))
}

func getIdeas(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(ideas)
}

func createIdea(w http.ResponseWriter, r *http.Request) {
    var newIdea Idea
    if err := json.NewDecoder(r.Body).Decode(&newIdea); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    newIdea.ID = uuid.New().String()
    newIdea.CreatedAt = time.Now()
    ideas = append([]Idea{newIdea}, ideas...) // Add to beginning of slice

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(newIdea)
}

func getIdea(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for _, item := range ideas {
        if item.ID == params["id"] {
            json.NewEncoder(w).Encode(item)
            return
        }
    }
    http.Error(w, "Idea not found", http.StatusNotFound)
}

func updateIdea(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var updatedIdea Idea
    if err := json.NewDecoder(r.Body).Decode(&updatedIdea); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    for i, item := range ideas {
        if item.ID == params["id"] {
            updatedIdea.ID = item.ID
            updatedIdea.CreatedAt = item.CreatedAt
            ideas[i] = updatedIdea
            json.NewEncoder(w).Encode(updatedIdea)
            return
        }
    }
    http.Error(w, "Idea not found", http.StatusNotFound)
}

func deleteIdea(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for i, item := range ideas {
        if item.ID == params["id"] {
            ideas = append(ideas[:i], ideas[i+1:]...)
            w.WriteHeader(http.StatusNoContent)
            return
        }
    }
    http.Error(w, "Idea not found", http.StatusNotFound)
}