package main

import (
	"database/sql"
	"goApi/config"
	"goApi/model"
	"log"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
	"goApi/controllers"
	"fmt"
)

var dbmap = initDb()

func initDb() *gorp.DbMap {
	fmt.Println(config.ConnectString)
	db, err := sql.Open("mysql", config.ConnectString)
	checkErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(model.User{}, "User").SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")
	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func main() {
	r := gin.Default()
	r.Use(Cors())
	c := controllers.NewUserController(dbmap)
	v1 := r.Group("api/v1")
	{
		v1.GET("/users", c.GetUsers)
		v1.GET("/users/:id", c.GetUser)
		v1.POST("/users", c.PostUser)
		v1.PUT("/users/:id", c.UpdateUser)
		v1.DELETE("/users/:id", c.DeleteUser)
		v1.OPTIONS("/users", OptionsUser)     // POST
		v1.OPTIONS("/users/:id", OptionsUser) // PUT, DELETE
	}

	r.Run(":2000")
}

func OptionsUser(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE,POST, PUT")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Next()
}
