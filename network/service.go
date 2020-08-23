package network

import (
	"Project_binary/log"
	"Project_binary/types"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	guuid "github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// GetDataHelper function to help service getData
func GetDataHelper(c *gin.Context, db *gorm.DB) (interface{}, types.ErrorRes) {
	userData := types.Data{}
	id := c.Param("id")
	err := validation.Validate(&id, validation.Required, is.UUID)
	if err != nil {
		log.Slogger.Infof("Validation failed Error: %s", err.Error())
		return nil, types.ErrorRes{HTTPCode: http.StatusBadRequest, Err: err}
	}
	ok := db.Where("Id = ?", id).First(&userData).RecordNotFound()
	if ok {
		log.Slogger.Infof("for Id: %s record not found in db", id)
		return nil, types.ErrorRes{HTTPCode: http.StatusNotFound, Err: errors.New("Record Not Found")}
	}
	return userData, types.ErrorRes{}
}

// WriteDataHelper function to help service writeData
func WriteDataHelper(c *gin.Context, db *gorm.DB) (interface{}, types.ErrorRes) {
	user := types.Data{ID: guuid.New(), Value: false, Key: ""}
	c.BindJSON(&user)
	err := validation.ValidateStruct(&user,
		validation.Field(&user.ID, validation.Empty),
		validation.Field(&user.Value, validation.Required),
		validation.Field(&user.Key, validation.Required),
	)
	if err != nil {
		log.Slogger.Infof("Validation failed Error: %s", err.Error())
		return nil, types.ErrorRes{HTTPCode: http.StatusBadRequest, Err: err}
	}
	db.Create(&user)
	return user, types.ErrorRes{}
}

// ModifyDataHelper function to help service ModifyData
func ModifyDataHelper(c *gin.Context, db *gorm.DB) (interface{}, types.ErrorRes) {
	var userData types.Data
	var user types.Data
	id := c.Param("id")
	err := validation.Validate(&id, validation.Required, is.UUID)
	if err != nil {
		log.Slogger.Infof("Validation failed Error: %s", err.Error())
		return nil, types.ErrorRes{HTTPCode: http.StatusBadRequest, Err: err}
	}
	c.BindJSON(&user)
	err = validation.ValidateStruct(&user,
		validation.Field(&user.ID, validation.Empty),
		validation.Field(&user.Value, validation.Required),
		validation.Field(&user.Key, validation.Required),
	)
	if err != nil {
		log.Slogger.Infof("Validation failed Error: %s", err.Error())
		return nil, types.ErrorRes{HTTPCode: http.StatusBadRequest, Err: err}
	}
	user.ID, err = guuid.Parse(id)
	if err != nil {
		log.Slogger.Infof("Error: %s", err.Error())
		return nil, types.ErrorRes{HTTPCode: http.StatusBadRequest, Err: err}
	}
	n := db.Model(&userData).Where("Id = ?", id).Update(map[string]interface{}{"Value": user.Value, "Key": user.Key}).RowsAffected
	if n == 0 {
		log.Slogger.Infof("for Id: %s, could not update the record in db", id)
		return nil, types.ErrorRes{HTTPCode: http.StatusNotModified, Err: errors.New("Check Id")}
	}
	return user, types.ErrorRes{}
}

// DeleteDataHelper function to help service DeleteData
func DeleteDataHelper(c *gin.Context, db *gorm.DB) (interface{}, types.ErrorRes) {
	var user types.Data
	id := c.Param("id")
	idUUID, err := guuid.Parse(id)
	if err != nil {
		return nil, types.ErrorRes{HTTPCode: http.StatusBadRequest, Err: errors.New("Wrong ID")}
	}
	db.Where("Id = ?", idUUID).First(&user)
	if user.ID != idUUID {
		return nil, types.ErrorRes{HTTPCode: http.StatusBadRequest, Err: errors.New("Wrong ID or no data exists")}
	}
	db.Delete(&user)
	return "No Content", types.ErrorRes{}
}

// LoginHelper function to help service Login
func LoginHelper(c *gin.Context, db *gorm.DB) (interface{}, types.ErrorRes) {
	var u types.User
	var user types.User
	err := c.ShouldBindJSON(&u)
	if err != nil {
		return nil, types.ErrorRes{HTTPCode: http.StatusBadRequest, Err: errors.New("Invalid json provided")}
	}
	err = validation.ValidateStruct(&u,
		validation.Field(&u.ID, validation.Empty),
		validation.Field(&u.Username, validation.Required),
		validation.Field(&u.Password, validation.Required),
	)
	if err != nil {
		log.Slogger.Infof("Validation failed Error: %s", err.Error())
		return nil, types.ErrorRes{HTTPCode: http.StatusBadRequest, Err: err}
	}
	ok := db.Where("username = ?", u.Username).First(&user).RecordNotFound()
	if ok {
		return nil, types.ErrorRes{HTTPCode: http.StatusBadRequest, Err: errors.New("UserName Incorrct")}
	}
	if user.Username != u.Username || user.Password != u.Password {
		return nil, types.ErrorRes{HTTPCode: http.StatusBadRequest, Err: errors.New("Invalid login details")}
	}
	token, err := CreateToken(user.ID)
	if err != nil {
		return nil, types.ErrorRes{HTTPCode: http.StatusBadRequest, Err: err}
	}
	CreateAuth(user.ID, token, db)
	return token.AccessToken, types.ErrorRes{}
}

//SignUpHelper function to help service SignUp
func SignUpHelper(c *gin.Context, db *gorm.DB) (interface{}, types.ErrorRes) {
	var userData types.User
	var userDbCheck types.User
	err := c.ShouldBindJSON(&userData)
	err = validation.ValidateStruct(&userData,
		validation.Field(&userData.ID, validation.Empty),
		validation.Field(&userData.Username, validation.Required),
		validation.Field(&userData.Password, validation.Required),
	)
	if err != nil {
		log.Slogger.Infof("Validation failed Error: %s", err.Error())
		return nil, types.ErrorRes{HTTPCode: http.StatusBadRequest, Err: err}
	}
	if err != nil {
		return nil, types.ErrorRes{HTTPCode: http.StatusBadRequest, Err: errors.New("Invalid json provided")}
	}
	db.Where("username = ?", userData.Username).First(&userDbCheck)
	if userDbCheck.Username == userData.Username {
		return nil, types.ErrorRes{HTTPCode: http.StatusBadRequest, Err: errors.New("Username already exists")}
	}
	db.Create(&userData)
	return "Account Created", types.ErrorRes{}
}

//LogoutHelepr function to help service Logout
func LogoutHelepr(c *gin.Context, db *gorm.DB) (interface{}, types.ErrorRes) {
	au, err := ExtractTokenMetadata(c.Request)
	if err != nil {
		return nil, types.ErrorRes{HTTPCode: http.StatusUnauthorized, Err: errors.New("unauthorized")}
	}
	deleted, err := DeleteAuth(au.AccessUuid, db)
	if deleted == 0 || err != nil {
		return nil, types.ErrorRes{HTTPCode: http.StatusUnauthorized, Err: errors.New("unauthorized")}
	}
	return "Successfully logged out", types.ErrorRes{}
}
