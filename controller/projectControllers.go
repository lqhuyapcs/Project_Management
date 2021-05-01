package controller

import (
	m "Projectmanagement_BE/models"
	u "Projectmanagement_BE/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

// RequestProjectID - form to get request user id
type RequestProjectID struct {
	ProjectID *uint `json:"project_id" sql:"-"`
}

// RequestUserProject
type RequestUserProject struct {
	ProjectID	*uint	`json:"project_id" sql:"-"`
	UserID		*uint 	`json:"user_id" sql:"-"`
}

// CreateProject - controller
var CreateProject = func(w http.ResponseWriter, r *http.Request) {

	project := &m.Project{}
	err := json.NewDecoder(r.Body).Decode(project) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		fmt.Println(err)
		return
	}
	user := r.Context().Value("user").(uint) //Grab the id of the user that send the request

	resp := project.Create(user) //Create project with user id
	u.Respond(w, resp)
}

// UpdateProject - controller
var UpdateProject = func(w http.ResponseWriter, r *http.Request) {

	project := &m.Project{}
	err := json.NewDecoder(r.Body).Decode(project) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		fmt.Println(err)
		return
	}
	UserID := r.Context().Value("user").(uint) //Grab the id of the user that send the request

	resp := project.Update(UserID)
	u.Respond(w, resp)
}

// AddMember2Project - controller
var AddMember2Project = func(w http.ResponseWriter, r *http.Request) {

	request := RequestUserProject{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	UserID := r.Context().Value("user").(uint)

	resp := m.AddMember2Project(UserID, *request.UserID, *request.ProjectID)

	u.Respond(w, resp)

}

// GetProjectByID - controller
var GetProjectByID = func(w http.ResponseWriter, r *http.Request) {
	request := &RequestProjectID{}
	err := json.NewDecoder(r.Body).Decode(request) // decode request body
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	UserID := r.Context().Value("user").(uint) //Grab the id of the user that send the request

	if request.ProjectID == nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	project, ok := m.GetProjectByID(*request.ProjectID)
	if !ok {
		u.Respond(w, u.Message(false, "Error when find project"))
		return
	}
	if ok {
		if project == nil {
			u.Respond(w, u.Message(false, "Project not found"))
			return
		}
	}

	// check relation between user and project
	userProject, ok := m.GetUserProject(UserID, *request.ProjectID)
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

	resp := u.Message(true, "")
	resp["project"] = project
	resp["role"] = userProject.Role

	u.Respond(w, resp)

}

// SearchUserInProject - controller
var SearchUserInProject = func(w http.ResponseWriter, r *http.Request) {
	request := &m.RequestSearchUser{}
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
	if request.ProjectID != nil {
		result, ok := m.SearchUserInProject(UserID, request.Query, request.ProjectID, request.PageSize, request.PageIndex)
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

// SearchTaskInProject - controller
var SearchTaskInProject = func(w http.ResponseWriter, r *http.Request) {
	request := &m.RequestSearchProjectTask{}
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
	if request.ProjectID != nil {
		result, ok := m.SearchTaskInProject(UserID, request.Query, request.ProjectID, request.Status, request.PageSize, request.PageIndex)
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

// SearchUserTaskInProject - controller
var SearchUserTaskInProject = func(w http.ResponseWriter, r *http.Request) {
	request := &m.RequestSearchUserTaskInProject{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	if request.PageIndex == nil || request.PageSize == nil {

		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	if request.UserID == nil {

		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	UserRequestID := r.Context().Value("user").(uint)

	if request.ProjectID != nil {
		result, ok := m.SearchUserTaskInProject(UserRequestID, request.UserID, request.ProjectID, request.Query, request.Status, request.PageSize, request.PageIndex)
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

// RemoveUserFromProject - controller
var RemoveUserFromProject = func(w http.ResponseWriter, r *http.Request) {
	request := &RequestUserProject{}
	err := json.NewDecoder(r.Body).Decode(request) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		fmt.Println(err)
		return
	}
	UserID := r.Context().Value("user").(uint)

	if request.ProjectID == nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := m.DeleteUserProject(UserID, *request.UserID, *request.ProjectID)

	u.Respond(w, resp)
}
