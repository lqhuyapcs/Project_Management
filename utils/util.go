package utils

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Message - util
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// Respond - util
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// MultipleRespond - util
func MultipleRespond(w http.ResponseWriter, data []map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Returns an int >= min, < max
func randomInt(min, max int) int {
    return min + rand.Intn(max-min)
}

// Random avatar
func RandomAvatarUrl() string {
	rand.Seed(time.Now().UnixNano())
	number := randomInt(1, 7)
	url := os.Getenv("api_icon") + "icon_" + strconv.Itoa(number) + ".png" 
	return url
}