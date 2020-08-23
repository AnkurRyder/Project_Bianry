package network

import (
	"Project_binary/db"
	"Project_binary/types"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// DeleteAuth function deletes metadata from db
func DeleteAuth(givenUUID string, db *gorm.DB) (int64, error) {
	var tokenTemp types.TokenMeta
	db.Where("acces_id = ?", givenUUID).First(&tokenTemp)
	if tokenTemp.AccesID != givenUUID {
		return 0, errors.New("Already Logged out")
	}
	db.Delete(&tokenTemp)
	return 1, nil
}

// CheckAuth for checking auth
func CheckAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenAuth, err := ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		_, err = FetchAuth(tokenAuth, db)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

// CreateToken function generating token for the user
func CreateToken(userID uint64) (*types.TokenDetails, error) {
	var err error
	td := &types.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute).Unix()
	td.AccessUuid = uuid.New().String()

	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userID
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(db.GoDotEnvVariable("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

// CreateAuth function for creating auth to db
func CreateAuth(userid uint64, td *types.TokenDetails, db *gorm.DB) {
	var tk types.TokenMeta
	tk.AccesID = td.AccessUuid
	tk.UserID = userid
	tk.ExpTime = td.AtExpires
	db.Create(&tk)
}

func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(db.GoDotEnvVariable("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// ExtractTokenMetadata return metadata encrypted in the token
func ExtractTokenMetadata(r *http.Request) (*types.AccessDetails, error) {
	token, err := verifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &types.AccessDetails{
			AccessUuid: accessUUID,
			UserId:     userID,
		}, nil
	}
	return nil, err
}

// FetchAuth returns the data stored in db corresponding to the token
func FetchAuth(authD *types.AccessDetails, db *gorm.DB) (uint64, error) {
	var tokenMeta types.TokenMeta
	db.Where("acces_id = ?", authD.AccessUuid).First(&tokenMeta)
	if tokenMeta.AccesID != authD.AccessUuid {
		return 0, errors.New("Token Expired")
	}
	return tokenMeta.UserID, nil
}
