package models

import (
	u "Projectmanagement_BE/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Project struct
type Project struct {
	gorm.Model
	Name        *string      `json:"name"`
	CreatorID   uint         `json:"creator_id"`
	Status      *uint        `json:"status"`
	Description *string      `json:"description"`
	Users		[]User		 `gorm:"many2many:user_projects" json:"users"`
	Tasks       []Task       `json:"tasks"`
}

// UserProject struct - project user relation
type UserProject struct {
	UserID        uint           `json:"user_id" gorm:"primaryKey"`
	ProjectID     uint           `json:"project_id" gorm:"primaryKey"`
	Role        string           `json:"role"`
	AddedByUserID uint           `json:"added_by_user_id"`
	CreatedAt     time.Time      `json:"date_created"`
	UpdateAt      time.Time      `json:"date_updated"`
	DeletedAt     gorm.DeletedAt `json:"date_deleted"`
}

// Create project
func (project *Project) Create(UserID uint) map[string]interface{} {

	if project.Name == nil {
		return u.Message(false, "Invalid request")
	}
	project.CreatorID = UserID
	GetDB().Create(project)

	if project.ID <= 0 {
		return u.Message(false, "Failed to create project, connection error.")
	}

	// create relation for creator to project
	userproject := &UserProject{
		ProjectID:     project.ID,
		UserID:        UserID,
		Role:          "admin",
		AddedByUserID: UserID,
	}

	GetDB().Create(userproject)


	resp := u.Message(true, "")
	resp["project"] = project
	resp["user_project"] = userproject
	return resp
}

// Update - Project model
func (project *Project) Update(UserID uint) map[string]interface{} {

	if project.Name == nil && project.Description == nil {
		return u.Message(false, "Invalid request")
	}
	// check role of user request and project
	roleUserReq, ok := GetRoleByUserProjectID(UserID, project.ID)
	if ok {
		if roleUserReq == "" {
			return u.Message(false, "No role between user request and project")
		}
	} else {
		return u.Message(false, "Connection error when query role between user request and project")
	}
	// check permissions
	if roleUserReq != "admin" {
		return u.Message(false, "Only admin can update project")
	}	

	updatedProject, ok := GetProjectByID(project.ID)
	if ok {
		if updatedProject == nil {
			return u.Message(false, "Project not found")
		}
	}
	if !ok {
		return u.Message(false, "Error when query project")
	}

	if project.Name != nil {
		updatedProject.Name = project.Name
	}
	if project.Description != nil {
		updatedProject.Description = project.Description
	}
	GetDB().Save(updatedProject)

	response := u.Message(true, "")
	response["project"] = updatedProject
	return response
}

// AddMember2Project - Add member to project
func AddMember2Project(UserRequestID uint, UserID uint, ProjectID uint) map[string]interface{} {

	// check role of user request and project
	roleUserReq, ok := GetRoleByUserProjectID(UserRequestID, ProjectID)
	if ok {
		if roleUserReq == "" {
			return u.Message(false, "No role between user request and project")
		}
	} else {
		return u.Message(false, "Connection error when query role between user request and project")
	}
	// check permissions
	if roleUserReq != "admin" {
		return u.Message(false, "Request user doesnt have permissions to add member")
	}

	// if user add is user request 
	if UserRequestID == UserID {
		// check relation between user removed and project
		userProjectRemoved, ok := GetUserProject(UserID, ProjectID)

		if userProjectRemoved == nil {
			if !ok {
				return u.Message(false, "Connection error when query relation between user and project")
			}
		}
		return u.Message(true, "User is already in project")
	}

	/*--------------- User request passed ----------------*/

	// check relation of user added and project
	userProjectAdded, ok := GetUserProject(UserID, ProjectID)
	if userProjectAdded != nil {
		return u.Message(false, "User already in project")
	}
	if userProjectAdded == nil {
		if !ok {
			return u.Message(false, "Connection error when query relation between user and project")
		}
	}

	// Add user to project
	userProject := &UserProject{
		UserID:        UserID,
		ProjectID:     ProjectID,
		Role:          "member",
		AddedByUserID: UserRequestID,
	}

	if GetDB().Create(userProject).Error != nil {
		return u.Message(false, "Error when add user to project")
	}

	resp := u.Message(true, "User has been added")
	return resp
}

// GetUserProject - get relation of user_id and project_id
func GetUserProject(UserID uint, ProjectID uint) (*UserProject, bool) {
	userProject := &UserProject{}
	err := GetDB().Table("user_projects").Where("user_id = ? AND project_id = ?", UserID, ProjectID).First(userProject).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}
	return userProject, true
}

// GetProjectByTaskID - get project_id by task_id
func GetProjectByTaskID(TaskID uint) (*Project, bool) {
	project := &Project{}
	err := GetDB().Table("projects").Joins("join tasks on tasks.project_id = projects.id").Where("tasks.id = ? ", TaskID).First(project).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}
	return project, true
}

// GetProjectByID - project model
func GetProjectByID(id uint) (*Project, bool) {
	project := &Project{}
	err := GetDB().Table("projects").Where("id = ?", id).First(project).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}
	users, _ := GetListUserByProjectID(id)
	project.Users = *users

	tasks, _ := GetListTaskByProjectID(id)
	project.Tasks= *tasks
	return project, true
}

// DeleteUserProject - UserProject model
func DeleteUserProject(UserRequestID uint, UserID uint, ProjectID uint) map[string]interface{} {

	// check role of user request and project
	roleUserReq, ok := GetRoleByUserProjectID(UserRequestID, ProjectID)
	if ok {
		if roleUserReq == "" {
			return u.Message(false, "user is not in this project")
		}
	} else {
		return u.Message(false, "Connection error when query role between user request and project")
	}
	userProjectRemoved := &UserProject{}
	
	// if user removed is user request => leave project
	if UserRequestID == UserID {
		// check relation between user removed and project
		userProjectRemoved, ok := GetUserProject(UserID, ProjectID)

		if userProjectRemoved == nil {
			if !ok {
				return u.Message(false, "Connection error when query relation between user and project")
			}
			if ok {
				return u.Message(false, "User is not in project")
			}
		}
		// delete user - project relation
		Err := GetDB().Where("user_id = ? AND project_id = ?", UserID, ProjectID).Delete(userProjectRemoved).Error
		if Err != nil {
			return u.Message(false, "Error when delete user from project")
		}
		
		return u.Message(true, "")
	}


	// if user removed is not user request, check permissions of user request
	if UserRequestID != UserID {
		// check permissions
		if roleUserReq != "admin" {
			return u.Message(false, "User request doesnt have permissions to delete member")
		}

		// check relation between user removed and project
		userProjectRemoved, ok := GetUserProject(UserID, ProjectID)

		if userProjectRemoved == nil {
			if !ok {
				return u.Message(false, "Connection error when query relation between user and project")
			}
			if ok {
				return u.Message(false, "User request is not in project")
			}
		}

		// check role of user removed
		roleUserRemoved, ok := GetRoleByUserProjectID(UserID, ProjectID)
		if ok {
			if roleUserRemoved == "" {
				return u.Message(false, "User remove is not in project")
			}
		} else {
			return u.Message(false, "Connection error when query role between user removed and project")
		}
		if roleUserRemoved == "admin" {
			return u.Message(false, "Can not remove admin from project")
		}

	}

	/*--------------- User request passed ----------------*/
	// delete user - project relation
	Err := GetDB().Where("user_id = ? AND project_id = ?", UserID, ProjectID).Delete(userProjectRemoved).Error
	if Err != nil {
		return u.Message(false, "Error when delete user from project")
	}

	return u.Message(true, "")
}

// GetProjectByUserID - project model
func GetProjectByUserID(UserID uint, Status *uint, PageSize *uint, PageIndex *uint) (*[]Project, bool) {
	project := &[]Project{}
	if Status == nil {
		if PageSize != nil && PageIndex != nil {
			pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
			err := GetDB().Table("projects").Joins("join user_projects on projects.id = user_projects.project_id").
				Where("user_projects.user_id = ?", UserID).
				Offset(offset).Limit(pageSize).Find(project).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return nil, true
				}
				return nil, false
			}
		}
		return project, true
	}
	if PageSize != nil && PageIndex != nil {
		fmt.Println("aa")
		pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
		err := GetDB().Table("projects").Joins("join user_projects on projects.id = user_projects.project_id").
			Where("user_projects.user_id = ? AND projects.status = ?", UserID, Status).
			Offset(offset).Limit(pageSize).Find(project).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, true
			}
			return nil, false
		}
	}
	fmt.Println("??")
	return project, true
}

// GetRoleByUserProjectID - query role by user_id and project_id
func GetRoleByUserProjectID(UserID uint, ProjectID uint) (string, bool) {
	userProject := &UserProject{}
	err := GetDB().Table("user_projects").Where("user_id = ? AND project_id = ?", UserID, ProjectID).First(userProject).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", true
		}
		return "", false
	}
	return userProject.Role, true
}

// GetListProjectByUserID - project model
func GetListProjectByUserID(UserID uint) (*[]Project, bool) {
	project := &[]Project{}

	err := GetDB().Table("projects").Joins("join user_projects on projects.id = user_projects.project_id").
		Where("user_projects.user_id = ?", UserID).Find(project).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}
	return project, true
}