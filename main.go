package main

import (
	app "Projectmanagement_BE/app"
	controller "Projectmanagement_BE/controller"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func handleRequests() {
	router := mux.NewRouter()

	// default api
	router.HandleFunc("/", nil)

	// users api
	router.Handle("/user", notImplement)
	router.HandleFunc("/api/user/login", controller.AuthenticateUser).Methods("POST") 
	router.HandleFunc("/api/user/register", controller.RegisterUser).Methods("POST")  
	router.HandleFunc("/api/user/update-info", controller.UpdateUser).Methods("POST") 
	router.HandleFunc("/api/user/get-by-id", controller.GetUserByID).Methods("POST")
	router.HandleFunc("/api/user/search-project", controller.SearchProject).Methods("POST")
	router.HandleFunc("/api/user/search-task", controller.SearchTask).Methods("POST")
	router.HandleFunc("/api/user/search-user", controller.SearchUser).Methods("POST")

	// projects api
	router.HandleFunc("/api/project/create", controller.CreateProject).Methods("POST")                
	router.HandleFunc("/api/project/add-user", controller.AddMember2Project).Methods("POST") 
	router.HandleFunc("/api/project/search-user", controller.SearchUserInProject).Methods("POST")
	router.HandleFunc("/api/project/search-task", controller.SearchTaskInProject).Methods("POST")
	router.HandleFunc("/api/project/update-info", controller.UpdateProject).Methods("POST")         
	router.HandleFunc("/api/project/remove-user", controller.RemoveUserFromProject).Methods("POST") 
	router.HandleFunc("/api/project/get-by-id", controller.GetProjectByID).Methods("POST")
	router.HandleFunc("/api/project/search-user-task", controller.SearchUserTaskInProject).Methods("POST")

	// tasks api
	router.HandleFunc("/api/task/create", controller.CreateTask).Methods("POST")         
	router.HandleFunc("/api/task/assign", controller.AssignTask).Methods("POST")          
	router.HandleFunc("/api/task/set-todo", controller.SetTODOTask).Methods("POST")       
	router.HandleFunc("/api/task/set-doing", controller.SetDOINGTask).Methods("POST")    
	router.HandleFunc("/api/task/set-done", controller.SetDONETask).Methods("POST")       
	router.HandleFunc("/api/task/set-waiting", controller.SetWAITINGTask).Methods("POST") 
	router.HandleFunc("/api/task/set-delete", controller.SetDELETETask).Methods("POST")   
	router.HandleFunc("/api/task/update-info", controller.UpdateTask).Methods("POST")     
	router.HandleFunc("/api/task/unassign-user", controller.UnassignTask).Methods("POST") 
	router.HandleFunc("/api/task/get-by-id", controller.GetTaskByID).Methods("POST")
	router.HandleFunc("/api/task/create-subtask", controller.CreateSubtask).Methods("POST") 
	router.HandleFunc("/api/task/update-subtask", controller.UpdateSubTask).Methods("POST")
	router.HandleFunc("/api/task/search-user", controller.SearchUserInTask).Methods("POST")



	router.Use(app.JwtAuthentication)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	err := http.ListenAndServe(":"+port, router)
	fmt.Println(err)
	if err == nil {
	}
}

func main() {
	handleRequests()
}

// in case api is not implemented yet
var notImplement = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not implemented"))
})
