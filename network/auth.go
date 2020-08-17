package network

import (
	"Project_binary/db"
	"Project_binary/types"
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

// Login function for helping user to login
func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var u types.User
		var user types.User
		err := c.ShouldBindJSON(&u)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
			return
		}
		db.Where("username = ?", u.Username).First(&user)
		//compare the user from the request, with the one we defined:
		if user.Username != u.Username || user.Password != u.Password {
			c.JSON(http.StatusUnauthorized, "Please provide valid login details")
			return
		}
		token, err := createToken(user.ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, err.Error())
			return
		}
		createAuth(user.ID, token, db)
		c.JSON(http.StatusOK, token.AccessToken)
	}
}

// CheckAuth for checking auth
func CheckAuth(c *gin.Context, db *gorm.DB) error {
	tokenAuth, err := ExtractTokenMetadata(c.Request)
	if err != nil {
		return err
	}
	_, err = FetchAuth(tokenAuth, db)
	if err != nil {
		return err
	}
	return nil
}

// CreateToken function generating token for the user
func createToken(userID uint64) (*types.TokenDetails, error) {
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

func createAuth(userid uint64, td *types.TokenDetails, db *gorm.DB) {
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
		// expTime, er := claims["exp"].(int64)
		// if !er {
		// 	return nil, err
		// }
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
	return tokenMeta.UserID, nil
}
