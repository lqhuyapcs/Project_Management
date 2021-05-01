package models

import (
	u "Projectmanagement_BE/utils"
	"time"

	"gorm.io/gorm"
)

// Task struct
type Task struct {
	gorm.Model
	Name        *string   `json:"name"`
	CreatorID   uint      `json:"creator_id"`
	ProjectID   *uint     `json:"project_id"`
	Deadline    *string   `json:"date_deadline"`
	Description *string   `json:"description"`
	Status      *uint     `json:"status"`
	Subtasks    []SubTask `json:"subtasks"`
	Users		[]User    `gorm:"many2many:user_projects" json:"users"`
}

// UserTask struct
type UserTask struct {
	UserID    uint           `json:"user_id" gorm:"primaryKey"`
	TaskID    uint           `json:"task_id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"date_created"`
	UpdateAt  time.Time      `json:"date_updated"`
	DeletedAt gorm.DeletedAt `json:"date_deleted"`
}

// SubTask struct
type SubTask struct {
	gorm.Model
	TaskID      uint    `json:"task_id"`
	Description *string `json:"description"`
	IsDone      *bool   `json:"is_done"`
}

// Create Task - model
func (task *Task) Create(UserID uint) map[string]interface{} {

	if task.ProjectID == nil || task.Description == nil || task.Name == nil {
		return u.Message(false, "Invalid request")
	}
	// check role of user request and project
	roleUserReq, ok := GetRoleByUserProjectID(UserID, *task.ProjectID)
	if ok {
		if roleUserReq == "" {
			return u.Message(false, "No role between user request and project")
		}
	} else {
		return u.Message(false, "Connection error when query role between user request and project")
	}
	// only admin can create task
	if roleUserReq != "admin" {
		return u.Message(false, "Only admin can create task")
	}

	// check valid task status
	statusType := [5]uint{1, 2, 3, 4, 5}
	errStatus := false
	for i := range statusType {
		if *task.Status == statusType[i] {
			errStatus = true
			break
		}
	}
	if !errStatus {
		return u.Message(false, "Invalid task status")
	}

	task.CreatorID = UserID
	GetDB().Create(task)

	if task.ID <= 0 {
		return u.Message(false, "Failed to create task, connection error??")
	}

	resp := u.Message(true, "Task has been created")
	resp["task"] = task
	return resp
}

// Update - model
func (task *Task) Update(UserID uint) map[string]interface{} {

	if task.Name == nil && task.Description == nil && task.Deadline == nil {
		return u.Message(false, "Invalid request")
	}

	// get project
	project, ok := GetProjectByTaskID(task.ID)
	if ok {
		if project == nil {
			return u.Message(false, "Task is not in this project")
		}
	} else {
		return u.Message(false, "Connection error when query project")
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
	// only admin can create task
	if roleUserReq != "admin" {
		return u.Message(false, "Only admin can create task")
	}

	updatedTask, ok := GetTaskByID(task.ID)
	if ok {
		if updatedTask == nil {
			return u.Message(false, "Task not found")
		}
	}
	if !ok {
		return u.Message(false, "Error when query task")
	}

	if task.Name != nil {
		updatedTask.Name = task.Name
	}
	if task.Description != nil {
		updatedTask.Description = task.Description
	}
	if task.Deadline != nil {
		updatedTask.Deadline = task.Deadline
	}
	GetDB().Save(updatedTask)

	response := u.Message(true, "")
	response["task"] = updatedTask
	return response
}

// AddMember2Task - model
func AddMember2Task(UserRequestID uint, UserID uint, TaskID uint) map[string]interface{} {

	// get project
	project, ok := GetProjectByTaskID(TaskID)
	if ok {
		if project == nil {
			return u.Message(false, "Task is not in this project")
		}
	} else {
		return u.Message(false, "Connection error when query project")
	}

	// check role of user request and project
	roleUserReq, ok := GetRoleByUserProjectID(UserRequestID, project.ID)
	if ok {
		if roleUserReq == "" {
			return u.Message(false, "No role between user request and project")
		}
	} else {
		return u.Message(false, "Connection error when query role between user request and project")
	}
	// only admin can add member
	if roleUserReq != "admin" {
		return u.Message(false, "Only admin can add member")
	}

	//--------------- User request passed ----------------

	// check relation of user added and project
	userProjectAdded, ok := GetUserProject(UserID, project.ID)
	if userProjectAdded == nil {
		if !ok {
			return u.Message(false, "Connection error when query relation between user and project")
		}
		return u.Message(false, "This user is not in project")
	}

	// check relation of user added and task
	userTaskAdded, ok := GetUserTask(UserID, TaskID)
	if userTaskAdded != nil {
		return u.Message(false, "User already in task")
	}
	if userTaskAdded == nil {
		if !ok {
			return u.Message(false, "Connection error when query relation between user and task")
		}
	}

	userTask := &UserTask{
		UserID: UserID,
		TaskID: TaskID,
	}

	if GetDB().Create(userTask).Error != nil {
		return u.Message(false, "Error when add user to task")
	}
	resp := u.Message(true, "User has been added to task")
	return resp
}

// RemoveUserFromTask - model
func RemoveUserFromTask(UserRequestID uint, UserID uint, TaskID uint) map[string]interface{} {

	// check relation of user removed and task
	userTaskRemoved, ok := GetUserTask(UserID, TaskID)
	if userTaskRemoved == nil {
		if !ok {
			return u.Message(false, "Connection error when query relation between user request and task")
		}
		if ok {
			return u.Message(false, "User not in task")
		}
	}

	// If user request is user removed, unassign from task
	if UserRequestID == UserID {
		Err := GetDB().Table("user_tasks").Where("user_id = ? AND task_id = ?", UserID, TaskID).Delete(userTaskRemoved).Error
		if Err != nil {
			return u.Message(false, "Error when unassign from task")
		}
		return u.Message(true, "")
	}

	// get project of task
	project, ok := GetProjectByTaskID(TaskID)
	if ok {
		if project == nil {
			return u.Message(false, "Task is not in this project")
		}
	} else {
		return u.Message(false, "Connection error when find project of this task")
	}

	// check role of user request and project
	roleUserReq, ok := GetRoleByUserProjectID(UserRequestID, project.ID)
	if ok {
		if roleUserReq == "" {
			return u.Message(false, "No role between user request and project")
		}
	} else {
		return u.Message(false, "Connection error when query role between user request and project")
	}
	// only admin can remove user from task
	if roleUserReq != "admin" {
		return u.Message(false, "Only admin can remove user from task")
	}

	/*--------Remove user from task----------*/
	Err := GetDB().Table("user_tasks").Where("user_id = ? AND task_id = ?", UserID, TaskID).Delete(userTaskRemoved).Error
	if Err != nil {
		return u.Message(false, "Error when unassign from task")
	}

	return u.Message(true, "")

}

// Create SubTask - model
func (subTask *SubTask) Create(UserID uint, TaskID uint) map[string]interface{} {
	// get project
	project, ok := GetProjectByTaskID(TaskID)
	if ok {
		if project == nil {
			return u.Message(false, "Task is not available")
		}
	} else {
		return u.Message(false, "Connection error when query project")
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
	// only admin can add subtask
	if roleUserReq != "admin" {
		return u.Message(false, "Only admin can add subtask")
	}

	GetDB().Create(subTask)

	if subTask.ID <= 0 {
		return u.Message(false, "Failed to create subtask, connection error.")
	}

	resp := u.Message(true, "Sub task has been created")
	resp["subtask"] = subTask
	return resp
}

// Update SubTask - model
func (subtask *SubTask) Update(UserID uint) map[string]interface{} {

	updatedSubTask, ok := GetSubtaskByID(subtask.ID)
	if ok {
		if updatedSubTask == nil {
			return u.Message(false, "Subtask not found")
		}
	}
	if !ok {
		return u.Message(false, "Error when query subtask")
	}
	TaskID := updatedSubTask.TaskID
	// get project
	project, ok := GetProjectByTaskID(TaskID)
	if ok {
		if project == nil {
			return u.Message(false, "Task is not available")
		}
	} else {
		return u.Message(false, "Connection error when query project")
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
	// only admin can update subtask
	if roleUserReq != "admin" {
		return u.Message(false, "Only admin can update subtask")
	}

	if subtask.Description != nil {
		updatedSubTask.Description = subtask.Description
	}
	if subtask.IsDone != nil {
		updatedSubTask.IsDone = subtask.IsDone
	}
	GetDB().Save(updatedSubTask)

	response := u.Message(true, "")
	response["subtask"] = updatedSubTask
	return response
}

// GetUserTask - get relation of user_id and task_id
func GetUserTask(UserID uint, TaskID uint) (*UserTask, bool) {
	userTask := &UserTask{}
	err := GetDB().Table("user_tasks").Where("user_id = ? AND task_id = ?", UserID, TaskID).First(userTask).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}
	return userTask, true
}

// GetSubTaskByTaskID - get sub tasks of task by task_id
func GetSubTaskByTaskID(TaskID uint) (*[]SubTask, bool) {
	subtask := &[]SubTask{}
	err := GetDB().Table("sub_tasks").Where("task_id = ?", TaskID).Find(subtask).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}
	return subtask, true

}

// GetNotDoneSubtask - get sub tasks which are not done
func GetNotDoneSubtask(TaskID uint) (*SubTask, bool) {
	subTask := &SubTask{}
	err := GetDB().Table("sub_tasks").Where("task_id = ? and is_done = ?", TaskID, false).First(subTask).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}
	return subTask, true
}

// SetStatusTODO - set status to to do
func SetStatusTODO(UserID uint, TaskID uint) map[string]interface{} {

	// check relation of user request and task
	// userTaskAdded, ok := GetUserTask(UserID, TaskID)
	// if userTaskAdded == nil {
	// 	if !ok {
	// 		return u.Message(false, "Connection error when query relation between user and task")
	// 	}
	// 	return u.Message(false, "User is not in this task")
	// }
	// get project
	project, ok := GetProjectByTaskID(TaskID)
	if ok {
		if project == nil {
			return u.Message(false, "Task not found")
		}
	} else {
		return u.Message(false, "Connection error when query project")
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

	// update task by id
	err := GetDB().Table("tasks").Where("ID = ?", TaskID).Update("status", 1).Error
	if err != nil {
		return u.Message(false, "Error when set status")
	}
	return u.Message(true, "Set status success")
}

// SetStatusDOING - set status to doing
func SetStatusDOING(UserID uint, TaskID uint) map[string]interface{} {

	// check relation of user request and task
	// userTaskAdded, ok := GetUserTask(UserID, TaskID)
	// if userTaskAdded == nil {
	// 	if !ok {
	// 		return u.Message(false, "Connection error when query relation between user and task")
	// 	}
	// 	return u.Message(false, "User is not in this task")
	// }
	project, ok := GetProjectByTaskID(TaskID)
	if ok {
		if project == nil {
			return u.Message(false, "Task not found")
		}
	} else {
		return u.Message(false, "Connection error when query project")
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
	// update task by id
	err := GetDB().Table("tasks").Where("ID = ?", TaskID).Update("status", 2).Error
	if err != nil {
		return u.Message(false, "Error when set status")
	}

	return u.Message(true, "")
}

// SetStatusWAITING - set status to waiting
func SetStatusWAITING(UserID uint, TaskID uint) map[string]interface{} {

	project, ok := GetProjectByTaskID(TaskID)
	if ok {
		if project == nil {
			return u.Message(false, "Task not found")
		}
	} else {
		return u.Message(false, "Connection error when query project")
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

	// update task by id
	err := GetDB().Table("tasks").Where("ID = ?", TaskID).Update("status", 4).Error
	if err != nil {
		return u.Message(false, "Error when set status")
	}
	return u.Message(true, "")
}

// SetStatusDELETE - set status to delete
func SetStatusDELETE(UserID uint, TaskID uint) map[string]interface{} {

	project, ok := GetProjectByTaskID(TaskID)
	if ok {
		if project == nil {
			return u.Message(false, "Task not found")
		}
	} else {
		return u.Message(false, "Connection error when query project")
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
	// update task by id
	err := GetDB().Table("tasks").Where("ID = ?", TaskID).Update("status", 5).Error
	if err != nil {
		return u.Message(false, "Error when set status")
	}
	return u.Message(true, "Set status success")
}

// SetStatusDONE - set status to done
func SetStatusDONE(UserID uint, TaskID uint) map[string]interface{} {

	project, ok := GetProjectByTaskID(TaskID)
	if ok {
		if project == nil {
			return u.Message(false, "Task not found")
		}
	} else {
		return u.Message(false, "Connection error when query project")
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

	// check subtasks done
	subTask, ok := GetNotDoneSubtask(TaskID)
	if subTask == nil {
		if !ok {
			return u.Message(false, "Connection error when query relation between subtask")
		}
	}
	if subTask != nil { //  not done subtask exists
		return u.Message(false, "Still remains subtask")
	}

	// update task by id
	err := GetDB().Table("tasks").Where("ID = ?", TaskID).Update("status", 3).Error
	if err != nil {
		return u.Message(false, "Error when set status")
	}

	return u.Message(true, "Set status success")
}

// GetTaskByID - task model
func GetTaskByID(id uint) (*Task, bool) {
	task := &Task{}
	err := GetDB().Table("tasks").Where("id = ?", id).Preload("Subtasks").First(task).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}
	users, _ := GetListUserByTaskID(id)

	task.Users = *users
	return task, true
}

// GetSubtaskByID - task model
func GetSubtaskByID(id uint) (*SubTask, bool) {
	subtask := &SubTask{}
	err := GetDB().Table("sub_tasks").Where("id = ?", id).First(subtask).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}

	return subtask, true
}

// GetTaskByUserID - task model
func GetTaskByUserID(UserID uint, Status *uint, PageSize *uint, PageIndex *uint) (*[]Task, bool) {
	task := &[]Task{}

	if Status == nil {
		if PageSize != nil && PageIndex != nil {
			pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
			err := GetDB().Table("tasks").Joins("join user_tasks on tasks.id = user_tasks.task_id").
				Where("user_tasks.user_id = ?", UserID).
				Offset(offset).Limit(pageSize).Preload("Subtasks").Find(task).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return nil, true
				}
				return nil, false
			}
		}
		return task, true
	}
	if PageSize != nil && PageIndex != nil {
		pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
		err := GetDB().Table("tasks").Joins("join user_tasks on tasks.id = user_tasks.task_id").
			Where("user_tasks.user_id = ? AND tasks.status = ?", UserID, Status).
			Offset(offset).Limit(pageSize).Preload("Subtasks").Find(task).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, true
			}
			return nil, false
		}
	}
	return task, true
}


// GetListTaskByProjectID - task model
func GetListTaskByProjectID(ProjectID uint) (*[]Task, bool) {
	listTask := &[]Task{}
	err := GetDB().Table("tasks").
	Where("tasks.project_id = ?", ProjectID).Preload("Subtasks").Find(listTask).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}
	return listTask, true
}