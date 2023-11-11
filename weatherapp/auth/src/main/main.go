package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/abohmeed/auth/authdb"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const secretkey string = "xco0sr0fh4e52x03g9mv"

var dbHost string
var dbRoot = "root"
var dbPassword string

type Token struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
}

func main() {
	if os.Getenv("DB_HOST") != "" {
		dbHost = os.Getenv("DB_HOST")
	}
	if os.Getenv("DB_PASSWORD") != "" {
		dbPassword = os.Getenv("DB_PASSWORD")
	}
	db := authdb.Connect(dbRoot, dbPassword, dbHost)
	authdb.CreateDB(db)
	authdb.CreateTables(db)
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowMethods("OPTIONS")
	router.Use(cors.New(corsConfig))
	router.GET("/", health)
	router.POST("/users/:id", loginUser)
	router.POST("/users", createUser)
	router.Run(":8080")
}

type UserCreds struct {
	Username string `json:"user_name"`
	Password string `json:"user_password"`
}

func health(c *gin.Context) {
	db := authdb.Connect(dbRoot, dbPassword, dbHost)
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to the database"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": "The auth is running"})
	}
}

func loginUser(c *gin.Context) {
	var uc UserCreds
	err := c.BindJSON(&uc)
	if err != nil {
		fmt.Println("Received invalid JSON for user login")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect or invalid JSON"})
		return
	}
	encPasswordb := md5.Sum([]byte(uc.Password))
	encPassword := hex.EncodeToString(encPasswordb[:])
	db := authdb.Connect(dbRoot, dbPassword, dbHost)
	u, err := authdb.GetUserByName(uc.Username, db)
	if err != nil {
		fmt.Println(err)
	}
	if u != (authdb.User{}) && u.Password == encPassword {
		token, err := GenerateJWT(u.Name)
		if err != nil {
			fmt.Println("Error while generating the token: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"JWT": token})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bad credentials"})
	}
}
func createUser(c *gin.Context) {
	var u authdb.User
	c.BindJSON(&u)
	db := authdb.Connect(dbRoot, dbPassword, dbHost)
	result, err := authdb.CreateUser(db, u)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while adding the user. Please check the logs"})
		return
	} else if !result {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "User already exists"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"success": "User added successfully"})
	}
}
func GenerateJWT(userName string) (string, error) {
	var mySigningKey = []byte(secretkey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["username"] = userName
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}
