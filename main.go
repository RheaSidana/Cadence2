package main

import (
    "encoding/json"
    "net/http"
    "strconv"
    "sync"
)

type Flower struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Color string `json:"color"`
    Price float64 `json:"price"`
}

var (
    flowers = []Flower{
        {ID: 1, Name: "Rose", Color: "Red", Price: 10.5},
        {ID: 2, Name: "Tulip", Color: "Yellow", Price: 8.0},
    }
    mu sync.Mutex
)

// Get all flowers
func getFlowers(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(flowers)
}

// Add a new flower
func addFlower(w http.ResponseWriter, r *http.Request) {
    var flower Flower
    if err := json.NewDecoder(r.Body).Decode(&flower); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    mu.Lock()
    flower.ID = len(flowers) + 1
    flowers = append(flowers, flower)
    mu.Unlock()
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(flower)
}

// Get a specific flower by ID
func getFlowerByID(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Path[len("/flowers/"):])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    for _, flower := range flowers {
        if flower.ID == id {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(flower)
            return
        }
    }
    http.Error(w, "Flower not found", http.StatusNotFound)
}

// Delete a flower by ID
func deleteFlower(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Path[len("/flowers/"):])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    mu.Lock()
    for i, flower := range flowers {
        if flower.ID == id {
            flowers = append(flowers[:i], flowers[i+1:]...)
            w.WriteHeader(http.StatusNoContent)
            mu.Unlock()
            return
        }
    }
    mu.Unlock()
    http.Error(w, "Flower not found", http.StatusNotFound)
}

func main() {
    http.HandleFunc("/flowers", getFlowers)
    http.HandleFunc("/flowers/", getFlowerByID)
    http.HandleFunc("/flowers/add", addFlower)
    http.HandleFunc("/flowers/delete/", deleteFlower)
    
    http.ListenAndServe(":8080", nil)
}
