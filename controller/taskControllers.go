package controller

import (
	m "Projectmanagement_BE/models"
	u "Projectmanagement_BE/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

// RequestListUserTask - struct to get form request create task
type RequestListUserTask struct {
	TaskID   *uint           `json:"task_id" sql:"-"`
	ListUserID []*uint `json:"user_ids" sql:"-"`
}

// RequestTaskID - struct to get form request set status
type RequestTaskID struct {
	TaskID *uint `json:"task_id" sql:"-"`
}

// RequestUserTask
type RequestUserTask struct {
	TaskID   *uint           `json:"task_id" sql:"-"`
	UserID *uint `json:"user_id" sql:"-"`
}

// CreateTask - controller
var CreateTask = func(w http.ResponseWriter, r *http.Request) {
	task := &m.Task{}
	err := json.NewDecoder(r.Body).Decode(task)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	UserID := r.Context().Value("user").(uint) //Grab the id of the user that send the request

	resp := task.Create(UserID) //Create task with user id
	u.Respond(w, resp)
}

// UpdateTask - controller
var UpdateTask = func(w http.ResponseWriter, r *http.Request) {

	task := &m.Task{}
	err := json.NewDecoder(r.Body).Decode(task) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	UserID := r.Context().Value("user").(uint) //Grab the id of the user that send the request

	resp := task.Update(UserID)
	u.Respond(w, resp)
}

// AssignTask - controller
var AssignTask = func(w http.ResponseWriter, r *http.Request) {

	UserID := r.Context().Value("user").(uint)

	// Assign users
	request := RequestUserTask{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	if request.TaskID == nil || request.UserID == nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := m.AddMember2Task(UserID, *request.UserID, *request.TaskID)

	u.Respond(w, resp)
}

// UnassignTask - controller
var UnassignTask = func(w http.ResponseWriter, r *http.Request) {
	UserID := r.Context().Value("user").(uint)

	// Assign users
	request := RequestUserTask{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	if request.TaskID == nil || request.UserID == nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := m.RemoveUserFromTask(UserID, *request.UserID, *request.TaskID)

	u.Respond(w, resp)
}

// SetTODOTask - controller
var SetTODOTask = func(w http.ResponseWriter, r *http.Request) {

	UserID := r.Context().Value("user").(uint)

	request := &RequestTaskID{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	if request.TaskID == nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := m.SetStatusTODO(UserID, *request.TaskID)
	u.Respond(w, resp)
}

// SetDOINGTask - controller
var SetDOINGTask = func(w http.ResponseWriter, r *http.Request) {

	UserID := r.Context().Value("user").(uint)

	request := &RequestTaskID{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	if request.TaskID == nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := m.SetStatusDOING(UserID, *request.TaskID)
	u.Respond(w, resp)
}

// SetDONETask - controller
var SetDONETask = func(w http.ResponseWriter, r *http.Request) {

	UserID := r.Context().Value("user").(uint)

	request := &RequestTaskID{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	if request.TaskID == nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := m.SetStatusDONE(UserID, *request.TaskID)
	u.Respond(w, resp)
}

// SetWAITINGTask - controller
var SetWAITINGTask = func(w http.ResponseWriter, r *http.Request) {

	UserID := r.Context().Value("user").(uint)

	request := &RequestTaskID{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	if request.TaskID == nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := m.SetStatusWAITING(UserID, *request.TaskID)
	u.Respond(w, resp)
}

// SetDELETETask - controller
var SetDELETETask = func(w http.ResponseWriter, r *http.Request) {

	UserID := r.Context().Value("user").(uint)

	request := &RequestTaskID{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	if request.TaskID == nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := m.SetStatusDELETE(UserID, *request.TaskID)
	u.Respond(w, resp)
}

// GetTaskByID - controller
var GetTaskByID = func(w http.ResponseWriter, r *http.Request) {
	request := &RequestTaskID{}
	err := json.NewDecoder(r.Body).Decode(request) // decode request body
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	UserID := r.Context().Value("user").(uint) //Grab the id of the user that send the request

	if request.TaskID == nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	task, ok := m.GetTaskByID(*request.TaskID)
	if !ok {
		u.Respond(w, u.Message(false, "Error when find task"))
		return
	}
	if ok {
		if task == nil {
			u.Respond(w, u.Message(false, "Task not found"))
			return
		}
	}

	// check relation between user and project
	if task.ProjectID == nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	userProject, ok := m.GetUserProject(UserID, *task.ProjectID)
	if !ok {
		u.Respond(w, u.Message(false, "Error when find project and user relation"))
		return
	}
	if ok {
		if userProject == nil {
			u.Respond(w, u.Message(false, "No relation between user and project"))
			return
		}
	}

	// if project found and user in project
	resp := u.Message(true, "")
	resp["task"] = task
	u.Respond(w, resp)
}

// SearchUserInTask - controller
var SearchUserInTask = func(w http.ResponseWriter, r *http.Request) {
	request := &m.RequestSearchUserInTask{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	if request.PageIndex == nil || request.PageSize == nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	UserID := r.Context().Value("user").(uint)
	if request.TaskID != nil {
		result, ok := m.SearchUserInTask(UserID, request.Query, request.TaskID, request.PageSize, request.PageIndex)
		if ok {
			if result != nil {
				resp := u.Message(true, "")
				resp["result"] = result
				u.Respond(w, resp)
				return
			}
		}
		if !ok {
			u.Respond(w, u.Message(false, "Error when connect to database"))
			return
		}
	}
	u.Respond(w, u.Message(true, ""))
	return
}

// CreateSubtask - controller
var CreateSubtask = func(w http.ResponseWriter, r *http.Request) {
	subtask := &m.SubTask{}
	err := json.NewDecoder(r.Body).Decode(subtask)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	UserID := r.Context().Value("user").(uint) //Grab the id of the user that send the request

	resp := subtask.Create(UserID, subtask.TaskID) //Create subtask with userid and taskid
	u.Respond(w, resp)
}

// UpdateSubtask - controller
var UpdateSubTask = func(w http.ResponseWriter, r *http.Request) {
	subtask := &m.SubTask{}
	err := json.NewDecoder(r.Body).Decode(subtask)

	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	UserID := r.Context().Value("user").(uint) //Grab the id of the user that send the request

	resp := subtask.Update(UserID) //Create subtask with userid and taskid
	u.Respond(w, resp)
}
