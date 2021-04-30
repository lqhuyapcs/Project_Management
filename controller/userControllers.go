package controller

import (
	m "Projectmanagement_BE/models"
	u "Projectmanagement_BE/utils"
	"encoding/json"
	"net/http"
)

// RequestUserID - form to get request user id
type RequestUserID struct {
	UserID *uint `json:"user_id" sql:"-"`
}

//RegisterUser - controller
var RegisterUser = func(w http.ResponseWriter, r *http.Request) {

	user := &m.User{}
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := user.Create() //Create account
	u.Respond(w, resp)
}

//AuthenticateUser - controller
var AuthenticateUser = func(w http.ResponseWriter, r *http.Request) {

	user := &m.User{}
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := m.UserAuthenticate(*user.Mail, *user.Password)
	u.Respond(w, resp)
}

//UpdateUser - controller
var UpdateUser = func(w http.ResponseWriter, r *http.Request) {
	user := &m.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	UserID := r.Context().Value("user").(uint)

	resp := user.Update(UserID) //Update user by ID
	u.Respond(w, resp)
}

// SearchProject - controller
var SearchProject = func(w http.ResponseWriter, r *http.Request) {
	request := &m.RequestSearchUserProject{}
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

	if request.Query == nil {
		result, ok := m.GetProjectByUserID(UserID, request.Status, request.PageSize, request.PageIndex)
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
	} else if request.Query != nil {
		result, ok := m.SearchProject(UserID, *request.Query, request.Status, request.PageSize, request.PageIndex)
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
	} else {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	u.Respond(w, u.Message(true, ""))
	return
}

// SearchTask - controller
var SearchTask = func(w http.ResponseWriter, r *http.Request) {
	request := &m.RequestSearchUserTask{}
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

	if request.Query != nil {
		result, ok := m.GetTaskByUserID(UserID, request.Status, request.PageSize, request.PageIndex)
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
	} else if request.Query != nil {
		result, ok := m.SearchTask(UserID, *request.Query, request.Status, request.PageSize, request.PageIndex)
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
	} else {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	u.Respond(w, u.Message(true, ""))
	return
}

// GetUserByID - controller
var GetUserByID = func(w http.ResponseWriter, r *http.Request) {
	request := &RequestUserID{}
	err := json.NewDecoder(r.Body).Decode(request) // decode request body
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	if request.UserID == nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	user, ok := m.GetUserByID(*request.UserID)
	if !ok {
		u.Respond(w, u.Message(false, "Error when find user"))
		return
	}
	if ok {
		if user == nil {
			u.Respond(w, u.Message(false, "User not found"))
			return
		}
	}
	resp := u.Message(true, "")
	resp["user"] = user
	u.Respond(w, resp)
}

// SearchUser - controller
var SearchUser = func(w http.ResponseWriter, r *http.Request) {
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

	result, ok := m.SearchUser(UserID, request.Query, request.ProjectID, request.PageSize, request.PageIndex)
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
