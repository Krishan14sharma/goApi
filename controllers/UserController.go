package controllers

import (
	"gopkg.in/gorp.v1"
	"github.com/gin-gonic/gin"
	"goApi/model"
	"log"
	"strconv"
)

type UserController struct {
	db *gorp.DbMap
}

func NewUserController(db *gorp.DbMap) *UserController {
	return &UserController{db} // creating the user controller and sending the reference back
}

func (u *UserController)GetUsers(c *gin.Context) {
	var users []model.User
	_, err := u.db.Select(&users, "SELECT * FROM User")

	if err == nil {
		c.JSON(200, users)
	} else {
		c.JSON(200, gin.H{"error": "no user(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/users
}
func (u *UserController)UpdateUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user model.User
	err := u.db.SelectOne(&user, "SELECT * FROM User WHERE id=?", id)

	if err == nil {
		var json model.User
		c.Bind(&json)

		user_id, _ := strconv.ParseInt(id, 0, 64)

		user := model.User{
			Id:        user_id,
			Firstname: json.Firstname,
			Lastname:  json.Lastname,
		}

		if user.Firstname != "" && user.Lastname != "" {
			_, err = u.db.Update(&user)

			if err == nil {
				c.JSON(200, user)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/users/1
}
func (u *UserController)PostUser(c *gin.Context) {
	var user model.User
	c.Bind(&user)

	log.Println(user)

	if user.Firstname != "" && user.Lastname != "" {

		if insert, _ := u.db.Exec(`INSERT INTO User (firstname, lastname) VALUES (?, ?)`, user.Firstname, user.Lastname); insert != nil {
			user_id, err := insert.LastInsertId()
			if err == nil {
				content := &model.User{
					Id:        user_id,
					Firstname: user.Firstname,
					Lastname:  user.Lastname,
				}
				c.JSON(201, content)
			} else {
				checkErr(err, "Insert failed")
			}
		}

	} else {
		c.JSON(400, gin.H{"error": "Fields are empty"})
	}
	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/users
}

func (u *UserController)GetUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user model.User
	err := u.db.SelectOne(&user, "SELECT * FROM User WHERE id=? LIMIT 1", id)
	if err == nil {
		user_id, _ := strconv.ParseInt(id, 0, 64)

		content := &model.User{
			Id:        user_id,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
		}
		c.JSON(200, content)
	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}

	// curl -i http://localhost:8080/api/v1/users/1
}

func (u *UserController)DeleteUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user model.User
	err := u.db.SelectOne(&user, "SELECT * FROM User WHERE id=?", id)
	if err == nil {
		_, err = u.db.Delete(&user)
		if err == nil {
			c.JSON(200, gin.H{"id #" + id: "deleted"})
		} else {
			checkErr(err, "Delete failed")
		}
	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}
	// curl -i -X DELETE http://localhost:8080/api/v1/users/1
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}