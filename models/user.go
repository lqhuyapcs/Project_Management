package models

import (
	u "Projectmanagement_BE/utils"
	"os"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Token struct
type Token struct {
	UserID uint
	jwt.StandardClaims
}

// UserID struct
type UserID struct {
	TaskID uint `json:"task_id" gorm:"-"`
	UserID uint `json:"user_id" gorm:"-"`
}

// User struct
type User struct {
	gorm.Model
	FullName *string	`json:"fullname"`
	Mail     *string 	`json:"mail"`
	Password *string   	`json:"password"`
	AvatarUrl 	 string	`json:"avatar"`
	Projects []Project  `gorm:"many2many:user_projects" json:"projects"`
	Tasks    []Task     `gorm:"many2many:user_tasks" json:"tasks"`
	Token    string     `json:"token" gorm:"-"`
}

// Create - user model
func (user *User) Create() map[string]interface{} {

	if user.Mail == nil || user.Password == nil ||  user.FullName == nil {
		return u.Message(false, "Invalid request")
	}


	if msg, status := user.Validate(); !status {
		return msg
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	*user.Password = string(hashedPassword)
	user.AvatarUrl = u.RandomAvatarUrl()

	GetDB().Create(user)

	if user.ID <= 0 {
		return u.Message(false, "Failed to create user, connection error.")
	}
	user.Password = nil //delete password

	response := u.Message(true, "User has been created")
	response["user"] = user
	return response
}

// Update - user model
func (user *User) Update(UserID uint) map[string]interface{} {

	// Get user by UserID
	updatedUser, ok := GetUserByID(UserID)
	if ok {
		if updatedUser == nil {
			return u.Message(false, "User not found")
		}
	}
	if !ok {
		return u.Message(false, "Error when query user")
	}

	// To update record, need to check all the valid request
	// mail
	if user.Mail != nil {
		temp := &User{}
		//check for errors and duplicate user name
		err := GetDB().Table("users").Where("mail = ?", user.Mail).First(temp).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return u.Message(false, "Connection error. Please retry.")
		}
		if temp.Mail != nil {
			return u.Message(false, "Mail exists.")
		}
		updatedUser.Mail = user.Mail
	}

	// password
	if user.Password != nil {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
		*updatedUser.Password = string(hashedPassword)
		updatedUser.Password = nil
	}
	GetDB().Save(updatedUser)

	// Respond
	response := u.Message(true, "")
	response["user"] = updatedUser
	return response
}

// UserAuthenticate - user model
func UserAuthenticate(mail string, password string) map[string]interface{} {

	user := &User{}
	user, status := GetUserByMail(mail)
	if status {
		if user == nil {
			return u.Message(false, "User not found")
		}
	} else {
		return u.Message(false, "Connection error")
	}

	err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Wrong password.")
	}

	//Worked! Logged In
	user.Password = nil
	projects, _ := GetListProjectByUserID(user.ID)
	user.Projects = *projects

	
	//Create JWT token
	tk := &Token{UserID: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS512"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString //Store the token in the response

	resp := u.Message(true, "Logged In")
	resp["user"] = user
	return resp
}

// GetUserByID - user model
func GetUserByID(id uint) (*User, bool) {
	user := &User{}
	err := GetDB().Table("users").Where("id = ?", id).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}

	user.Password = nil
	return user, true
}

// GetUserByMail - user model
func GetUserByMail(mail string) (*User, bool) {
	user := &User{}
	err := GetDB().Table("users").Where("mail = ?", mail).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}
	return user, true
}

// GetListUserByTaskID - user model
func GetListUserByTaskID(TaskID uint) (*[]User, bool) {
	listUser := &[]User{}
	err := GetDB().Table("users").Select("id", "mail", "full_name", "avatar_url").Joins("join user_tasks on users.id = user_tasks.user_id").
		Where("user_tasks.task_id = ?", TaskID).Find(listUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}
	return listUser, true
}

// GetListUserByProjectID - user model
func GetListUserByProjectID(ProjectID uint) (*[]User, bool) {
	listUser := &[]User{}
	err := GetDB().Table("users").Select("id", "mail", "full_name", "avatar_url").Joins("join user_projects on users.id = user_projects.user_id").
		Where("user_projects.project_id = ?", ProjectID).Find(listUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}

	return listUser, true
}

// Validate - user model
func (user *User) Validate() (map[string]interface{}, bool) {

	temp := &User{}

	//check for errors and duplicate user name
	err := GetDB().Table("users").Where("mail = ?", user.Mail).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry."), false
	}
	if temp.Mail != nil {
		return u.Message(false, "Mail already in use by another user."), false
	}

	return u.Message(false, "Requirement passed."), true
}

