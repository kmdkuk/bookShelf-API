package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	gorm.Model
	Name     string    `json:"name"`
	Age      int       `json:age`
	Birthday time.Time `json:birthday`
}

func gormConnect() *gorm.DB {
	DBMS := "mysql"
	USER := "user"
	PASS := "password"
	PROTOCOL := "tcp(localhost:3306)"
	DBNAME := "sample_db"
	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME
	db, err := gorm.Open(DBMS, CONNECT)

	if err != nil {
		panic(err.Error())
	}
	fmt.Println("db connected: ", &db)
	return db
}

func main() {

	db := gormConnect()

	defer db.Close()

	db.LogMode(true)
	db.Set("gorm:table_options", "ENGINE=InnoDB")
	db.AutoMigrate(&User{})

	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world")
	})

	//CREATE
	r.POST("/user", func(c *gin.Context) {
		user := User{}
		now := time.Now()
		user.CreatedAt = now
		user.UpdatedAt = now

		err := c.BindJSON(&user)
		if err != nil {
			c.String(http.StatusBadRequest, "Request is failed: "+err.Error())
		}
		db.NewRecord(user)
		db.Create(&user)
		if db.NewRecord(user) == false {
			c.JSON(http.StatusOK, user)
		}
	})

	r.Run(":8080")
}
