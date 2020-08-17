package network

import (
	"Project_binary/types"
	"net/http"
	"os"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Login function for helping user to login
func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var u types.User
		var user types.User
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
			return
		}
		db.Where("Username = ?", u.Username).First(&user)
		//compare the user from the request, with the one we defined:
		if user.Username != u.Username || user.Password != u.Password {
			c.JSON(http.StatusUnauthorized, "Please provide valid login details")
			return
		}
		token, err := CreateToken(user.ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, err.Error())
			return
		}
		c.JSON(http.StatusOK, token)
	}
}

// CreateToken function generating token for the user
func CreateToken(userID uint64) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userID
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
