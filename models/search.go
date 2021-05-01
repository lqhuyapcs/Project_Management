package models

import (
	"gorm.io/gorm"
)

// RequestSearchUserProject - model
type RequestSearchUserProject struct {
	Query     *string `json:"query" sql:"-"`
	Status    *uint   `json:"status" sql:"-"`
	PageSize  *uint   `json:"page_size" sql:"-"`
	PageIndex *uint   `json:"page_index" sql:"-"`
}

// RequestSearchUserTask - model
type RequestSearchUserTask struct {
	Query     *string `json:"query" sql:"-"`
	Status    *uint   `json:"status" sql:"-"`
	PageSize  *uint   `json:"page_size" sql:"-"`
	PageIndex *uint   `json:"page_index" sql:"-"`
}

// RequestSearchProjectTask - model
type RequestSearchProjectTask struct {
	ProjectID *uint   `json:"project_id" sql:"-"`
	Query     *string `json:"query" sql:"-"`
	Status    *uint   `json:"status" sql:"-"`
	PageSize  *uint   `json:"page_size" sql:"-"`
	PageIndex *uint   `json:"page_index" sql:"-"`
}

// RequestSearchUser - model
type RequestSearchUser struct {
	ProjectID *uint   `json:"project_id" sql:"-"`
	Query     *string `json:"query" sql:"-"`
	PageSize  *uint   `json:"page_size" sql:"-"`
	PageIndex *uint   `json:"page_index" sql:"-"`
}

// RequestSearchUserInTask - model
type RequestSearchUserInTask struct {
	TaskID    *uint   `json:"task_id" sql:"-"`
	Query     *string `json:"query" sql:"-"`
	PageSize  *uint   `json:"page_size" sql:"-"`
	PageIndex *uint   `json:"page_index" sql:"-"`
}

// RequestSearchUserTaskInProject - model
type RequestSearchUserTaskInProject struct {
	UserID    *uint   `json:"user_id" sql:"-"`
	ProjectID *uint   `json:"project_id" sql:"-"`
	Query     *string `json:"query" sql:"-"`
	Status    *uint   `json:"status_id" sql:"-"`
	PageSize  *uint   `json:"page_size" sql:"-"`
	PageIndex *uint   `json:"page_index" sql:"-"`
}

// SearchProject - model
func SearchProject(UserID uint, Query string, Status *uint, PageSize *uint, PageIndex *uint) (*[]Project, bool) {
	project := &[]Project{}

	if Status == nil {
		if PageSize != nil && PageIndex != nil {
			pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
			err := GetDB().Table("projects").Joins("join user_projects on projects.id = user_projects.project_id").
				Where("user_projects.user_id = ? AND to_tsvector('english', projects.name) @@ plainto_tsquery('english', ?) and user_projects.deleted_at IS NULL", UserID, Query).
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
		pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
		err := GetDB().Table("projects").Joins("join user_projects on projects.id = user_projects.project_id").
			Where("user_projects.user_id = ? AND projects.status = ? AND to_tsvector('english', projects.name) @@ plainto_tsquery('english', ?) and user_projects.deleted_at IS NULL", UserID, Status, Query).
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

// SearchTask - model
func SearchTask(UserID uint, Query string, Status *uint, PageSize *uint, PageIndex *uint) (*[]Task, bool) {
	task := &[]Task{}

	if Status == nil {
		if PageSize != nil && PageIndex != nil {
			pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
			err := GetDB().Table("tasks").Joins("join user_tasks on tasks.id = user_tasks.task_id").
				Where("user_tasks.user_id = ? AND to_tsvector('english', tasks.name) @@ plainto_tsquery('english', ?) and user_tasks.deleted_at IS NULL", UserID, Query).
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
		err := GetDB().Table("tasks.").Joins("join user_tasks. on tasks..id = user_tasks.task_id").
			Where("user_tasks.user_id = ? AND tasks.status = ? AND to_tsvector('english', tasks.name) @@ plainto_tsquery('english', ?) and user_tasks.deleted_at IS NULL", UserID, Status, Query).
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

// SearchUser - model
func SearchUser(UserID uint, Query *string, ProjectID *uint, PageSize *uint, PageIndex *uint) (*[]User, bool) {
	user := &[]User{}
	if ProjectID != nil {
		// check role of user request and project
		roleUserReq, ok := GetRoleByUserProjectID(UserID, *ProjectID)
		if ok {
			if roleUserReq == "" {
				return nil, false
			}
		} else {
			return nil, false
		}
	}

	if Query != nil {
		if ProjectID != nil {
			if PageSize != nil && PageIndex != nil {
				pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
				err := GetDB().Table("users").Distinct("id").
					Where("users.id not in (select user_id from user_projects where user_projects.project_id = ? and user_projects.deleted_at IS NULL) AND users.mail LIKE ?", *ProjectID, "%"+*Query+"%").
					Offset(offset).Limit(pageSize).Preload("Employee").Find(user).Error
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						return nil, true
					}
					return nil, false
				}
			}
		} else {
			if PageSize != nil && PageIndex != nil {
				pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
				err := GetDB().Table("users").
					Where("users.mail LIKE ? AND users.id <> ?", "%"+*Query+"%", UserID).
					Offset(offset).Limit(pageSize).Preload("Employee").Find(user).Error
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						return nil, true
					}
					return nil, false
				}
			}
		}
	} else {
		if ProjectID != nil {
			if PageSize != nil && PageIndex != nil {
				pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
				err := GetDB().Table("users").Distinct("id").
					Where("users.id not in (SELECT user_id from user_projects where user_projects.project_id = ? and user_projects.deleted_at IS NULL) AND users.id <> ?", *ProjectID, UserID).
					Offset(offset).Limit(pageSize).Preload("Employee").Find(user).Error
				if err != nil {

					if err == gorm.ErrRecordNotFound {
						return nil, true
					}
					return nil, false
				}
			}
		} else {
			if PageSize != nil && PageIndex != nil {
				pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
				err := GetDB().Table("users").Where("id <> ?", UserID).Offset(offset).Limit(pageSize).Preload("Employee").Find(user).Error
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						return nil, true
					}
					return nil, false
				}
			}
		}
	}
	return user, true
}

// SearchUserInProject - model
func SearchUserInProject(UserID uint, Query *string, ProjectID *uint, PageSize *uint, PageIndex *uint) (*[]User, bool) {
	user := &[]User{}
	if ProjectID != nil {
		// check role of user request and project
		roleUserReq, ok := GetRoleByUserProjectID(UserID, *ProjectID)
		if ok {
			if roleUserReq == "" {
				return nil, false
			}
		} else {
			return nil, false
		}
	}

	if Query != nil {
		if PageSize != nil && PageIndex != nil {
			pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
			err := GetDB().Table("users").Distinct("id").
				Where("users.id in (select user_id from user_projects where user_projects.project_id = ? and user_projects.deleted_at IS NULL ) AND users.mail LIKE ? AND users.id <> ?", *ProjectID, "%"+*Query+"%", UserID).
				Offset(offset).Limit(pageSize).Preload("Employee").Find(user).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return nil, true
				}
				return nil, false
			}
		}
	} else {
		if PageSize != nil && PageIndex != nil {
			pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
			err := GetDB().Table("users").Distinct("id").
				Where("users.id in (select user_id from user_projects where user_projects.project_id = ? and user_projects.deleted_at IS NULL ) AND users.id <> ?", *ProjectID, UserID).
				Offset(offset).Limit(pageSize).Preload("Employee").Find(user).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return nil, true
				}
				return nil, false
			}
		}
	}
	return user, true
}

// SearchTaskInProject - model
func SearchTaskInProject(UserID uint, Query *string, ProjectID *uint, Status *uint, PageSize *uint, PageIndex *uint) (*[]Task, bool) {
	task := &[]Task{}
	if ProjectID != nil {
		// check role of user request and project
		roleUserReq, ok := GetRoleByUserProjectID(UserID, *ProjectID)
		if ok {
			if roleUserReq == "" {
				return nil, false
			}
		} else {
			return nil, false
		}
	}

	if Query != nil {
		if Status != nil {
			if PageSize != nil && PageIndex != nil {
				pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
				err := GetDB().Table("tasks").
					Where("tasks.project_id = ? AND tasks.status = ? AND to_tsvector('english', tasks.name) @@ plainto_tsquery('english', ?)", *ProjectID, *Query).
					Offset(offset).Limit(pageSize).Preload("Subtasks").Find(task).Error
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						return nil, true
					}
					return nil, false
				}
			}
		} else {
			if PageSize != nil && PageIndex != nil {
				pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
				err := GetDB().Table("tasks").
					Where("tasks.project_id = ? AND to_tsvector('english', tasks.name) @@ plainto_tsquery('english', ?)", *ProjectID, *Query).
					Offset(offset).Limit(pageSize).Preload("Subtasks").Find(task).Error
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						return nil, true
					}
					return nil, false
				}
			}
		}
	} else {
		if Status != nil {
			if PageSize != nil && PageIndex != nil {
				pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
				err := GetDB().Table("tasks").
					Where("tasks.project_id = ? AND tasks.status = ?", *ProjectID, *Status).
					Offset(offset).Limit(pageSize).Preload("Subtasks").Find(task).Error
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						return nil, true
					}
					return nil, false
				}
			}
		} else {
			if PageSize != nil && PageIndex != nil {
				pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
				err := GetDB().Table("tasks").
					Where("tasks.project_id = ?", *ProjectID).
					Offset(offset).Limit(pageSize).Preload("Subtasks").Find(task).Error
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						return nil, true
					}
					return nil, false
				}
			}
		}
	}
	return task, true
}

// CalculatePaginate - model
func CalculatePaginate(PageSize uint, PageIndex uint) (int, int) {
	if PageIndex == 0 {
		PageIndex = 1
	}

	switch {
	case PageSize > 100:
		PageSize = 100
	case PageSize <= 0:
		PageSize = 10
	}
	pageIndex := int(PageIndex)
	pageSize := int(PageSize)
	offset := int((pageIndex - 1) * pageSize)
	return pageSize, offset
}

// SearchUserTaskInProject - model
func SearchUserTaskInProject(RequestUserID uint, UserID *uint, ProjectID *uint, Query *string, Status *uint, PageSize *uint, PageIndex *uint) (*[]Task, bool) {
	task := &[]Task{}

	// check role of user request and project
	roleUserReq, ok := GetRoleByUserProjectID(*UserID, *ProjectID)
	if ok {
		if roleUserReq == "" {
			return nil, false
		}
	} else {
		return nil, false
	}

	// if user request is in project, search task of user
	// check whether it is in this project or not
	if Query != nil {
		if Status == nil {
			if PageSize != nil && PageIndex != nil {
				pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
				err := GetDB().Table("tasks").Joins("join user_tasks on tasks.id = user_tasks.task_id").
					Where("user_tasks.user_id = ? AND tasks.project_id = ? AND to_tsvector('english', tasks.name) @@ plainto_tsquery('english', ?) and user_tasks.deleted_at IS NULL", *UserID, *ProjectID, *Query).
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
			err := GetDB().Table("tasks.").Joins("join user_tasks on tasks.id = user_tasks.task_id").
				Where("user_tasks.user_id = ? AND tasks.status = ? AND asks.project_id = ? AND to_tsvector('english', tasks.name) @@ plainto_tsquery('english', ?) and user_tasks.deleted_at IS NULL", *UserID, *Status, *ProjectID, *Query).
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
	if Status == nil {
		if PageSize != nil && PageIndex != nil {
			pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
			err := GetDB().Table("tasks").Joins("join user_tasks on tasks.id = user_tasks.task_id").
				Where("user_tasks.user_id = ? AND tasks.project_id = ? and user_tasks.deleted_at IS NULL", *UserID, *ProjectID).
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
		err := GetDB().Table("tasks.").Joins("join user_tasks. on tasks.id = user_tasks.task_id").
			Where("user_tasks.user_id = ? AND tasks.status = ? AND tasks.project_id = ? and user_tasks.deleted_at IS NULL", *UserID, *Status, *ProjectID).
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

// SearchUserInTask - model
func SearchUserInTask(UserID uint, Query *string, TaskID *uint, PageSize *uint, PageIndex *uint) (*[]User, bool) {
	user := &[]User{}
	// get project
	project, ok := GetProjectByTaskID(*TaskID)
	if ok {
		if project == nil {
			return nil, false
		}
	} else {
		return nil, false
	}
	// check role of user request and project
	roleUserReq, ok := GetRoleByUserProjectID(UserID, project.ID)
	if ok {
		if roleUserReq == "" {
			return nil, false
		}
	} else {
		return nil, false
	}

	if Query != nil {
		pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
		err := GetDB().Table("users").Distinct("id").
			Where("users.id in (select user_id from user_tasks where user_tasks.task_id = ? and user_tasks.deleted_at IS NULL ) AND users.mail LIKE ? AND users.id <> ?", *TaskID, "%"+*Query+"%", UserID).
			Offset(offset).Limit(pageSize).Preload("Employee").Find(user).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, true
			}
			return nil, false
		}
	} else {
		pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
		err := GetDB().Table("users").
			Where("users.id in (select user_id from user_tasks where user_tasks.task_id = ? and user_tasks.deleted_at IS NULL) AND users.id <> ?", *TaskID, UserID).
			Offset(offset).Limit(pageSize).Preload("Employee").Find(user).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, true
			}
			return nil, false
		}

	}
	return user, true
}
